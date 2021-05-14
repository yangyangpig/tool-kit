package queue

import (
	"fmt"
	"gather/tool-kitcl/pkg/threading"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const queueName = "queue"

type (
	Queue struct {
		name string
		producerFactory ProducerFactory
		consumerFactory ConsumerFactory
		producerRoutineGroup *threading.RoutineGroup
		consumerRoutineGroup *threading.RoutineGroup

		producerCount int
		consumerCount int
		active int32
		channel chan string
		quit chan struct{}

		listeners []Listener
		eventLock sync.Mutex
		eventChannels []chan interface{}

	}

	Listener interface {
		OnPause()
		OnResume()
	}
)

func NewQueue(producerFactory ProducerFactory, consumerFactory ConsumerFactory) *Queue  {
	queue := &Queue{
		producerFactory: producerFactory,
		consumerFactory: consumerFactory,
		producerRoutineGroup: threading.NewRoutineGroup(),
		consumerRoutineGroup: threading.NewRoutineGroup(),
		producerCount: runtime.NumCPU(),
		consumerCount: runtime.NumCPU()<<1,
		channel: make(chan string),
		quit: make(chan struct{}),
	}

	queue.SetQueueName(queueName)
	return queue
}

func (q *Queue) SetQueueName(name string)  {
	q.name = name
}

func (q *Queue)AddListener(listener Listener)  {
	q.listeners = append(q.listeners, listener)
}

func (q *Queue) Broadcast(message interface{}) {
	go func() {
		q.eventLock.Lock()
		defer q.eventLock.Unlock()

		for _, channel := range q.eventChannels {
			channel <- message
		}
	}()
}

func (q *Queue) SetNumConsumer(count int)  {
	q.consumerCount = count
}

func (q *Queue) SetNumProducer(count int) {
	q.producerCount = count
}

func (q *Queue) Start() {
	q.StartProducers(q.producerCount)
	q.StartConsumer(q.consumerCount)
	q.producerRoutineGroup.Wait()
	close(q.channel)
	q.consumerRoutineGroup.Wait()
}

func (q *Queue) Close()  {
	close(q.quit)
}

func (q *Queue) StartProducers(number int) {
	for i:=0; i < number; i++ {
		q.producerRoutineGroup.Run(func() {
			// TODO produce()
			q.produce()
		})
	}
}

func (q *Queue) produce() {

	var producer Producer

	// producer工厂生产一个producer
	for {
		var err error
		if producer, err = q.producerFactory(); err != nil {
			fmt.Printf("Error on creating producer: %v", err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	// 记录一下活跃的producer个数
	atomic.AddInt32(&q.active, 1)
	producer.AddListener(routineListener{
		queue: q,
	})

	// 生产数据
	for {
		select {
		case <- q.quit:
			fmt.Printf("Quitting producer")
		default:
			if v, ok := q.produceOne(producer); ok {
				q.channel <- v
			}
		}
		return
	}
}

func (q *Queue) produceOne(p Producer) (string, bool){
	defer func() {
		if p := recover(); p != nil {

		}
	}()

	return p.Produce()
}

func (q *Queue) StartConsumer(number int)  {
	for i:=0; i< number; i++ {
		eventChan := make(chan interface{})
		q.eventLock.Lock()
		q.eventChannels = append(q.eventChannels, eventChan)
		q.eventLock.Unlock()

		q.consumerRoutineGroup.Run(func() {
			q.consume(eventChan)
		})
	}
}

func (q *Queue) consume(eventChan chan interface{})  {
	var consumer Consumer

	for {
		var err error

		// consumer工厂生产一个consumer
		if consumer, err = q.consumerFactory(); err != nil {
			fmt.Printf("Error on creationg consumer: %v", err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	for  {
		select {
		case message, ok := <- q.channel:
			if ok {
				q.consumeOne(consumer, message)
			} else {
				fmt.Printf("Task channel was closed")
				return
			}
		}
	}
}

func (q *Queue) consumeOne(consumer Consumer, message string)  {
	// 发生panic不用管
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// TODO
			}
		}()

		if err := consumer.Consume(message); err != nil {
			fmt.Printf("consume message error: %v", err)
		}
	}()
}

func (q *Queue) pause() {

}

func (q *Queue) resume() {

}

type routineListener struct {
	queue *Queue
}

func (rl routineListener) OnProducerPause()  {
	if atomic.AddInt32(&rl.queue.active, -1) <= 0 {
		rl.queue.pause()
	}
}

func (rl routineListener) OnProducerResume() {
	if atomic.AddInt32(&rl.queue.active, 1) == 1 {
		rl.queue.resume()
	}
}