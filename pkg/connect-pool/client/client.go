package client

import (
	"context"
	"fmt"
	"gather/tool-kitcl/pkg/connect-pool/protocbuf"
	"github.com/0x5010/grpcp"
	"google.golang.org/grpc"
	"log"
)

type HelloWorldCli struct {
	gcon *grpc.ClientConn
	sub  tool_pkg_pool.HelloWorldClient
}

func NewHelloWorldClient() *HelloWorldCli {
	conn, _ := grpcp.GetConn("127.0.0.1:8080")
	client := tool_pkg_pool.NewHelloWorldClient(conn)

	return &HelloWorldCli{
		gcon: conn,
		sub:  client,
	}
}

func (c *HelloWorldCli) SayHello(ctx context.Context, in *tool_pkg_pool.HelloRequest) (*tool_pkg_pool.HelloResp, error)  {
	resp, err := c.sub.SayHello(ctx, in)
	if err != nil {
		log.Printf("sayHello error %v", err)
		return resp, err
	}

	fmt.Println(resp)
}
