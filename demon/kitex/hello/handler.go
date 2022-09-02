package main

import (
	"context"
	hello "toolkit/demon/kitex/hello/kitex_gen/kitex/hello"
)

// HelloWorldImpl implements the last service interface defined in the IDL.
type HelloWorldImpl struct{}

// SayHello implements the HelloWorldImpl interface.
func (s *HelloWorldImpl) SayHello(ctx context.Context, req *hello.HelloRequest) (resp *hello.HelloResp, err error) {
	// TODO: Your code here...
	resp = &hello.HelloResp{
		Msg:  req.GetName(),
		Code: 0,
	}
	return
}
