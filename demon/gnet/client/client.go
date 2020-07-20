package main

import (
	"fmt"
	"net"
)

func main()  {
	c, err := net.Dial("tcp", ":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	_, err = c.Write([]byte("this is test"))
	if err != nil {
		fmt.Println("write the failed", err)
		return
	}
	buff := make([]byte, 512)

	n, err := c.Read(buff)
	if err != nil {
		fmt.Println("Read failed:", err)
		return
	}
	fmt.Println("count:", n, "msg:", string(buff))

}
