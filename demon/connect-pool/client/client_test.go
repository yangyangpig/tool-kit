package client

import (
	"context"
	"fmt"
	"gather/tool-kitcl/demon/connect-pool/protocbuf"
	"google.golang.org/grpc"
	"log"
	"testing"
)

var dailF = func(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}

func TestHelloWorldCli_SayHello(t *testing.T) {
	cli := NewHelloWorldClient()
	resp, err := cli.SayHello(context.Background(), &tool_pkg_pool.HelloRequest{Name: "zhangsan"})
	if err != nil {
		log.Fatalf("invoke server error %v", err)
		return
	}
	fmt.Println(resp)
}

func BenchmarkHelloWorldCli_SayHello(b *testing.B) {
	// cli := NewHelloWorldClient()
	// b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli := NewHelloWorldClient()
		_, err := cli.SayHello(context.Background(), &tool_pkg_pool.HelloRequest{Name: "zhangsan"})
		if err != nil {
			log.Fatalf("invoke server error %v", err)
			return
		}
		//fmt.Println(resp)
	}
}

// connect with grpc dail
func BenchmarkHelloWorldCli_SayHello2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cc, err := dailF(":8080")
		if err != nil {
			b.Fatal(err)
		}
		stub := tool_pkg_pool.NewHelloWorldClient(cc)
		_, err = stub.SayHello(context.Background(), &tool_pkg_pool.HelloRequest{Name: "name"})
		if err != nil {
			b.Fatal(err)
		}
	}
}
