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
			fmt.Println("conn closed")
			return
		}
		msg := string(buf[0:n])
		fmt.Println("recv msg:", msg)
        var newMsg Msg
        if err:=json.Unmarshal([]byte(msg), newMsg);err!=nil{
            log.Println("recv msg convert to json err")
        }

        if newMsg.From=="bind"{
            p:=Player{Num:newMsg.Num,GloveNum:newMsg.GloveNum,Team:"blue"}
            if !HasPlayer(players,p) {
                players = append(players ,p )
            }
        }else if newMsg.From=="jhd" {
            if newMsg.Num=="09"|| newMsg.Num=="10" {
                SetTeamByJHD(newMsg.GloveNum,"red")
            }else {
                jhd:=GetJHDByNum(newMsg.Num)
                player:=GetPlayerByGloveNum(newMsg.GloveNum)
                if newMsg.Status=="on"{
                    if player.Team!=jhd.Color {
                        //开启timer
                        t:=time.AfterFunc(5*time.Second,func(){ChangeJHDColor(jhd.Num,player.Num)})
                        timerMap[player.Num]=t
                    }
                }else if newMsg.Status=="off"{
                    t:=timerMap[player.Num]
                    t.Stop()
                    delete(timerMap,player.Num)
                }

            }

        }else if newMsg.From=="zdf" {

        }

	}
}
func ChangeJHDColor(jhdNum string,playerNum string){
    jhd:=GetJHDByNum(jhdNum)
    if jhd.Color=="red" {
        jhd.Color="blue"
    }else if jhd.Color=="blue" {
        jhd.Color="red"
        CheckFinish()
    }
    delete(timerMap, playerNum)
}
func CheckFinish(){
    for index,jhd :=range jhds{

    }
}
func SetTeamByJHD(gloveNum string,team string){
    for index,p := range players{
        if p.GloveNum ==  gloveNum {
            players[index].setTeam(team)
            conn.Write([]byte("ZDF"+p.Num+"=0"))
        }
    }
}
func HasPlayer(players []Player,player Player) bool{
    for _,p := range players{
        if p.Num ==  player.Num {
            return true
        }
    }
    return false
}
type JHD struct {
	Num      string
	Color    string
	Active   bool //受on off控制
	Disabled bool //失效
}
type Player struct {
	WxId     string
	Num      string
	GloveNum string
	Team     string
	Status   string
}
func GetJHDByNum(num string) *JHD{
    for index,jhd:= range jhds{
        if jhd.Num==num {
            return &jhds[index]
        }
    }
    return nil
}
func GetPlayerByGloveNum(gloveNum string) Player{
    for _,player:= range players{
        if player.GloveNum==gloveNum {
            return player
        }
    }
    return nil
}
func (player *Player)setTeam(team string){
    player.Team=team
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
var timerMap map[string] *time.Timer
func main() {

	http.HandleFunc("/reset", ResetHandler)
    http.HandleFunc("/start", StartHandler)
	//http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Println("server running.")
	http.ListenAndServe(":5568", nil)

	listen_sock, err := net.Listen("tcp", "192.168.1.21:5567")
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
	defer listen_sock.Close()
	fmt.Println("tcp server running v2")
	for {
		new_conn, err := listen_sock.Accept()
		if err != nil {
			continue
		}
		fmt.Println("conn ok")
        conn = new_conn
		//go recvConnMsg(new_conn)
	}

}

func StartHandler(w http.ResponseWriter, req *http.Request) {
    log.Println("into StartHandler")
    for _,player:= range players{
        if player.Team=="red"{
            conn.Write([]byte("ZDF"+player.Num+"=1"))
        }else if player.Team=="blue" {
            conn.Write([]byte("ZDF"+player.Num+"=2"))
        }
    }
    InitJHD()
    conn.Write([]byte("JHD00=1"))
    fmt.Fprint(w, "start ok")
    log.Printf("start ok")
}
func InitJHD(){
    timerMap = make(map[string] *time.Timer)
    jhds = make([]JHD,8)
    jhds=append(jhds,JHD{Num:"01",Color:"red"})
    jhds=append(jhds,JHD{Num:"02",Color:"red"})
    jhds=append(jhds,JHD{Num:"03",Color:"red"})
    jhds=append(jhds,JHD{Num:"04",Color:"red"})
    jhds=append(jhds,JHD{Num:"05",Color:"red"})
    jhds=append(jhds,JHD{Num:"06",Color:"red"})
    jhds=append(jhds,JHD{Num:"07",Color:"red"})
    jhds=append(jhds,JHD{Num:"08",Color:"red"})
}
func ResetHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into ResetHandler")
	ResetGameStatus()
	fmt.Fprint(w, "reset ok")
	log.Printf("reset ok")
}

func ResetGameStatus() {
    conn.Write([]byte("ZDF00=03"))
    recvConnMsg(conn)

}