package mr

import "sync/atomic"

type AtomicStoreErr struct {
	err atomic.Value
}

func (as *AtomicStoreErr) Set(err error) {
	as.err.Store(err)
}

func (as *AtomicStoreErr) Load() error  {
	if v := as.err.Load(); v != nil {
		return v.(error)
	}
	return nil
}
