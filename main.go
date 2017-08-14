package main
import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"time"
	"github.com/gorilla/websocket"
	"starroom/music_controller"
	"io/ioutil"
	"strings"
)
type GameStatus struct {
	Step   int    //阶段
	Status string //状态
	From   string //起点
	To     string //目标点
}
type BtnStatus struct {
	ABtnStatus string
	BBtnStatus string
}
var currentStatus GameStatus
var bs BtnStatus
var conn *websocket.Conn
func loadConf()string{//有问题
	data, err := ioutil.ReadFile("conf/server.conf")
	if err != nil {
		fmt.Println("readConf error:", err)
	}
	log.Print("server_conf data:"+string(data))
	server_path:=strings.Split(string(data),"=")[1]
	return string(server_path)
}
func main() {
	http.HandleFunc("/pt", PintuHandler)
	http.HandleFunc("/set", ButtonHandler)
	http.HandleFunc("/pass", PassHandler)
	http.HandleFunc("/query", QueryHandler)
	http.HandleFunc("/wsquery", WsQueryHandler)
	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Print("server running.")
	http.ListenAndServe(":5566", nil)
}
func PassHandler(w http.ResponseWriter, req *http.Request){
	req.ParseForm()
	if len(req.Form["game"]) > 0 {
		if string(req.Form["game"][0]) == "zhefanpao" {
			passStatus, _ := json.Marshal(GameStatus{Status:"end"})
			conn.WriteMessage(1, passStatus)
			time.Sleep(4 * time.Second)
			music_controller.PlayFinish()
			log.Println("pass game zhefanpao")
		}
	}else{
		fmt.Fprint(w, "PassHandler param must more than 1")
	}
}
func ButtonHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into ButtonHandler")
	req.ParseForm()
	if len(req.Form["aBtn"]) > 0 {
		if string(req.Form["aBtn"][0]) == "true" {
			if bs.ABtnStatus != "true" { //只在A按钮状态改变时进入，屏蔽重复输入
				//状态改变啦
				if bs.ABtnStatus == "" { //游戏第一轮开始时进入
					//游戏开始
					bs.ABtnStatus = "true"
					currentStatus.Status = "start"
					currentStatus.Step = 1
					currentStatus.From = "A"
					currentStatus.To = "B"
					music_controller.PlayMusicToB()
					sendCurrentStatus()

				}
				//非第一轮，ready状态时
				if currentStatus.To == "B" && currentStatus.Status == "ready" {
					currentStatus.Status = "start"
					bs.BBtnStatus="false"//去谁那就把谁变false
					music_controller.PlayMusicToB()
					sendCurrentStatus()
				}
				if currentStatus.Status == "start" && currentStatus.To == "A" {
					currentStatus.Status = "ready"
					bs.BBtnStatus = "false"
					bs.ABtnStatus = "false"
					currentStatus.Step++
					if currentStatus.Step == 2 {
						music_controller.PlayStep1ok()
					} else if currentStatus.Step == 3 {
						music_controller.PlayStep2ok()
					} else if currentStatus.Step == 4 {
						currentStatus.Status = "end"
						sendCurrentStatus()
						time.Sleep(4 * time.Second)
						music_controller.PlayFinish()
						return
					}
					reverseDestination()
					sendCurrentStatus()
				}
			}
		}
	}
	if len(req.Form["bBtn"]) > 0 {
		if string(req.Form["bBtn"][0]) == "true" {

			if bs.BBtnStatus != "true" { //只在B按钮状态改变时进入，屏蔽重复输入
				//状态改变啦

				//非第一轮，ready状态时
				if currentStatus.To == "A" && currentStatus.Status == "ready" {
					currentStatus.Status = "start"
					bs.ABtnStatus = "false"//去谁那就把谁变false
					music_controller.PlayMusicToA()
					sendCurrentStatus()

				}
				if currentStatus.Status == "start" && currentStatus.To == "B" {
					bs.BBtnStatus = "false"
					bs.ABtnStatus = "false"
					currentStatus.Step++
					if currentStatus.Step == 2 {
						music_controller.PlayStep1ok()
					} else if currentStatus.Step == 3 {
						music_controller.PlayStep2ok()
					}
					if currentStatus.Step == 4 {
						currentStatus.Status = "end"
						sendCurrentStatus()
						time.Sleep(4 * time.Second)
						music_controller.PlayFinish()
						return
					} else {
						currentStatus.Status = "ready"
					}
					reverseDestination()
					sendCurrentStatus()
				}
			}
		}
	}
	//fmt.Fprint(w, btns.ABtnStatus+","+btns.BBtnStatus+","+btns.Step)
}
func ResetHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into ResetHandler")
	ResetGameStatus()
	fmt.Fprint(w, "reset ok")
	log.Printf("reset ok")
}

func reverseDestination() {
	if currentStatus.From == "A" {
		currentStatus.From = "B"
	} else if currentStatus.From == "B" {
		currentStatus.From = "A"
	}
	if currentStatus.To == "A" {
		currentStatus.To = "B"
	} else if currentStatus.To == "B" {
		currentStatus.To = "A"
	}
}
func ResetGameStatus() {
	currentStatus.Step = 0
	currentStatus.Status = "ready"
	currentStatus.From = ""
	currentStatus.To = ""
	bs.ABtnStatus = ""
	bs.BBtnStatus = ""
}
func QueryHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into QueryHandler")
	s, err := json.Marshal(currentStatus)
	if err != nil {
		fmt.Println("json.Mashal error:", err)
	}
	fmt.Fprint(w, string(s))
}
func WsQueryHandler(w http.ResponseWriter, req *http.Request) {
	var upgrader = websocket.Upgrader{}
	conn, _ = upgrader.Upgrade(w, req, nil)
	ResetGameStatus()
	sendCurrentStatus()
	defer conn.Close()
	for {
		_, msg, _ := conn.ReadMessage()
		if string(msg) == "failed" {
			currentStatus.Status = "ready"
			reverseDestination()
			sendCurrentStatus()
			if currentStatus.To == "A" {
				music_controller.PlayFailedToB()
			} else if currentStatus.To == "B" {
				music_controller.PlayFailedToA()
			}
		}
		log.Println(string(msg))
	}
}
func sendCurrentStatus() {
	jsonCurrentStatus, _ := json.Marshal(currentStatus)
	conn.WriteMessage(1, jsonCurrentStatus)
	jsonbs,_:=json.Marshal(bs)
	log.Println(string(jsonbs))
}
func PintuHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into PintuHandler")
	req.ParseForm()
	if len(req.Form["player"]) > 0 {
		if string(req.Form["player"][0]) == "right" {
			log.Println("right http.get")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E6%8B%BC%E5%9B%BER%E4%BE%A7%E5%AE%8C%E6%88%90&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
			fmt.Fprint(w, "right finish")
			log.Printf("right ok")
		}
		if string(req.Form["player"][0]) == "left" {
			log.Println("left http.get")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E6%8B%BC%E5%9B%BEL%E4%BE%A7%E5%AE%8C%E6%88%90&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
			fmt.Fprint(w, "left finish")
			log.Println("left ok")
		}
	}
}
