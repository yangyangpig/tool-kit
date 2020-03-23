package main

import (
	"fmt"
	"log"
	"time"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
)

type echoServer struct {
	*gnet.EventServer
	pool *goroutine.Pool
}
// 当 server 初始化完成之后调用
func (es *echoServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

// 当连接被打开的时候调用。
func (es *echoServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action)  {
	cx := c.Context()
	if v, ok := cx.(string); ok {
		log.Println("context value %s", v)
	}
	addr := c.LocalAddr()
	log.Printf("network %s addrress %s", addr.Network(), addr.String())

	remoteAddr := c.RemoteAddr()
	log.Printf("remote addr network %s remote addrress %s", remoteAddr.Network(), remoteAddr.String())

	// TODO 为什么读不到数据??
	data := c.Read()

	log.Printf("server receive data %s", string(data))
	return
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	data := append([]byte{}, frame...)

	log.Printf("react frame %s", string(frame))
	// Use ants pool to unblock the event-loop.
	_ = es.pool.Submit(func() {
		time.Sleep(1 * time.Second)
		c.AsyncWrite(data)
	})

	return
}

type echoDecode struct {}

func (d *echoDecode) Encode(c gnet.Conn, buf []byte) ([]byte, error)  {
	if c.Context() == nil {
		return buf, nil
	}
	var msg string
	if e, ok := c.Context().(string); ok {
		msg = e
	}

	return []byte(fmt.Sprintf("context is not nil erro msg %s", msg)), nil
}

func (d *echoDecode) Decode(c gnet.Conn) ([]byte, error)  {
	buf := c.Read()
	c.ResetBuffer()
	log.Printf("server receive data %s", string(buf))
	return buf, nil
}

func main() {
	p := goroutine.Default()
	defer p.Release()

	echo := &echoServer{pool: p}
	echoDecode := &echoDecode{}
	log.Fatal(gnet.Serve(echo, "tcp://:9000", gnet.WithMulticore(true)), gnet.WithCodec(echoDecode))
}
