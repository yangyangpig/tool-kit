package mr

import "sync"

func once(fn func(error)) func(error) {
	once := new(sync.Once)
	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}
// 目的清空channel，避免阻塞和泄漏
func drain(channel <-chan interface{}) {
	// drain the channel
	for range channel {
	}
}

