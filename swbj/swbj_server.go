package main
import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"encoding/json"
)

func main() {
	http.HandleFunc("/swbj", SwbjHandler)
	http.HandleFunc("/wsquery", WsQueryHandler)
	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Print("server running.")
	http.ListenAndServe(":5573", nil)
}

func SwbjHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into SwbjHandler")
	req.ParseForm()
	if len(req.Form["v"]) > 0 {
		val:=string(req.Form["v"][0])
		if val=="btncall" {
			log.Println("  SwbjHandler v is btncall")

			//修改为触发播放剧情的url
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%8B%A8%E9%80%9A110%E7%94%B5%E8%AF%9D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()


		}else if(val=="end"){
			log.Println(" SwbjHandler v is end")
			conn.WriteMessage(1, []byte("end"))

		}else if(val=="note1"){
			log.Println(" SwbjHandler v is note1")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%8B%A8%E9%80%9A110%E7%94%B5%E8%AF%9D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note2"){
			log.Println(" SwbjHandler v is note2")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%8B%A8%E9%80%9A110%E7%94%B5%E8%AF%9D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note3"){
			log.Println(" SwbjHandler v is note3")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%8B%A8%E9%80%9A110%E7%94%B5%E8%AF%9D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note4"){
			log.Println(" SwbjHandler v is note4")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%8B%A8%E9%80%9A110%E7%94%B5%E8%AF%9D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note5"){
			log.Println(" SwbjHandler v is note5")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%8B%A8%E9%80%9A110%E7%94%B5%E8%AF%9D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note6"){
			log.Println(" SwbjHandler v is note6")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%8B%A8%E9%80%9A110%E7%94%B5%E8%AF%9D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}
	}

}

var conn *websocket.Conn
func WsQueryHandler(w http.ResponseWriter, req *http.Request) {
	var upgrader = websocket.Upgrader{}
	conn, _ = upgrader.Upgrade(w, req, nil)

	defer conn.Close()
	for {
		_, msg, _ := conn.ReadMessage()
		if string(msg) == "success" {
			log.Println("SwbjHandler finish")
		//	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.43&group=action_st&st=%E6%89%8B%E6%9C%BA%E7%BB%88%E7%AB%AF%E6%8B%BC%E5%9B%BE%E6%B8%B8%E6%88%8F%E5%AE%8C%E6%88%90&user_action=true")
		//	if err != nil {
		//		print(err)
		//	}
		//	resp.Body.Close()
		//	break
		}
	}
	log.Println("finish.close conn")
}
