package main

import (
	"net"

	"os"
	"time"
	"log"
)
var tm map[string]*time.Timer
var cm map[string]net.Conn
func main() {
	listen_sock, err := net.Listen("tcp", "127.0.0.1:5567")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer listen_sock.Close()
	log.Println("tcp server running.listen:5567")
	for {
		//192.168.1.201-220 zdf
		//192.168.1.72 JHD
		new_conn, err := listen_sock.Accept()
		if err != nil {
			continue
		}
		//new_conn.Write([]byte("bbbbb"))
		cm=make(map[string]net.Conn)
		cm["aaa"]=new_conn
		for _,v := range cm {
			v.Write([]byte("cccc"))
		}



		log.Println("wrte finish")

		//go recvConnMsg(new_conn)
	}
}