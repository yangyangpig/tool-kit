package task_pool

type Handle func()

type PoolInterface interface {
	Fill(handle Handle)
	PoolStatus() (idleWorkerNum int64, busyWorkerNum int64)
	FreeIdleWorkers()
	IncrementIdleWorkerNum()
	DecrementIdleWorkerNum()
	IncrementBusyWorkerNum()
	DecrementBusyWorkerNum()
}
