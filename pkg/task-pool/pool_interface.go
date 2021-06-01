package task_pool

import "runtime"

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

var taskChanCap = func() int {
	// Use blocking channel if GOMAXPROCS=1.
	// This switches context from sender to receiver immediately,
	// which results in higher performance (under go1.5 at least).
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}

	// Use non-blocking workerChan if GOMAXPROCS>1,
	// since otherwise the sender might be dragged down if the receiver is CPU-bound.
	return 1
}()
