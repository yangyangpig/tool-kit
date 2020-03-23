package use_case

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

type HandleFunc func()
type CommonUseOption func(*CommonUse)
type HandleWithParam func(g interface{})

func DefaultHandleFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

func DefaultHandleWithParam(g interface{}) {
	fmt.Println("hello world with param!!!")
	if k, ok := g.(*sync.WaitGroup); ok {
		k.Done()
	}

}

type CommonUse struct {
	wg              *sync.WaitGroup
	handle          HandleFunc
	handleWithParam HandleWithParam
}

func WithHandleFunc(hf HandleFunc) CommonUseOption {
	return func(use *CommonUse) {
		use.handle = hf
	}
}

func WithHandleWithParam(hfp HandleWithParam) CommonUseOption {
	return func(use *CommonUse) {
		use.handleWithParam = hfp
	}
}

func NewCommonUse(opts ...CommonUseOption) *CommonUse {
	cu := &CommonUse{
		wg:              &sync.WaitGroup{},
		handle:          DefaultHandleFunc,
		handleWithParam: DefaultHandleWithParam,
	}

	for _, opt := range opts {
		opt(cu)
	}

	return cu
}

func (c *CommonUse) Run() {
	defer ants.Release()

	runTime := 1000

	//for i := 0; i < runTime; i++ {
	//	c.wg.Add(1)
	//	_ = ants.Submit(c.handle) // 不能传递带参数的函数，有瓶颈
	//}
	//
	//fmt.Printf("running goroutines: %d\n", ants.Running())
	//fmt.Printf("finish all tasks.\n")


	// 一个函数绑定一个池子的goroutine处理
	var p *ants.PoolWithFunc
	if c.handleWithParam != nil {
		p, _ = ants.NewPoolWithFunc(10, c.handleWithParam)

		for i := 0; i < runTime; i++ {
			c.wg.Add(1)
			_ = p.Invoke(c.wg)
		}
	}
	c.wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())

	// 声明一个池子，处理不同函数
	prep, _ := ants.NewPool(1000, ants.WithPreAlloc(true))

	defer prep.Release()

	for i := 0; i < runTime; i++ {
		c.wg.Add(i)
		// 唯一缺点，这个方式，指定的函数不能够传参数
		_ = prep.Submit(func() {
			fmt.Println("this is a test")
		})
	}
	prep.Running()
}
