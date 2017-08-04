package main

import (
	"net"
	"fmt"
	"os"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:5567")
	if  err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	conn.Write([]byte("Hello,I'm client!"))

	fmt.Println("send msg")
	buf := make([]byte, 50)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn closed")
		return
	}
	fmt.Println("recv msg:", string(buf[0:n]))
}