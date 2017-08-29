package main
import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"encoding/json"
)

func main() {
	http.HandleFunc("/dati", DatiHandler)
	http.HandleFunc("/msg", ButtonHandler)
	http.HandleFunc("/wsquery", WsQueryHandler)
	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Print("server running.")
	http.ListenAndServe("192.168.1.21:5569", nil)
}
func DatiHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into DatiHandler")
	req.ParseForm()
	if len(req.Form["v"]) > 0 {
		v:=string(req.Form["v"][0])
		if v=="1ture"{
			log.Println("1true event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E8%AE%A1%E6%97%B6%E5%BC%80%E5%A7%8B%E5%89%8D%E5%BE%80B%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="1false"{
			log.Println("1false event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E8%AE%A1%E6%97%B6%E5%BC%80%E5%A7%8B%E5%89%8D%E5%BE%80B%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="2true"{
			log.Println("2true event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E8%AE%A1%E6%97%B6%E5%BC%80%E5%A7%8B%E5%89%8D%E5%BE%80B%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="2false"{
			log.Println("2false event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E8%AE%A1%E6%97%B6%E5%BC%80%E5%A7%8B%E5%89%8D%E5%BE%80B%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="3true"{
			log.Println("3true event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E8%AE%A1%E6%97%B6%E5%BC%80%E5%A7%8B%E5%89%8D%E5%BE%80B%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="3false"{
			log.Println("3false event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E8%AE%A1%E6%97%B6%E5%BC%80%E5%A7%8B%E5%89%8D%E5%BE%80B%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}
	}

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

