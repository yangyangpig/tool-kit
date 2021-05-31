package sync

import (
	"fmt"
	"sync"
	"time"
)

type (
	PoolOption func(*Pool)
	
	node struct {
		item interface{}
		next *node
		lastUsed time.Duration
	}
	
	Pool struct {
		limit int // 控制池子大小
		created int // 已经创建的元素
		maxAge time.Duration // 控制池子中，元素的最大生命周期
		lock sync.Locker // 互斥锁
		cond *sync.Cond // 协程同步
		head *node // 单向链表
		create func() interface{} // 外部定义创建方法
		destroy func(interface{}) // 外部定义删除方法
	}
)

func NewPool(n int, create func()interface{}, destroy func(interface{}), opts ...PoolOption) *Pool {
	if n <= 0 {
		panic("pool size can't be negative or zero")
	}
	lock := new(sync.Mutex)
	pool := &Pool{
		limit: n,
		lock: lock,
		cond: sync.NewCond(lock),
		create: create,
		destroy: destroy,
	}

	for _, opt := range opts {
		opt(pool)
	}
	return pool
}
// get一定要放到goroutine中，才能发挥sync.cond的作用
func (p *Pool) Get() interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()

	for {
		if p.head != nil {
			head := p.head
			p.head = head.next
			if p.maxAge >0 && p.maxAge + head.lastUsed < time.Since(time.Now().AddDate(-1,-1,-1)) {
				p.created--
				p.destroy(head.item)
				continue
			} else {
				return head.item
			}
		}

		if p.created < p.limit {
			p.created++
			return p.create()
		}
		fmt.Printf("current created %d\n",p.created)
		// 池子元素全部被拿来使用还没归还，并且池子已经满了，此时，会组塞当前的goroutine，直到有元素重新返回再唤醒
		p.cond.Wait()
	}
}

func (p *Pool) Put(x interface{}) {
	if x == nil {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()

	p.head = &node{
		item: x,
		next: p.head,
		lastUsed: time.Since(time.Now().AddDate(-1,-1,-1)),
	}

	p.cond.Signal()
}

func WithMaxAge(duration time.Duration) PoolOption {
	return func(pool *Pool) {
		pool.maxAge = duration
	}
}
