package threading

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)
func TestRoutinePool(t *testing.T)  {
	pool := NewTaskRunner(1)
	times := 1000
	var counter uint32
	var waitGroup sync.WaitGroup
	for i:=0; i < times; i++ {
		waitGroup.Add(1)
		pool.Schedule(func() {
			atomic.AddUint32(&counter, 1)
			fmt.Printf("current goroutine number: %d\n", counter)
			waitGroup.Done()
		})
	}
	waitGroup.Wait()
}

func BenchmarkRoutinePool(b *testing.B)  {
	queue := NewTaskRunner(runtime.NumCPU())

	for i:=0; i < b.N; i++ {
		queue.Schedule(func() {

		})
	}
}
