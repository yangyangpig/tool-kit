package main

import (
	"log"
	api "toolkit/demon/kitex/echo/kitex_gen/api/echo"
)

func main() {
	svr := api.NewServer(new(EchoImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
