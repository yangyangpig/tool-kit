package main

import (
	"context"
	"log"
	"time"
	"toolkit/demon/kitex/hello/kitex_gen/kitex/hello"
	"toolkit/demon/kitex/hello/kitex_gen/kitex/hello/helloworld"

	"github.com/cloudwego/kitex/client"
)

func main() {
	client, err := helloworld.NewClient("helloworld", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		log.Fatal(err)
	}
	for {
		resp, err := client.SayHello(context.Background(), &hello.HelloRequest{
			Name: "hello kitex",
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
		time.Sleep(time.Second)

	}
}
