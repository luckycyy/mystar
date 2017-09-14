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
	http.ListenAndServe(":5570", nil)
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
			log.Println("pintugt finish")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.43&group=action_st&st=%E6%89%8B%E6%9C%BA%E7%BB%88%E7%AB%AF%E6%8B%BC%E5%9B%BE%E6%B8%B8%E6%88%8F%E5%AE%8C%E6%88%90&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
			break
		}
	}
	log.Println("finish.close conn")
}


/*音效
http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.43&group=action_st&st=%E8%BD%AC%E6%8B%BC%E5%9B%BE%E9%9F%B3%E6%95%88&user_action=true
 */
