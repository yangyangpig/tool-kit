package main

import (
	"fmt"
	"gather/tool-kitcl/pkg/task-pool"
	"sync"
)

func main()  {
	pool, err := task_pool.NewPool(func(option *task_pool.Option) {
		option.InitWorkerNum = 10
	})
	if err != nil {
		fmt.Printf("init pool fail %v", err)
		return
	}
	fmt.Println(pool.PoolStatus())
	var wg sync.WaitGroup
	var handle = func() {
		fmt.Println("click.....")
		wg.Done()
	}
	wg.Add(1)
	pool.Fill(handle)
	pool.PoolStatus()
	wg.Wait()

	fmt.Println(pool.PoolStatus())

}
