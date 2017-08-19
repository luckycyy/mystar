package other

import (
	"fmt"
	"net"
	"os"
	"net/http"
	"log"
	"encoding/json"
	"time"
	"strings"
)

func recvConnMsg(conn net.Conn) {
	buf := make([]byte, 50)
	defer conn.Close()
	for {
		n, err := conn.Read(buf)
		if err != nil {
            log.Println("conn closed")
			return
		}
		msg := string(buf[0:n])
        log.Println("recv msg:", msg)
		var newMsg Msg
		if err := json.Unmarshal([]byte(msg), newMsg); err != nil {
			log.Println("recv msg convert to json err")
		}

		if newMsg.From == "bind" {
			p := Player{Num: newMsg.Num, GloveNum: newMsg.GloveNum, Team: "blue", Active: true}
			if !HasPlayer(players, p) {
				players = append(players, p)
			}
		} else if newMsg.From == "jhd" {
			if newMsg.Num == "09" || newMsg.Num == "10" {
				SetTeamByJHD(newMsg.GloveNum, "red")
			} else {
				jhd := GetJHDByNum(newMsg.Num)
				player,exist:= GetPlayerByGloveNum(newMsg.GloveNum)
                if !exist {
                    log.Println("GetPlayerByGloveNum error")
                }
				if newMsg.Status == "on" {
					if player.Dying {
						t := zdfTimerMap[player.Num]
						t.Stop()
						delete(zdfTimerMap, player.Num)
						SendBroadcastMsg("ZDF" + player.Num + "=2\r\n")
					}
					if player.Team != jhd.Color {
						//开启timer
						t := time.AfterFunc(5*time.Second, func() { ChangeJHDColor(jhd.Num, player.Num) })
						jhdTimerMap[player.Num] = t
					}
				} else if newMsg.Status == "off" {
					t := jhdTimerMap[player.Num]
					t.Stop()
					delete(jhdTimerMap, player.Num)
				}

			}

		} else if newMsg.From == "zdf" {
			attackedPlayer := GetPlayerByNum(newMsg.Num)
			attacker ,exist:= GetPlayerByGloveNum(newMsg.GloveNum)
            if !exist {
                log.Println("GetPlayerByGloveNum error")
            }
			if (attackedPlayer.Team != attacker.Team) && attacker.Active {
				if attackedPlayer.Team == "red" {
					attackedPlayer.Active = false
					SendBroadcastMsg("ZDF" + attackedPlayer.Num + "=0\r\n")
					if v, ok := zdfTimerMap[attackedPlayer.Num]; ok {
						v.Stop()
						delete(zdfTimerMap, attackedPlayer.Num)
					}
					t := time.AfterFunc(15*time.Second, func() { ReActive(attackedPlayer.Num) })
					zdfTimerMap[attackedPlayer.Num] = t
				} else if attackedPlayer.Team == "blue" {
					attackedPlayer.Dying = true
					SendBroadcastMsg("ZDF" + attackedPlayer.Num + "=4\r\n")
					t := time.AfterFunc(15*time.Second, func() { ChangeToRed(attackedPlayer.Num) })
					zdfTimerMap[attackedPlayer.Num] = t
				}
			}
		}

	}
}
func ChangeToRed(playerNum string) {
	player := GetPlayerByNum(playerNum)
	player.Team = "red"
	delete(zdfTimerMap, playerNum)
	SendBroadcastMsg("ZDF" + playerNum + "=1\r\n")
	redWin := CheckRedWin()
	if redWin {
		//完成逻辑处理
	}
}
func CheckRedWin() bool {
	for _, player := range players {
		if player.Team == "blue" {
			return false
		}
	}
	return true
}
func ReActive(playerNum string) {
	player := GetPlayerByNum(playerNum)
	player.Active = true
	delete(zdfTimerMap, playerNum)
	SendBroadcastMsg("ZDF" + playerNum + "=1\r\n")
}
func ChangeJHDColor(jhdNum string, playerNum string) {
	jhd := GetJHDByNum(jhdNum)
	if jhd.Color == "red" {
		jhd.Color = "blue"
		blueWin := CheckBlueWin()
		if blueWin {
			//完成逻辑处理
		}
	} else if jhd.Color == "blue" {
		jhd.Color = "red"
	}
	delete(jhdTimerMap, playerNum)
}
func CheckBlueWin() bool {
	for _, jhd := range jhds {
		if jhd.Color == "red" {
			return false
		}
	}
	return true
}
func SetTeamByJHD(gloveNum string, team string) {
	for index, p := range players {
		if p.GloveNum == gloveNum {
			players[index].setTeam(team)
			SendBroadcastMsg("ZDF" + p.Num + "=0\r\n")
		}
	}
}
func HasPlayer(players []Player, player Player) bool {
	for _, p := range players {
		if p.Num == player.Num {
			return true
		}
	}
	return false
}

type JHD struct {
	Num    string
	Color  string
	Active bool
}
type Player struct {
	Num      string
	GloveNum string
	Team     string
	Active   bool
	Dying    bool
}

