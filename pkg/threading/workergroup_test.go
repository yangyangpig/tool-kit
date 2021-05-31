package threading

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

func TestWorkerGroup_Start(t *testing.T) {
	var wg sync.WaitGroup
	var lock sync.Mutex
	m := make(map[string]SinglePlaceholderType)
	wg.Add(runtime.NumCPU())
	group := NewWorkerGroup(func() {
		lock.Lock()
		m[fmt.Sprint(RoutineId())] = SinglePlaceholder
		fmt.Printf("current m %+v\n", m)
		lock.Unlock()
		wg.Done()
	}, runtime.NumCPU())
	go group.Start()
	wg.Wait()
}


func RoutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	// if error, just return 0
	n, _ := strconv.ParseUint(string(b), 10, 64)

	return n
}

// 测试一下数组裁剪通用逻辑

type Elem struct {

}
