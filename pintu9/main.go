package main
import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"encoding/json"
)

func main() {
	http.HandleFunc("/msg", ButtonHandler)
	http.HandleFunc("/wsquery", WsQueryHandler)
	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Print("server running.")
	http.ListenAndServe("192.168.1.21:5569", nil)
}

func ButtonHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into ButtonHandler")
	req.ParseForm()
	if len(req.Form["btn"]) > 0 {
		num:=string(req.Form["btn"][0])
		btn:=Btn{Btn:num}
		jsonCurrentStatus, _ := json.Marshal(btn)
		conn.WriteMessage(1, jsonCurrentStatus)
	}

}
type Btn struct {
	Btn string
}

var conn *websocket.Conn
func WsQueryHandler(w http.ResponseWriter, req *http.Request) {
	var upgrader = websocket.Upgrader{}
	conn, _ = upgrader.Upgrade(w, req, nil)

	defer conn.Close()
	for {
		_, msg, _ := conn.ReadMessage()
		if string(msg) == "success" {
			log.Println(string(msg))
			break
		}
	}
	log.Println("finish.close conn")
}