func GetJHDByNum(num string) *JHD {
	for index, jhd := range jhds {
		if jhd.Num == num {
			return &jhds[index]
		}
	}
	return nil
}
func GetPlayerByNum(num string) *Player {
	for index, player := range players {
		if player.Num == num {
			return &players[index]
		}
	}
	return nil
}

func GetPlayerByGloveNum(gloveNum string) (Player,bool) {
	for _, player := range players {
		if player.GloveNum == gloveNum {
			return player,true
		}
	}
	return Player{},false
}
func (player *Player) setTeam(team string) {
	player.Team = team
}

type Msg struct {
	From     string
	Num      string //编号
	GloveNum string
	Status   string //on off
}

var jhds []JHD
var players []Player
var jhdTimerMap map[string]*time.Timer
var zdfTimerMap map[string]*time.Timer
var connPool map[string]net.Conn

func main() {
    connPool = make(map[string] net.Conn)
    go RunHttpApi()
	listen_sock, err := net.Listen("tcp", "127.0.0.1:5567")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer listen_sock.Close()
    log.Println("tcp server running.listen:5567")
	for {
		new_conn, err := listen_sock.Accept()
		if err != nil {
			continue
		}
		clientIP:=strings.Split(new_conn.RemoteAddr().String(),":")[0]
		SaveToConnPool(clientIP,new_conn)
        log.Println("save conn to pool ok")
	}

}
func SaveToConnPool(clientIP string,conn net.Conn){
	log.Println("SaveToConnPool,clientIP=" +clientIP)
	switch clientIP {
	case "192.168.1.72":
		connPool["JHD"]=conn
	case "192.168.1.201":
		connPool["ZDF01"]=conn
	case "192.168.1.202":
		connPool["ZDF02"]=conn
	case "192.168.1.203":
		connPool["ZDF03"]=conn
	case "192.168.1.204":
		connPool["ZDF04"]=conn
	case "192.168.1.205":
		connPool["ZDF05"]=conn
	case "192.168.1.206":
		connPool["ZDF06"]=conn
	case "192.168.1.207":
		connPool["ZDF07"]=conn
	case "192.168.1.208":
		connPool["ZDF08"]=conn
	case "192.168.1.209":
		connPool["ZDF09"]=conn
	case "192.168.1.210":
		connPool["ZDF10"]=conn
	case "192.168.1.211":
		connPool["ZDF11"]=conn
	case "192.168.1.212":
		connPool["ZDF12"]=conn
	case "192.168.1.213":
		connPool["ZDF13"]=conn
	case "192.168.1.214":
		connPool["ZDF14"]=conn
	case "192.168.1.215":
		connPool["ZDF15"]=conn
	case "192.168.1.216":
		connPool["ZDF16"]=conn
	case "192.168.1.217":
		connPool["ZDF17"]=conn
	case "192.168.1.218":
		connPool["ZDF18"]=conn
	case "192.168.1.219":
		connPool["ZDF19"]=conn
	case "192.168.1.220":
		connPool["ZDF20"]=conn
	}

}
func RunHttpApi(){
	http.HandleFunc("/reset", ResetHandler)
	http.HandleFunc("/start", StartHandler)
	//http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Println("http api server running.listen:5568")
	http.ListenAndServe("127.0.0.1:5568", nil)
}
func StartHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into StartHandler")
	for _, player := range players {
		if player.Team == "red" {
			SendBroadcastMsg("ZDF" + player.Num + "=1\r\n")
		} else if player.Team == "blue" {
			SendBroadcastMsg("ZDF" + player.Num + "=2\r\n")
		}
	}
	InitTimerMap()
	InitJHD()
	SendBroadcastMsg("JHD00=1\r\n")
	fmt.Fprint(w, "start ok")
	log.Println("start ok")
}
func InitTimerMap() {
	jhdTimerMap = make(map[string]*time.Timer)
	zdfTimerMap = make(map[string]*time.Timer)
}
func InitJHD() {
	jhds = make([]JHD, 8)
	jhds = append(jhds, JHD{Num: "01", Color: "red", Active: true})
	jhds = append(jhds, JHD{Num: "02", Color: "red", Active: true})
	jhds = append(jhds, JHD{Num: "03", Color: "red", Active: true})
	jhds = append(jhds, JHD{Num: "04", Color: "red", Active: true})
	jhds = append(jhds, JHD{Num: "05", Color: "red", Active: true})
	jhds = append(jhds, JHD{Num: "06", Color: "red", Active: true})
	jhds = append(jhds, JHD{Num: "07", Color: "red", Active: true})
	jhds = append(jhds, JHD{Num: "08", Color: "red", Active: true})
}
func ResetHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into ResetHandler")
	ResetGameStatus()
	fmt.Fprint(w, "reset ok")
	log.Println("reset ok")
}

func ResetGameStatus() {
	connPool=make(map[string] net.Conn)
	for _,c := range connPool {
		c.Write([]byte("ZDF00=03\r\n"))
		go recvConnMsg(c)
	}
}
func SendBroadcastMsg(msg string){
	log.Println("sendBroad:"+msg)
	for _,c := range connPool {
		c.Write([]byte(msg))
	}
}
