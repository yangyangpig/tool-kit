package threading

import "runtime"

var SinglePlaceholder SinglePlaceholderType
type SinglePlaceholderType = struct {}

// 使用场景：指定容量
type TaskRunner struct {
	limitChan chan SinglePlaceholderType
}

func NewTaskRunner(concurrency int) *TaskRunner  {
	if concurrency == 0 {
		// 根据cpu的个数，消费者处理能力指定concurrency
		concurrency = taskChanCap
	}
	return &TaskRunner{
		limitChan: make(chan SinglePlaceholderType, concurrency),
	}
}

// 如果concurrency=1,可以保证goroutine顺序执行,否则，无法保证
func (tr *TaskRunner) Schedule(task func()) {
	tr.limitChan <- SinglePlaceholder

	go func() {
		defer Recover(func() {
			<-tr.limitChan
		})
	}()

	task()
}

var taskChanCap = func() int {
	// Use blocking channel if GOMAXPROCS=1.
	// This switches context from sender to receiver immediately,
	// which results in higher performance (under go1.5 at least).
	// Inspired by fasthttp at
	// https://github.com/valyala/fasthttp/blob/master/workerpool.go#L139
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}

	// Use non-blocking workerChan if GOMAXPROCS>1,
	// since otherwise the sender might be dragged down if the receiver is CPU-bound.
	return 1
}()


