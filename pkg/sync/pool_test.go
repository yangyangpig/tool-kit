package sync

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const limit = 10

var Placeholder PlaceholderType
type (
	GenericType     = interface{}
	PlaceholderType = struct{}
)

type PprofHttpServer struct {
	pprofAddr string
	httpServer *http.Server
	HTTPListener net.Listener

	quit chan struct{}
}

func init()  {
	//pp := &PprofHttpServer{
	//	pprofAddr: ":15080",
	//	quit: make(chan struct{}),
	//}
	//err := pp.init()
	//if err != nil {
	//	fmt.Printf("init pprof http server error %v", err)
	//	return
	//}
	//pp.Start()

}

func (s *PprofHttpServer) init() error {
	s.httpServer = &http.Server{}
	var err error
	s.HTTPListener, err = net.Listen("tcp", s.pprofAddr)
	if err != nil {
		return err
	}
	return nil
}

func (s *PprofHttpServer) Start()  {
	go func() {
		err := s.httpServer.Serve(s.HTTPListener)
		fmt.Printf("http server done: %v", err)
	}()

}

func (s *PprofHttpServer) Close()  {
	s.quit <- struct{}{}
}

func TestPoolGet(t *testing.T)  {

	stack := NewPool(limit, create, destroy)
	ch := make(chan interface{})

	for i:=0; i< limit; i++ {
		go func() {
			v := stack.Get()
			if v.(int) != 1 {
				t.Fatal("unmatch value")
			}
			ch <- struct{}{}
		}()

		select {
		case <-ch:
		case <-time.After(time.Second * 3):
			t.Fail()
		}
	}

}

func TestPoolPopMore(t *testing.T)  {

	stack := NewPool(limit, create, destroy)
	ch := make(chan PlaceholderType, 1)
	for i:=0; i < limit; i++ {
		var wait sync.WaitGroup
		wait.Add(1)
		go func() {
			stack.Get()
			ch <- Placeholder
			wait.Done()
		}()

		wait.Wait()
		select {
		case <-ch:
			fmt.Printf("here is done\n")
		default:
			t.Fail()
		}
	}

	// 主goroutine等待上面那些goroutine
	var waitGroup, pushWait sync.WaitGroup//waitGroup.Add(1)
	pushWait.Add(1)
	//waitGroup.Add(1)
	go func() {
		pushWait.Done()
		v := stack.Get()
		fmt.Printf("current v %d", v)
		waitGroup.Done()
	}()
	pushWait.Wait()
	stack.Put(1)
	//waitGroup.Wait()
}

func TestPoolPopFirst(t *testing.T)  {
	var value int32
	stack := NewPool(limit, func() interface{} {
		return atomic.AddInt32(&value, 1)
	}, destroy, WithMaxAge(time.Second * 12))

	for i:=0; i < limit; i++ {
		v := stack.Get().(int32)
		fmt.Printf("current value %d\n", v)
		stack.Put(v)
	}
}

func create() interface{} {
	return 1
}

func destroy(_ interface{})  {

}