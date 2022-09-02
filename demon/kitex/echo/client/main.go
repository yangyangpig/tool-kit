package main

import (
	"context"
	"log"
	"time"
	"toolkit/demon/kitex/echo/kitex_gen/api"
	"toolkit/demon/kitex/echo/kitex_gen/api/echo"

	"github.com/cloudwego/kitex/client"
)
func main() {
	client, err := echo.NewClient("echo", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		log.Fatal(err)
	}
	for {
		req := &api.Request{Message: "hell ketix"}
		resp, err := client.Echo(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
		time.Sleep(time.Second)

	}

}
