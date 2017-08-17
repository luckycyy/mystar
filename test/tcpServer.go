package main

import (
	"fmt"
	"net"
	"os"
	"net/http"
	"log"
	"encoding/json"
	"time"
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
						conn.Write([]byte("ZDF" + player.Num + "=2"))
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
					conn.Write([]byte("ZDF" + attackedPlayer.Num + "=0"))
					if v, ok := zdfTimerMap[attackedPlayer.Num]; ok {
						v.Stop()
						delete(zdfTimerMap, attackedPlayer.Num)
					}
					t := time.AfterFunc(15*time.Second, func() { ReActive(attackedPlayer.Num) })
					zdfTimerMap[attackedPlayer.Num] = t
				} else if attackedPlayer.Team == "blue" {
					attackedPlayer.Dying = true
					conn.Write([]byte("ZDF" + attackedPlayer.Num + "=4"))
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
	conn.Write([]byte("ZDF" + playerNum + "=1"))
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
	conn.Write([]byte("ZDF" + playerNum + "=1"))
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
			conn.Write([]byte("ZDF" + p.Num + "=0"))
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
var conn net.Conn
var jhdTimerMap map[string]*time.Timer
var zdfTimerMap map[string]*time.Timer

func main() {

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
        log.Println("conn ok")
		conn = new_conn
		//go recvConnMsg(new_conn)
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
			conn.Write([]byte("ZDF" + player.Num + "=1"))
		} else if player.Team == "blue" {
			conn.Write([]byte("ZDF" + player.Num + "=2"))
		}
	}
	InitTimerMap()
	InitJHD()
	conn.Write([]byte("JHD00=1"))
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
	conn.Write([]byte("ZDF00=03"))
	recvConnMsg(conn)

}
