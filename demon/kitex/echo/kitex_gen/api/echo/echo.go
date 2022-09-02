// Code generated by Kitex v0.4.2. DO NOT EDIT.

package echo

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	api "toolkit/demon/kitex/echo/kitex_gen/api"
)

func serviceInfo() *kitex.ServiceInfo {
	return echoServiceInfo
}

var echoServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "Echo"
	handlerType := (*api.Echo)(nil)
	methods := map[string]kitex.MethodInfo{
		"echo": kitex.NewMethodInfo(echoHandler, newEchoEchoArgs, newEchoEchoResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "api",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.4.2",
		Extra:           extra,
	}
	return svcInfo
}

func echoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.EchoEchoArgs)
	realResult := result.(*api.EchoEchoResult)
	success, err := handler.(api.Echo).Echo(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newEchoEchoArgs() interface{} {
	return api.NewEchoEchoArgs()
}

func newEchoEchoResult() interface{} {
	return api.NewEchoEchoResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Echo(ctx context.Context, req *api.Request) (r *api.Response, err error) {
	var _args api.EchoEchoArgs
	_args.Req = req
	var _result api.EchoEchoResult
	if err = p.c.Call(ctx, "echo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}