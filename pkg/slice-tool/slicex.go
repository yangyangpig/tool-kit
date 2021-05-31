package slice_tool

import (
	"fmt"
	"time"
)

// worker是一个具体的执行实体对象
type Worker struct {
	Fn func()
	RecycleTime time.Time
}

type Elem struct {
	contain []*Worker
	size int
}

func NewElem(size int) *Elem {
	return &Elem{
		contain: make([]*Worker, 0, size),
		size: size,
	}
}

func (e *Elem) Len() int{
	return len(e.contain)
}

func (e *Elem) IsEmpty() bool  {
	return len(e.contain) == 0
}

func (e *Elem) Insert(data *Worker) error {
	e.contain = append(e.contain, data)
	return nil
}

func (e *Elem) detach() interface{}  {
	l := e.Len()
	if l == 0 {
		return 0
	}

	el := e.contain[l-1]

	e.contain[l-1] = nil //
	e.contain = e.contain[:l-1]
	return el
}

func (e *Elem) RetrieveExpiry(duration time.Duration) []*Worker {
	num := e.Len()
	// 获取指定时间之前的worker，为过期时间
	expireTime := time.Now().Add(-duration)

	index := e.binarySearch(0, num-1, expireTime)

	if expireTime.Before(e.contain[index].RecycleTime) {
		fmt.Printf("expireTime woker %+v", e.contain[index])
	}

	if expireTime.Before(e.contain[index+1].RecycleTime) {
		fmt.Printf("not expire time woker %+v", e.contain[index+1])
	}
	var expireWorker []*Worker
	//expireWorker = e.contain[:0] // 声明一个空的容器
	if index != -1 {
		expireWorker = append(expireWorker, e.contain[:index+1]...)
		m := copy(e.contain, e.contain[index+1:]) // 复制未超时的,只会按照长度小对应的位置复制
		for i:=m; i< num; i++ {
			e.contain[i] = nil //gc 回收
		}
		e.contain = e.contain[:m]
	}
	return expireWorker
}

// 二分查找,可能contain的池子expire不是有序的，可以通过sort方法排序
func (e *Elem) binarySearch(l, r int, expireTime time.Time) int {
	var mid int
	for l <= r {
		mid = (l + r) / 2
		// TODO maybe need to lock
		if expireTime.Before(e.contain[mid].RecycleTime) {
			r = mid -1
		} else {
			l = mid + 1
		}
	}
	return r
}
