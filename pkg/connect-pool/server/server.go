package server

import (
	"context"
	"gather/tool-kitcl/pkg/connect-pool/protocbuf"
)



type HelloWorld struct {

}

func NewHelloWorld() *HelloWorld {
	return &HelloWorld{}
}

func (h *HelloWorld)SayHello(ctx context.Context, in *tool_pkg_pool.HelloRequest) (*tool_pkg_pool.HelloResp, error) {
	resp := &tool_pkg_pool.HelloResp{}
	resp.Msg = in.Name
	resp.Code = 200
	return resp, nil
}

