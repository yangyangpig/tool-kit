package threading

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRoutineGroup_Run(t *testing.T) {
	var count int32

	group := NewRoutineGroup()

	for i:=0; i < 3; i++ {
		group.Run(func() {
			atomic.AddInt32(&count, 1)
		})
	}

	group.Wait()

	assert.Equal(t, int32(3), count, "the two worlds should be same")
}

func TestRoutineGroup_RunSafe(t *testing.T) {
	var count int32

	group := NewRoutineGroup()
	var once sync.Once
	for i:=0; i< 3; i++ {
		group.RunSafe(func() {
			once.Do(func() {
				panic("") // 这个panic可以被捕获不报错
			})

			atomic.AddInt32(&count, 1)
		})
	}

	group.Wait()

	assert.Equal(t, int32(2), count)
}