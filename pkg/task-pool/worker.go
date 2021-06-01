package task_pool

type Worker struct {
	taskChan chan Handle
	p        PoolInterface
}

func NewWorker(p PoolInterface, taskCap int) *Worker {
	return &Worker{
		taskChan: make(chan Handle, taskCap), // 指定的大小
		p:        p,
	}
}

func (w *Worker) Start() {
	go func() {
		// TODO 考虑放回的处理逻辑
		defer func() {
			w.p.DecrementBusyWorkerNum()
			w.p.PutWorkerCache(w)
			// TODO 处理pain的recover
			w.p.SignalCond()
		}()
		for {
			handle, ok := <-w.taskChan
			// TODO 可以考虑taskChan被close时候处理
			if !ok {
				break
			}
			handle()
			if ok := w.p.RevertWorker(w); ok {
				return
			}
		}

	}()
}

func (w *Worker) Stop() {
	close(w.taskChan)
}

func (w *Worker) Go(handle Handle) {
	w.taskChan <- handle
}


