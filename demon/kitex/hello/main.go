package main

import (
	"log"
	hello "toolkit/demon/kitex/hello/kitex_gen/kitex/hello/helloworld"
)

func main() {
	svr := hello.NewServer(new(HelloWorldImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
