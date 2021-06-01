package task_pool

import (
	"sync"
	"sync/atomic"
)

// 这个pool主要是用来处理task，可以发送task到pool，然后就不控制了
type Pool struct {
	m              sync.Mutex
	idleWorkerNum  int64
	busyWorkerNum  int64
	idleWorkerList []*Worker
}

type Option struct {
	InitWorkerNum int
}

var defaultOption = Option{
	InitWorkerNum: 0,
}

type ModOption func(option *Option)

// init some idle worker pool
func NewPool(modOptions ...ModOption) (PoolInterface, error) {
	option := defaultOption
	for _, fn := range modOptions {
		fn(&option)
	}

	// validate the option
	if err := Validate(option); err != nil {
		return nil, err
	}

	var busyNum int64 = 0
	initPool := &Pool{idleWorkerList: make([]*Worker, 0)} // 临时值
	for i := 0; i < option.InitWorkerNum; i++ {
		// 创建没有handle的channel的
		w := NewWorker(initPool,taskChanCap)
		// 启动worker的协程
		w.Start()
		initPool.idleWorkerList = append(initPool.idleWorkerList, w)
	}
	return &Pool{
		idleWorkerNum:  initPool.idleWorkerNum,
		busyWorkerNum:  busyNum,
		idleWorkerList: initPool.idleWorkerList,
	}, nil

}

func (p *Pool) Fill(handle Handle) {
	var w *Worker
	p.m.Lock()
	if len(p.idleWorkerList) != 0 {
		w = p.idleWorkerList[len(p.idleWorkerList)-1]
		p.idleWorkerList = p.idleWorkerList[0 : len(p.idleWorkerList)-1]
		p.DecrementIdleWorkerNum()
		p.IncrementBusyWorkerNum()
	}
	p.m.Unlock()

	if w == nil {
		w = NewWorker(p, taskChanCap)
		w.Start()
		p.IncrementBusyWorkerNum()
		p.DecrementIdleWorkerNum()
	}
	w.Go(handle)
}

func (p *Pool) PoolStatus() (idleWorkerNum int64, busyWorkerNum int64) {
	return p.idleWorkerNum, p.busyWorkerNum
}

func (p *Pool) FreeIdleWorkers() {
	p.m.Lock()
	for i := range p.idleWorkerList {
		p.idleWorkerList[i].Stop()
	}
	p.DecrementIdleWorkerNum()
}

func (p *Pool) IncrementIdleWorkerNum() {
	idleNum := atomic.AddInt64(&p.idleWorkerNum, 1)
	p.idleWorkerNum = idleNum
}

func (p *Pool) DecrementIdleWorkerNum() {
	idleNum := atomic.AddInt64(&p.idleWorkerNum, -1)
	p.idleWorkerNum = idleNum
}

func (p *Pool) IncrementBusyWorkerNum() {
	busyNum := atomic.AddInt64(&p.busyWorkerNum, 1)
	p.busyWorkerNum = busyNum
}

func (p *Pool) DecrementBusyWorkerNum() {
	busyNum := atomic.AddInt64(&p.busyWorkerNum, -1)
	p.busyWorkerNum = busyNum
}
