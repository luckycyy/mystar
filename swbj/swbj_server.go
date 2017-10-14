package main
import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/swbj", SwbjHandler)
	http.HandleFunc("/wsquery", WsQueryHandler)
	http.HandleFunc("/wsquerygamestatus", WsQueryGameStatusHandler)
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
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E6%8E%A5%E9%80%9A%E7%94%B5%E8%AF%9D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()


		}else if(val=="juqingend"){
			log.Println(" SwbjHandler v is juqingend")
			conn.WriteMessage(1, []byte("juqingend"))

		}else if(val=="mobileend"){
			log.Println(" SwbjHandler v is mobileend")
			conn.WriteMessage(1, []byte("mobileend"))

		}else if(val=="note1"){
			log.Println(" SwbjHandler v is note1")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%AC%AC%E4%B8%80%E6%9D%A1%E7%95%99%E8%A8%80&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note2"){
			log.Println(" SwbjHandler v is note2")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%AC%AC%E4%BA%8C%E6%9D%A1%E7%95%99%E8%A8%80&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note3"){
			log.Println(" SwbjHandler v is note3")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%AC%AC%E4%B8%89%E6%9D%A1%E7%95%99%E8%A8%80&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note4"){
			log.Println(" SwbjHandler v is note4")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%AC%AC%E5%9B%9B%E6%9D%A1%E7%95%99%E8%A8%80&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note5"){
			log.Println(" SwbjHandler v is note5")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%AC%AC%E4%BA%94%E6%9D%A1%E7%95%99%E8%A8%80&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="note6"){
			log.Println(" SwbjHandler v is note6")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%AC%AC%E5%85%AD%E6%9D%A1%E7%95%99%E8%A8%80&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="gamefinish"){
			log.Println(" SwbjHandler v is gamefinish")
			conn2.WriteMessage(1, []byte("gamefinish"))
			log.Println(" ws write to browser gamefinish")
		}else if(val=="poyifailed"){
			log.Println(" SwbjHandler v is poyifailed")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%A0%B4%E8%AF%91%E9%94%99%E8%AF%AF&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="poyisuccess"){
			log.Println(" SwbjHandler v is poyisuccess")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%A0%B4%E8%AF%91%E6%88%90%E5%8A%9F&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="boxingchange"){
			log.Println(" SwbjHandler v is boxingchange")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E6%B3%A2%E5%BD%A2%E5%88%87%E6%8D%A2%E9%9F%B3%E6%95%88&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="mimasuccess"){
			log.Println(" SwbjHandler v is mimasuccess")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%94%B5%E8%84%91%E5%B1%8F%E5%B9%95%E5%AF%86%E7%A0%81%E6%AD%A3%E7%A1%AE&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="mimafailed"){
			log.Println(" SwbjHandler v is mimafailed")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%94%B5%E8%84%91%E5%B1%8F%E5%B9%95%E5%AF%86%E7%A0%81%E9%94%99%E8%AF%AF&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="huakuang4click"){
			log.Println(" SwbjHandler v is huakuang4click")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%82%B9%E5%87%BB%E7%94%BB%E6%A1%86%E4%B8%80&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="huakuang3click"){
			log.Println(" SwbjHandler v is huakuang3click")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%82%B9%E5%87%BB%E7%94%BB%E6%A1%86%E4%B8%89&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="huakuang2click"){
			log.Println(" SwbjHandler v is huakuang2click")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%82%B9%E5%87%BB%E7%94%BB%E6%A1%86%E4%BA%8C&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if(val=="huakuang1click"){
			log.Println(" SwbjHandler v is huakuang1click")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E7%AC%94%E8%AE%B0-%E7%82%B9%E5%87%BB%E7%94%BB%E6%A1%86%E5%9B%9B&user_action=true")
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


var conn2 *websocket.Conn
func WsQueryGameStatusHandler(w http.ResponseWriter, req *http.Request) {
	var upgrader = websocket.Upgrader{}
	conn2, _ = upgrader.Upgrade(w, req, nil)

	defer conn2.Close()
	for {
		_, msg, _ := conn2.ReadMessage()
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
	log.Println("finish.close conn2")
}