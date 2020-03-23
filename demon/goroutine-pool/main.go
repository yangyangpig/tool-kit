package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"sync/atomic"
	"time"
)

var sum int32

func myFunc(i interface{})  {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func demoFunc()  {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

func main()  {
	//defaultPool := use_case.NewCommonUse(
	//	//	use_case.WithHandleFunc(func() {
	//	//	fmt.Println("hello world!!!!!!")
	//	//}),
	//	//use_case.WithHandleWithParam(
	//	//	func(g interface{}) {
	//	//		fmt.Println("hello world with the param !!!! \n")
	//	//		if k, ok := g.(*sync.WaitGroup); ok {
	//	//			k.Done()
	//	//		}
	//	//	}),
	//	//)
	//	//defaultPool.Run()

	commonPool()
}

func commonPool()  {
	p, err := ants.NewPool(100, ants.WithNonblocking(true))
	if err != nil {
		fmt.Println("new pool fail error %v", err)
		return
	}

	err = p.Submit(func() {
		fmt.Println("test test test!!!")
	})

	if err != nil {
		fmt.Println("submit method fail error %v", err)
		return
	}
	fmt.Println(p.Running())
	fmt.Printf("pool cap %d", p.Cap())

}

func common()  {
	defer ants.Release()

	runTimes := 1000

	// Use the common pool

	var wg sync.WaitGroup
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}

	for i:=0; i< runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum)
	}

	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all task.\n")

	// use the pool with a function
	// set 10 to the capacity of goroutine pool and 1 second for expired duration
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
		wg.Done()
	})

	defer p.Release()

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}
	wg.Wait()
	fmt.Printf("running gorotuines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
}


