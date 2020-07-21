package main

import (
	"fmt"
	"net"
)

func main() {
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

type Client struct {
	Cn net.Conn
}

func NewClient() (*Client, error) {
	client := &Client{}
	c, err := net.Dial("tcp", ":9000")
	if err != nil {
		fmt.Println(err)
		return client, err
	}
	client.Cn = c
	return client, nil
}

func (c *Client) PingPong() {
	_, err := c.Cn.Write([]byte("this is test"))
	if err != nil {
		fmt.Println("write the failed", err)
		return
	}
	// defer c.cli.Close()
	buff := make([]byte, 512)

	n, err := c.Cn.Read(buff)
	if err != nil {
		fmt.Println("Read failed:", err)
		return
	}
	fmt.Println("count:", n, "msg:", string(buff))

}
