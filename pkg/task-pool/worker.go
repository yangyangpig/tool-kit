package task_pool

type Worker struct {
	taskChan chan Handle
	p        PoolInterface
}

func NewWorker(p PoolInterface, taskCap int) *Worker {
	return &Worker{
		taskChan: make(chan Handle, 1), // 指定的大小
		p:        p,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			handle, ok := <-w.taskChan
			if !ok {
				break
			}
			handle()
		}

	}()
	w.p.IncrementIdleWorkerNum()
}

func (w *Worker) Stop() {
	close(w.taskChan)
}

func (w *Worker) Go(handle Handle) {
	w.taskChan <- handle
}
