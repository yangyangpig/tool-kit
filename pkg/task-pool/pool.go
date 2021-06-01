package task_pool

import (
	"sync"
	"sync/atomic"
)

// 这个pool主要是用来处理task，可以发送task到pool，然后就不控制了
// 使用场景：可以指定一组可用于执行task的goroutine的pool，可以复用goroutine，避免过多goroutine
type Pool struct {
	m              sync.Locker
	idleWorkerNum  int64
	busyWorkerNum  int64
	blockWorkerNum uint32
	idleWorkerList []*Worker
	workerCache  sync.Pool
	cond *sync.Cond
	opt *Option
}

type Option struct {
	MaxWorkerNum int64
	Nonblocking bool
	MaxBlockingTasks uint32
}

var defaultOption = &Option{
	MaxWorkerNum: 0,
}

type ModOption func(option *Option)

func WithMaxWorkerNum(num int64) ModOption {
	return func(option *Option) {
		option.MaxWorkerNum = num
	}
}

func WithNonblocking(nonblocking bool) ModOption  {
	return func(option *Option) {
		option.Nonblocking = nonblocking
	}
}

func WithMaxBlockingTasks(num uint32) ModOption {
	return func(option *Option) {
		option.MaxBlockingTasks = num
	}
}

// init some idle worker pool
func NewPool(modOptions ...ModOption) (PoolInterface, error) {
	option := defaultOption
	for _, fn := range modOptions {
		fn(option)
	}

	// validate the option
	if err := Validate(option); err != nil {
		return nil, err
	}

	var busyNum int64 = 0
	initPool := &Pool{idleWorkerList: make([]*Worker, 0, option.MaxWorkerNum)} // 临时值
	initPool.workerCache.New = func() interface{} {
		return &Worker{
			taskChan: make(chan Handle, taskChanCap), // 指定的大小
			p:        initPool,
		}
	}
	for i := 0; i < int(option.MaxWorkerNum); i++ {
		// 创建没有handle的channel的
		w := NewWorker(initPool,taskChanCap)
		// 启动worker的协程
		w.Start()
		initPool.idleWorkerList = append(initPool.idleWorkerList, w)
	}
	initPool.busyWorkerNum = busyNum
	initPool.opt = option
	initPool.m = new(sync.Mutex)
	initPool.cond = sync.NewCond(initPool.m)
	return initPool, nil

}

func (p *Pool) Fill(handle Handle) {
	var w *Worker
	p.m.Lock()
	if len(p.idleWorkerList) != 0 {
		w = p.idleWorkerList[len(p.idleWorkerList)-1]
		p.idleWorkerList = p.idleWorkerList[0 : len(p.idleWorkerList)-1]
		p.DecrementIdleWorkerNum()
		p.IncrementBusyWorkerNum()
	}
	p.m.Unlock()

	// TODO 这里有个问题，如果idleWorkerList取完了，会不断创建新的worker，worker会不断增加，失去了池子复用特点
	if w == nil {
		w = NewWorker(p, taskChanCap) // 可以考虑使用sync.pool存储这个worker，避免多次声明对象，减少gc压力
		w.Start()
		p.IncrementBusyWorkerNum()
		p.DecrementIdleWorkerNum()
	}
	w.Go(handle)
}

func (p *Pool) FillV2(handle Handle)  {
	spanWorker := func() {
		w := p.workerCache.Get().(*Worker)
		w.Go(handle)
	}
	p.m.Lock()
	if len(p.idleWorkerList) != 0 {
		w := p.idleWorkerList[len(p.idleWorkerList)-1]
		if w != nil {
			// 池子里面还有数据
			p.m.Unlock()
			spanWorker()
			return
		}
	}

	if _, busy :=p.PoolStatus(); busy >= p.opt.MaxWorkerNum {
		if p.opt.Nonblocking {
			p.m.Unlock()
			return
		}
		// wait等待single
		retry:
			if p.opt.MaxBlockingTasks !=0 && p.opt.MaxBlockingTasks <= p.blockWorkerNum {
				return
			}

			p.blockWorkerNum++
			p.cond.Wait()
			p.blockWorkerNum--

			var nw int64
			if nw = p.busyWorkerNum; nw == 0 {
				p.m.Unlock()
				spanWorker()
				return
			}
			if len(p.idleWorkerList) == 0  {
				// 还有空闲的
				goto retry
			}
	}
	return
}

func (p *Pool) PoolStatus() (idleWorkerNum int64, busyWorkerNum int64) {
	return p.idleWorkerNum, p.busyWorkerNum
}

func (p *Pool) RevertWorker(w *Worker) bool {
	p.m.Lock()
	if p.opt.MaxWorkerNum >0 && p.busyWorkerNum > p.opt.MaxWorkerNum {
		p.m.Unlock()
		return false
	}
	p.idleWorkerList = append(p.idleWorkerList, w)
	p.cond.Signal()
	p.m.Unlock()
	return true
}

func (p *Pool) PutWorkerCache(w *Worker)  {
	p.workerCache.Put(w)
}

func (p *Pool) SignalCond()  {
	p.cond.Signal()
}

func (p *Pool) FreeIdleWorkers() {
	p.m.Lock()
	defer p.m.Unlock()
	for i := range p.idleWorkerList {
		p.idleWorkerList[i].Stop()
	}
	p.DecrementIdleWorkerNum()
}

func (p *Pool) IncrementIdleWorkerNum() {
	idleNum := atomic.AddInt64(&p.idleWorkerNum, 1)
	p.idleWorkerNum = idleNum
}

func (p *Pool) DecrementIdleWorkerNum() {
	idleNum := atomic.AddInt64(&p.idleWorkerNum, -1)
	p.idleWorkerNum = idleNum
}

func (p *Pool) IncrementBusyWorkerNum() {
	busyNum := atomic.AddInt64(&p.busyWorkerNum, 1)
	p.busyWorkerNum = busyNum
}

func (p *Pool) DecrementBusyWorkerNum() {
	busyNum := atomic.AddInt64(&p.busyWorkerNum, -1)
	p.busyWorkerNum = busyNum
}
