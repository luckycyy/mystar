package main
import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

func main() {
	connMap=make(map[string]*websocket.Conn)
	http.HandleFunc("/reset", KtResetHandler)
	http.HandleFunc("/msg", MsgHandler)
	http.HandleFunc("/ws", WsHandler)
	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Print("kt 5571 server running.")
	http.ListenAndServe(":5571", nil)
}

func MsgHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into ButtonHandler")
	req.ParseForm()
	//11 12 13 failed 1ok
	if len(req.Form["v"]) > 0 {
		v:=string(req.Form["v"][0])
		log.Println("v is:"+v)
		if(v=="1ok"){
			log.Println("gt_kt 1 ok")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.43&group=action_st&st=%E6%8A%95%E6%94%BE%E6%8A%97%E4%BD%93%E5%B7%A6%E8%BE%B9%E5%AE%8C%E6%88%90&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
			connMap["kt1"].WriteMessage(1, []byte("success"))
		}else if(v=="2ok"){
			log.Println("gt_kt 2 ok")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.43&group=action_st&st=%E6%8A%95%E6%94%BE%E6%8A%97%E4%BD%93%E4%B8%AD%E9%97%B4%E5%AE%8C%E6%88%90&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
			connMap["kt2"].WriteMessage(1, []byte("success"))
		} else if(v=="3ok"){
			log.Println("gt_kt 3 ok")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.43&group=action_st&st=%E6%8A%95%E6%94%BE%E6%8A%97%E4%BD%93%E5%8F%B3%E8%BE%B9%E5%AE%8C%E6%88%90&user_action=true")
			if err != nil {
				print(err)
			}
			connMap["kt3"].WriteMessage(1, []byte("success"))
			resp.Body.Close()
		}else if(v=="11")||(v=="12")||(v=="13")||(v=="1f"){
			connMap["kt1"].WriteMessage(1, []byte(v))
		}else if(v=="21")||(v=="22")||(v=="23")||(v=="2f"){
			connMap["kt2"].WriteMessage(1, []byte(v))
		}else if(v=="31")||(v=="32")||(v=="33")||(v=="3f"){
			connMap["kt3"].WriteMessage(1, []byte(v))
		}
	}

}
var connMap map[string]*websocket.Conn

func WsHandler(w http.ResponseWriter, req *http.Request) {
	var upgrader = websocket.Upgrader{}
	var conn *websocket.Conn
	conn, _ = upgrader.Upgrade(w, req, nil)

	defer conn.Close()
	for {
		_, msg, _ := conn.ReadMessage()//此处有问题
		if string(msg) == "kangti1" {
			log.Println("kangti1 conn")
			connMap["kt1"]= conn
		}else if string(msg) == "kangti2" {
			log.Println("kangti2 conn")
			connMap["kt2"]= conn
		}else if string(msg) == "kangti3" {
			log.Println("kangti3 conn")
			connMap["kt3"]= conn
		}
	}
	log.Println("finish.close conn")
}

func KtResetHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into kt reset")
	connMap["kt1"].WriteMessage(1, []byte("reload"))
	connMap["kt2"].WriteMessage(1, []byte("reload"))
	connMap["kt3"].WriteMessage(1, []byte("reload"))
}
