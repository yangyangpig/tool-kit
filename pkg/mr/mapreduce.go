package mr

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

const (
	defaultRoutine = 10
	minRoutine     = 1
)

var (
	ErrDefaultCancelWithNil  = errors.New("mapreduce cancelled with nil")
	ErrDefaultReduceNoOutput = errors.New("reduce not writing value")
)

// generate--(source chan)->map--(write)-->reduce--(write)-->out
type (
	Write interface {
		Pipe(v interface{})
	}
	GenerateFunc func(source chan<- interface{})
	MapFunc      func(item interface{}, write Write)

	MapperFunc func(item interface{}, write Write, cancel func(error))

	ReducerFunc     func(pipe <-chan interface{}, write Write, cancel func(error))
	VoidReducerFunc func(pipe <-chan interface{}, cancel func(error))

	Option           func(opts *mapReduceOptions)
	mapReduceOptions struct {
		routines int
	}
)

func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc,
	opts ...Option) (interface{}, error) {

	// 生产数据
	rawData := processRawData(generate)

	fmt.Printf("process raw data %d\n", rawData)
	return MapReduceWithRawData(rawData, mapper, reducer, opts...)
}

// 没有值返回
func MapReduceVoid(generate GenerateFunc, mapper MapperFunc, reducer VoidReducerFunc, opts ...Option) error {
	_, err := MapReduce(generate, mapper, func(input <-chan interface{}, write Write, cancel func(error)) {
		reducer(input, cancel)
		drain(input)
		write.Pipe(SinglePlaceholder) // 由于这里reducer函数没有返回数值，所以MapReduceWithRawData中没有向返回的channel输入，根据单向发送通道类型,会阻塞
	}, opts...)
	return err
}

func MapReduceWithRawData(rawData <-chan interface{},
	mapper MapperFunc, reducer ReducerFunc, opts ...Option) (interface{}, error) {
	options := buildOptions(opts...)
	output := make(chan interface{}, options.routines)
	collector := make(chan interface{}, options.routines)

	done := NewDoneChan()

	outputPipeline := NewPipeLine(output, done.Done()) //output作为pipe的channel

	var closeOnce sync.Once
	var logVal AtomicStoreErr

	finish := func() {
		closeOnce.Do(func() {
			done.Close()
			close(output)
		})
	}

	cancel := once(func(err error) {
		if err != nil {
			fmt.Printf("cacel happen error: %+v", err)
			logVal.Set(err)
		} else {
			logVal.Set(ErrDefaultCancelWithNil)
		}
		drain(rawData)
		finish()
	})

	// 执行reducer
	go func() {
		defer func() {
			if r := recover(); r != nil {
				cancel(fmt.Errorf("%v", r))
			} else {
				finish()
			}
		}()
		reducer(collector, outputPipeline, cancel)
		drain(collector)
	}()

	// 执行mapper
	go executeMapper(func(item interface{}, w Write) {
		fmt.Printf("execute mapper %+v\n", item)
		mapper(item, w, cancel)
	}, rawData, collector, done.Done(), options.routines)

	value, ok := <-output

	if err := logVal.Load(); err != nil {
		// TODO 打印日志，从外层输出日志内容
		fmt.Printf("happen error: %+v", err)
		return nil, err
	} else if ok {
		// 外层需要对value进行断言
		return value, nil
	} else {
		return nil, ErrDefaultReduceNoOutput
	}

}

func executeMapper(mapper MapFunc, input <-chan interface{}, collector chan<- interface{},
	done <-chan SinglePlaceholderType, routines int) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(collector)
	}()

	launcher := make(chan SinglePlaceholderType, routines)

	collectorPipe := NewPipeLine(collector, done)

	for {
		select {
		case <-done:
			return
		case launcher <- SinglePlaceholder:
			item, ok := <-input
			if !ok {
				// input的channel被关闭
				<-launcher
				return
			}
			wg.Add(1)
			go func() {
				defer func() {
					if p := recover(); p != nil {
						// TODO
						fmt.Printf("mapper defer happen error: %+v", p)
					}
					wg.Done()
					<-launcher
				}()
				fmt.Printf("mapper run here: %+v\n", item)

				mapper(item, collectorPipe)
			}()
		}
	}
}

type PipeLine struct {
	channel chan<- interface{}
	done    <-chan SinglePlaceholderType
}

func NewPipeLine(channel chan<- interface{}, done <-chan SinglePlaceholderType) PipeLine {
	return PipeLine{
		channel: channel,
		done:    done,
	}
}

func (p PipeLine) Pipe(v interface{}) {
	select {
	case <-p.done:
		return
	default:
		p.channel <- v

	}
}

func processRawData(generate GenerateFunc) chan interface{} {
	source := make(chan interface{})
	go func() {
		defer func() {
			if p := recover(); p != nil {
				log.Fatal(p)
			}
			// 为了安全起见，由发送者主动关闭channel
			close(source)
		}()

		generate(source)

	}()

	return source
}

func buildOptions(opts ...Option) *mapReduceOptions {

	options := &mapReduceOptions{
		routines: defaultRoutine,
	}
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithRoutines(number int) Option {
	return func(opts *mapReduceOptions) {
		if number < minRoutine {
			opts.routines = minRoutine
		} else {
			opts.routines = number
		}
	}
}
