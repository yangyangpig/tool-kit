package queue

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	consumers = 1
	rounds = 100
)

func TestQueue(t *testing.T)  {
	mockProducer := newMockedProducer(rounds)
	mockConsumer := newMockedConsumer()
	q := NewQueue(func() (Producer, error) {
		return mockProducer, nil
	}, func() (Consumer, error) {
		return mockConsumer, nil
	})
	mockConsumer.wait.Add(consumers)
	q.AddListener(newMockProducerListener(mockProducer))
	q.SetQueueName("mockqueue")
	q.SetNumProducer(1)
	q.SetNumConsumer(consumers)

	var mainWait sync.WaitGroup
	mainWait.Add(1)
	go func() {
		defer mainWait.Done()
		q.Start()
	}()


	mainWait.Wait()

	fmt.Printf("current rounds %d consumer count %d", rounds, mockConsumer.count)

	//assert.Equal(t, int32(rounds), atomic.LoadInt32(&mockConsumer.count))
}

type (
	MockProducer struct {
		total int32 // 用于发送次数控制
		wg sync.WaitGroup
		count int32
		pause int32
	}

	MockConsumer struct {
		count int32
		events int32
		wait sync.WaitGroup
	}

	MockProducerListener struct {
		pm *MockProducer
	}
)

func newMockedConsumer() *MockConsumer {
	return new(MockConsumer)
}
func (c *MockConsumer) Consume(string) error  {
	atomic.AddInt32(&c.count, 1)
	return nil
}

func (c *MockConsumer) OnEvent(event interface{}) {
	if atomic.AddInt32(&c.count, 1) < consumers {
		c.wait.Done()
	}
}

func newMockedProducer(total int32) *MockProducer {
	p := new(MockProducer)
	p.total = total
	p.wg.Add(int(total))
	return p
}

func (p *MockProducer) AddListener(listener ProduceListener) {

}

func (p *MockProducer) Produce() (string, bool)  {
	if p.pause ==1 {
		return "", false
	}
	if atomic.AddInt32(&p.count, 1) <= p.total {
		p.wg.Done()
		return "item:" + strconv.Itoa(int(p.pause)), true
	} else {
		time.Sleep(time.Second)
		return "", false
	}
}

func newMockProducerListener(producer *MockProducer) *MockProducerListener {
	return &MockProducerListener{
		pm:producer,
	}
}

func (l *MockProducerListener) OnPause()  {
	atomic.AddInt32(&l.pm.pause, 1)

}

func (l *MockProducerListener) OnResume()  {
	atomic.AddInt32(&l.pm.pause, -1)
}


