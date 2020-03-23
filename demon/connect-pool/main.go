package main

import (
	"gather/tool-kitcl/demon/connect-pool/protocbuf"
	"gather/tool-kitcl/demon/connect-pool/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()

	tool_pkg_pool.RegisterHelloWorldServer(s, server.NewHelloWorld())

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen tcp :8080")
	}

	log.Println("serving on :8080")
	log.Println(s.Serve(l))
}
