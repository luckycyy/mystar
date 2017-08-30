package main
import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"encoding/json"
	"time"
)

func main() {
	http.HandleFunc("/dati", DatiHandler)
	http.HandleFunc("/msg", ButtonHandler)
	http.HandleFunc("/wsquery", WsQueryHandler)
	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Print("server running.")
	http.ListenAndServe(":5569", nil)
}
func DatiHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into DatiHandler")
	time.Sleep(2*time.Second)
	req.ParseForm()
	if len(req.Form["v"]) > 0 {
		v:=string(req.Form["v"][0])
		if v=="1true"{
			log.Println("1true event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%AC%AC%E4%B8%80%E9%A2%98%E5%9B%9E%E7%AD%94%E6%AD%A3%E7%A1%AE&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="1false"{
			log.Println("1false event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%AC%AC%E4%B8%80%E9%A2%98%E5%9B%9E%E7%AD%94%E9%94%99%E8%AF%AF&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="2true"{
			log.Println("2true event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%AC%AC%E4%BA%8C%E9%A2%98%E5%9B%9E%E7%AD%94%E6%AD%A3%E7%A1%AE&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="2false"{
			log.Println("2false event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%AC%AC%E4%BA%8C%E9%A2%98%E5%9B%9E%E7%AD%94%E9%94%99%E8%AF%AF&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="3true"{
			log.Println("3true event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%AC%AC%E4%B8%89%E9%A2%98%E5%9B%9E%E7%AD%94%E6%AD%A3%E7%A1%AE&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if v=="3false"{
			log.Println("3false event")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%AC%AC%E4%B8%89%E9%A2%98%E5%9B%9E%E7%AD%94%E9%94%99%E8%AF%AF&user_action=true")
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
			log.Println("pintu9 finish")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E4%B9%9D%E6%8C%89%E9%92%AE%E8%BD%AC%E6%8B%BC%E5%9B%BE%E6%B8%B8%E6%88%8F%E5%AE%8C%E6%88%90&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
			break
		}
	}
	log.Println("finish.close conn")
}

