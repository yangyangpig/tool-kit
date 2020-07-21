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
// test benchmen
// go test -bench=. -benchmem=true -run=none
// 测试制定函数
// go test -bench="funcName" -benchmem=true -run=none
// bench 标记接受一个表达式作为参数，.表示运行所有的基准测试
// -benchtime 指定运行时长单位秒
// -run 匹配一个从来没有的单元测试方法过滤掉单元测试的输出
// go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out
// -memprofile 输出内存监控文件
// -cpuprofile 输出cpu监控文件
// go tool pprof profile.out 查看监控文件
/**
result:
BenchmarkHelloWorldCli_SayHello-4           6850            177907 ns/op            4828 B/op         99 allocs/op 使用了grpc池子
BenchmarkHelloWorldCli_SayHello2-4          1234           1036429 ns/op          118972 B/op        286 allocs/op 普通的grpc连接
ns/op 平均每次迭代锁消耗纳秒数
B/op 平均每次迭代内存所分配的字节数
allocs/op 平均每次迭代的内存分配次数
 */
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
