package timewheel

import (
"container/list"
"fmt"
"time"
)

const drainWorkers = 8

type (
	// key为click的标识，value为参数
	Execute func(key, value interface{})

	TimingWheel struct {
		interval      time.Duration
		ticker        Ticker
		slots         []*list.List
		timers        *SafeMap // 用来存储click对象
		tickedPos     int // 时间轮时针位置
		numSlots      int // 时间轮槽
		execute       Execute // 执行的函数
		setChannel    chan timingEntry
		moveChannel   chan baseEntry
		removeChannel chan interface{}
		drainChannel  chan func(key, value interface{})
		stopChannel   chan PlaceholderType
	}

	timingEntry struct {
		baseEntry
		value   interface{}
		circle  int // 圈数
		diff    int
		removed bool
	}

	baseEntry struct {
		delay time.Duration
		key   interface{}
	}

	positionEntry struct {
		pos  int
		item *timingEntry
	}

	timingTask struct {
		key   interface{}
		value interface{}
	}
)

func NewTimingWheel(interval time.Duration, numSlots int, execute Execute) (*TimingWheel, error) {
	if interval <= 0 || numSlots <= 0 || execute == nil {
		return nil, fmt.Errorf("interval: %v, slots: %d, execute: %p", interval, numSlots, execute)
	}

	return newTimingWheelWithClock(interval, numSlots, execute, NewTicker(interval))
}

func newTimingWheelWithClock(interval time.Duration, numSlots int, execute Execute, ticker Ticker) (
	*TimingWheel, error) {
	tw := &TimingWheel{
		interval:      interval, // 时间轮子间隔，单位时间
		ticker:        ticker,
		slots:         make([]*list.List, numSlots),
		timers:        NewSafeMap(),
		tickedPos:     numSlots - 1,
		execute:       execute,
		numSlots:      numSlots,
		setChannel:    make(chan timingEntry),
		moveChannel:   make(chan baseEntry),
		removeChannel: make(chan interface{}),
		drainChannel:  make(chan func(key, value interface{})),
		stopChannel:   make(chan PlaceholderType),
	}

	// 初始化轮子
	tw.initSlots()
	go tw.run()

	return tw, nil
}

func (tw *TimingWheel) Drain(fn func(key, value interface{})) {
	tw.drainChannel <- fn
}

func (tw *TimingWheel) MoveTimer(key interface{}, delay time.Duration) {
	if delay <= 0 || key == nil {
		return
	}

	tw.moveChannel <- baseEntry{
		delay: delay,
		key:   key,
	}
}

func (tw *TimingWheel) RemoveTimer(key interface{}) {
	if key == nil {
		return
	}

	tw.removeChannel <- key
}
// 设置定时任务
func (tw *TimingWheel) SetTimer(key, value interface{}, delay time.Duration) {
	if delay <= 0 || key == nil {
		return
	}

	// 延时执行
	tw.setChannel <- timingEntry{
		baseEntry: baseEntry{
			delay: delay,
			key:   key, // 定时器的标识
		},
		value: value,
	}
}

func (tw *TimingWheel) Stop() {
	close(tw.stopChannel)
}
// 所有轮中，槽的执行内容，无序
func (tw *TimingWheel) drainAll(fn func(key, value interface{})) {
	runner := NewTaskRunner(drainWorkers)
	for key, slot := range tw.slots {
		fmt.Printf("slot index %d slot %+v\n",key, slot)
		for e := slot.Front(); e != nil; {
			task := e.Value.(*timingEntry)
			fmt.Printf("task %+v\n",task)
			next := e.Next()
			slot.Remove(e)
			e = next
			if !task.removed {
				runner.Schedule(func() {
					fn(task.key, task.value)
				})
			}
		}
	}
}
// 根据延迟时间算出定时任务在时间轮子上位置，位置有两部分组成(圈数+当前位置)
func (tw *TimingWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	steps := int(d / tw.interval)
	pos = (tw.tickedPos + steps) % tw.numSlots
	circle = (steps - 1) / tw.numSlots

	fmt.Printf("getPositionAndCircle d %d tw.interval %d steps %d pos %d \n", d, tw.interval,steps,pos)
	return
}

