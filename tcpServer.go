package main

import (
	"fmt"
	"net"
	"os"
)

func recvConnMsg(conn net.Conn) {
	buf := make([]byte, 50)
	defer conn.Close()
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn closed")
			return
		}
		fmt.Println("recv msg:", string(buf[0:n]))
		conn.Write([]byte("Hi,I'm server!"))
	}
}

func main() {
	listen_sock, err := net.Listen("tcp", "127.0.0.1:5567")
	if  err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
	defer listen_sock.Close()
	fmt.Println("tcp server running v2")
	for {
		new_conn, err := listen_sock.Accept()
		if err != nil {
			continue
		}
		fmt.Println("conn ok")
		go recvConnMsg(new_conn)
	}

}