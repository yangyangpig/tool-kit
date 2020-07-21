package main

import (
	"fmt"
	"testing"
)

var cli *Client

func init() {
	c, err := NewClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	cli = c
}
// go test -v client_test.go client.go -test.run TestClient_PingPong 指定测试某个函数
func TestClient_PingPong(t *testing.T) {
	cli.PingPong()
}

// go test bench=. -benchmem=true -run=none
func BenchmarkClient_PingPong(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli.PingPong()
	}
}
