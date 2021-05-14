package threading

import (
	"sync"
)

// 封装一组goroutine同步控制，通过全局的sync.WaitGroup，控制全部goroutine结束后，才退出来

type RoutineGroup struct {
	waitGroup sync.WaitGroup
}

func NewRoutineGroup() *RoutineGroup {
	return new(RoutineGroup)
}

// 不要在fn中把外部参数传递进来，因为外部参数会被其它goroutine改变，如果想在goroutine中传递参数，
// golang推荐是通过channel，而不是变量
func (g *RoutineGroup) Run(fn func()) {
	g.waitGroup.Add(1)

	go func() {
		defer g.waitGroup.Done()

		fn()
	}()
}

func (g *RoutineGroup) RunSafe(fn func()) {
	g.waitGroup.Add(1)

	go func() {
		defer func() {
			if p := recover(); p != nil {

			}
		}()

		defer g.waitGroup.Done()
		fn()
	}()
}

func (g *RoutineGroup) Wait() {
	g.waitGroup.Wait()
}