func (tw *TimingWheel) initSlots() {
	for i := 0; i < tw.numSlots; i++ {
		tw.slots[i] = list.New()
	}
}
// 为什么只更新了safemap的item，没有对list中timingEntry更新
func (tw *TimingWheel) moveTask(task baseEntry) {
	val, ok := tw.timers.Get(task.key)
	if !ok {
		return
	}

	timer := val.(*positionEntry)
	// 如果延迟时间比时间间隔小，就直接执行定时任务算了，一般不会出现
	if task.delay < tw.interval {
		GoSafe(func() {
			tw.execute(timer.item.key, timer.item.value)
		})
		return
	}

	pos, circle := tw.getPositionAndCircle(task.delay)

	fmt.Printf("moveTask current timer.pos %+v pos %d circle %d\n",timer, pos, circle)
	if pos >= timer.pos {
		timer.item.circle = circle
		timer.item.diff = pos - timer.pos
	} else if circle > 0 {
		circle--
		timer.item.circle = circle
		timer.item.diff = tw.numSlots + pos - timer.pos
	} else {
		// 不够一圈并且比当前位置往回
		timer.item.removed = true
		newItem := &timingEntry{
			baseEntry: task,
			value:     timer.item.value,
		}
		tw.slots[pos].PushBack(newItem)
		tw.setTimerPosition(pos, newItem)
	}

	fmt.Printf("moveTask last key %s value %+v\n",task.key, task)
}
// 时针位置
func (tw *TimingWheel) onTick() {
	// 每次走一步
	tw.tickedPos = (tw.tickedPos + 1) % tw.numSlots
	l := tw.slots[tw.tickedPos]
	tw.scanAndRunTasks(l)
}

func (tw *TimingWheel) removeTask(key interface{}) {
	val, ok := tw.timers.Get(key)
	if !ok {
		return
	}

	timer := val.(*positionEntry)
	timer.item.removed = true
}

func (tw *TimingWheel) run() {
	for {
		select {
		// 时间轮按照指定的间隔滚动起来
		case <-tw.ticker.Chan():
			// fmt.Printf("timingwheel run next step\n")
			tw.onTick()
		case task := <-tw.setChannel:
			tw.setTask(&task) // 设置新的定时任务
		case key := <-tw.removeChannel:
			tw.removeTask(key) // 删除指定的定时任务
		case task := <-tw.moveChannel:
			tw.moveTask(task) // 移动定时任务
		case fn := <-tw.drainChannel:
			tw.drainAll(fn) // 清空所有轮子上的定时任务
		case <-tw.stopChannel:
			tw.ticker.Stop() // 停住轮子 不需要释放list吗
			return
		}
	}
}

func (tw *TimingWheel) runTasks(tasks []timingTask) {
	if len(tasks) == 0 {
		return
	}

	go func() {
		// TODO 要改成多个goroutine
		for i := range tasks {
			RunSafe(func() {
				tw.execute(tasks[i].key, tasks[i].value)
			})
		}
	}()
}

func (tw *TimingWheel) scanAndRunTasks(l *list.List) {
	var tasks []timingTask

	//TODO 可以改成其他数据结构 每次都要遍历全链表
	for e := l.Front(); e != nil; {
		task := e.Value.(*timingEntry)
		fmt.Printf("当前的key %s value %+v\n",task.key, task)
		if task.removed {
			next := e.Next()
			l.Remove(e)
			tw.timers.Del(task.key)
			e = next
			continue
		} else if task.circle > 0 {
			// 不是当前圈
			task.circle--
			e = e.Next()
			continue
		} else if task.diff > 0 {
			next := e.Next()
			l.Remove(e)
			// (tw.tickedPos+task.diff)%tw.numSlots
			// cannot be the same value of tw.tickedPos
			pos := (tw.tickedPos + task.diff) % tw.numSlots
			tw.slots[pos].PushBack(task)
			tw.setTimerPosition(pos, task)
			task.diff = 0
			e = next
			continue
		}

		tasks = append(tasks, timingTask{
			key:   task.key,
			value: task.value,
		})
		next := e.Next()
		l.Remove(e)
		tw.timers.Del(task.key)
		e = next
	}

	tw.runTasks(tasks)
}

func (tw *TimingWheel) setTask(task *timingEntry) {
	// 如果定时任务延迟时间比时间轮间隔时间短，以时间轮间隔为准
	if task.delay < tw.interval {
		task.delay = tw.interval
	}

	// 如果定时任务存在，
	if val, ok := tw.timers.Get(task.key); ok {
		fmt.Printf("safemap中key %s value %+v\n",task.key, task)
		entry := val.(*positionEntry)
		entry.item.value = task.value
		tw.moveTask(task.baseEntry)
	} else {
		// 定时任务还没存在
		fmt.Printf("第一次task %+v value %+v\n",task.key, task)
		pos, circle := tw.getPositionAndCircle(task.delay)
		task.circle = circle
		tw.slots[pos].PushBack(task) // 这里的task和safemap中的item指向同一个指针
		tw.setTimerPosition(pos, task)
	}
}

func (tw *TimingWheel) setTimerPosition(pos int, task *timingEntry) {
	if val, ok := tw.timers.Get(task.key); ok {
		timer := val.(*positionEntry)
		timer.pos = pos
	} else {
		tw.timers.Set(task.key, &positionEntry{
			pos:  pos,
			item: task,
		})
	}
}

