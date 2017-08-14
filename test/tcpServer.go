package main

import (
	"fmt"
	"net"
	"os"

	"net/http"
	"log"
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
		msg := string(buf[0:n])
		fmt.Println("recv msg:", msg)
		if msg == "" {

		}
		conn.Write([]byte("Hi,I'm server!"))
	}
}

type JHD struct {
	Num      string //净化点编号
	Color    string
	Active   bool //点亮，未点亮。受on off控制
	Disabled bool //失效
}
type Player struct {
	WxId     string //微信id
	Num      string //衣服编号
	GloveNum string //手套编号
	Team     string
	Status   string
}
type Msg struct {
	From     string //消息来自,ZDF战斗服 JHD净化点
	Num      string //编号
	GloveNum string //手套编号
	Status   string //on off
}

var jhds []JHD
var players []Player
func main() {

	http.HandleFunc("/reset", ResetHandler)
	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Print("server running.")
	http.ListenAndServe(":5568", nil)

	listen_sock, err := net.Listen("tcp", "192.168.1.21:5567")
	if err != nil {
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
func ResetHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into ResetHandler")
	ResetGameStatus()
	fmt.Fprint(w, "reset ok")
	log.Printf("reset ok")
}
func ResetGameStatus() {
	currentStatus.Step = 0
	currentStatus.Status = "ready"
	currentStatus.From = ""
	currentStatus.To = ""
	bs.ABtnStatus = ""
	bs.BBtnStatus = ""
}