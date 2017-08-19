package main

import (
    "net"
    "errors"
    "log"
    "strings"
    "net/http"
    "fmt"
    "encoding/json"
    "time"
)

type ConnPool struct {
    Pool map[string]net.Conn
}

func (connPool *ConnPool) add(key string, conn net.Conn) {
    connPool.Pool[key] = conn
}
func (connPool *ConnPool) remove(key string) {
    if conn := connPool.Pool[key]; conn != nil {
        conn.Close()
        log.Println("conn close")
    }
    delete(connPool.Pool, key)
    log.Println("remove" +key+" from connPool")
}
func (connPool *ConnPool) sendBroadcast(msg string) {
    for key, conn := range connPool.Pool {
        go func() {
            if _, err := conn.Write([]byte(msg)); err != nil {
                log.Println(err.Error())
                myConnPool.remove(key)
                log.Println("write msg:"+msg+" to "+key+" err")
            }
        }()
    }
}
func (connPool *ConnPool) get(key string) (net.Conn, error) {
    if conn := connPool.Pool[key]; conn != nil {
        return connPool.Pool[key], nil
    }
    return nil, errors.New("connPool not found key:" + key)
}
func NewConnPool() *ConnPool {
    return &ConnPool{Pool: make(map[string]net.Conn)}
}

type Msg struct {
    From     string
    Num      string //编号
    GloveNum string
    Status   string //on off
}
type Player struct {
    Num      string
    GloveNum string
    Team     string
    Active   bool
    Dying    bool
}
type JHD struct {
    Num    string
    Color  string
    Active bool
}
func (player *Player) setTeam(team string) {
    player.Team = team
}

var myConnPool *ConnPool
var serverIP = "127.0.0.1"
var players []Player
var jhds []JHD
var zdfAddress map[string]string
var jhdTimerMap map[string]*time.Timer
var zdfTimerMap map[string]*time.Timer

var lockFlag bool
var doorFlag bool
var lockPlayer string
func main() {
    myConnPool = NewConnPool()
    go runHttpApi()
    listener, _ := net.Listen("tcp", serverIP+":5567")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println(err.Error())
        }
        clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
        myConnPool.add(clientIP, conn)
        log.Println(clientIP + ",added to connPool")
        go recvMsg(clientIP,conn)
    }
}
func runHttpApi(){
    http.HandleFunc("/reset", resetHandler)
    http.HandleFunc("/start", startHandler)
    //http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
    log.Println("http api server running.listen:5568")
    http.ListenAndServe(serverIP+":5568", nil)
}
func resetHandler(w http.ResponseWriter, req *http.Request) {
    log.Println("into ResetHandler")
    myConnPool.sendBroadcast("ZDF00=03\r\n")
    InitJHD()
    InitTimerMap()
    InitZDFAddress()
    fmt.Fprint(w, "reset ok")
    log.Println("reset ok")
}
func startHandler(w http.ResponseWriter, req *http.Request) {
    log.Println("into StartHandler")
    for _, player := range players {
        conn,_:=myConnPool.get(zdfAddress[player.Num])
        if player.Team == "red" {
            conn.Write([]byte("ZDF" + player.Num + "=1\r\n"))
        } else if player.Team == "blue" {
            conn.Write([]byte("ZDF" + player.Num + "=2\r\n"))
        }
    }
    myConnPool.sendBroadcast("JHD00=1\r\n")
    fmt.Fprint(w, "start ok")
    log.Println("start ok")
}

func recvMsg(clientIP string,conn net.Conn){
    buf := make([]byte, 50)
    for{
        n, err := conn.Read(buf)
        if err != nil {
            log.Println("read msg err")
            myConnPool.remove(clientIP)
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
        }else if newMsg.From == "jhd" {
            if newMsg.Num == "09" || newMsg.Num == "10" {
                SetTeamByJHD(newMsg.GloveNum, "red")
                if !lockFlag{
                    lockFlag=true
                    lockPlayer=newMsg.GloveNum
                    log.Println("open lock")
                }else if lockFlag&&(!doorFlag)&&(lockPlayer!=newMsg.GloveNum) {
                    doorFlag=true
                    log.Println("open door")
                }

            } else {
                jhd := GetJHDByNum(newMsg.Num)
                player,exist:= GetPlayerByGloveNum(newMsg.GloveNum)
                if !exist {
                    log.Println("GetPlayerByGloveNum err")
                }
                if newMsg.Status == "on" {
                    if player.Dying {
                        t := zdfTimerMap[player.Num]
                        t.Stop()
                        delete(zdfTimerMap, player.Num)
                        conn,_:=myConnPool.get(zdfAddress[player.Num])
                        conn.Write([]byte("ZDF" + player.Num + "=2\r\n"))
                    }else if player.Team != jhd.Color {
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
        }else if newMsg.From == "zdf" {
            attackedPlayer := GetPlayerByNum(newMsg.Num)
            attacker ,exist:= GetPlayerByGloveNum(newMsg.GloveNum)
            if !exist {
                log.Println("GetPlayerByGloveNum err")
            }
            if (attackedPlayer.Team != attacker.Team) && attacker.Active {
                conn,_:=myConnPool.get(zdfAddress[attackedPlayer.Num])
                if attackedPlayer.Team == "red" {
                    attackedPlayer.Active = false
                    conn.Write([]byte("ZDF" + attackedPlayer.Num + "=0\r\n"))
                    if v, ok := zdfTimerMap[attackedPlayer.Num]; ok {
                        v.Stop()
                        delete(zdfTimerMap, attackedPlayer.Num)
                    }
                    t := time.AfterFunc(15*time.Second, func() { ReActive(attackedPlayer.Num) })
                    zdfTimerMap[attackedPlayer.Num] = t
                } else if attackedPlayer.Team == "blue" {
                    attackedPlayer.Dying = true
                    conn.Write([]byte("ZDF" + attackedPlayer.Num + "=4\r\n"))
                    t := time.AfterFunc(15*time.Second, func() { ChangeToRed(attackedPlayer.Num) })
                    zdfTimerMap[attackedPlayer.Num] = t
                }
            }
        }
    }
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
func InitTimerMap() {
    jhdTimerMap = make(map[string]*time.Timer)
    zdfTimerMap = make(map[string]*time.Timer)
}
func InitZDFAddress() {
    zdfAddress = make(map[string]string)
    zdfAddress["01"]="192.168.1.201"
    zdfAddress["02"]="192.168.1.202"
    zdfAddress["03"]="192.168.1.203"
    zdfAddress["04"]="192.168.1.204"
    zdfAddress["05"]="192.168.1.205"
    zdfAddress["06"]="192.168.1.206"
    zdfAddress["07"]="192.168.1.207"
    zdfAddress["08"]="192.168.1.208"
    zdfAddress["09"]="192.168.1.209"
    zdfAddress["10"]="192.168.1.210"
    zdfAddress["11"]="192.168.1.211"
    zdfAddress["12"]="192.168.1.212"
    zdfAddress["13"]="192.168.1.213"
    zdfAddress["14"]="192.168.1.214"
    zdfAddress["15"]="192.168.1.215"
    zdfAddress["16"]="192.168.1.216"
    zdfAddress["17"]="192.168.1.217"
    zdfAddress["18"]="192.168.1.218"
    zdfAddress["19"]="192.168.1.219"
    zdfAddress["20"]="192.168.1.220"
}

func HasPlayer(players []Player, player Player) bool {
    for _, p := range players {
        if p.Num == player.Num {
            return true
        }
    }
    return false
}
func SetTeamByJHD(gloveNum string, team string) {
    for index, p := range players {
        if p.GloveNum == gloveNum {
            players[index].setTeam(team)
            //!!SendBroadcastMsg("ZDF" + p.Num + "=0\r\n")
            log.Println("player num:"+p.Num+" setTeam:"+team)
        }
    }
}
func GetJHDByNum(num string) *JHD {
    for index, jhd := range jhds {
        if jhd.Num == num {
            return &jhds[index]
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
func GetPlayerByNum(num string) *Player {
    for index, player := range players {
        if player.Num == num {
            return &players[index]
        }
    }
    return nil
}

func ChangeToRed(playerNum string) {
    player := GetPlayerByNum(playerNum)
    player.Team = "red"
    delete(zdfTimerMap, playerNum)
    conn,_:=myConnPool.get(zdfAddress[playerNum])
    conn.Write([]byte("ZDF" + playerNum + "=1\r\n"))
    redWin := CheckRedWin()
    if redWin {
        //完成逻辑处理
        log.Println("red win")
    }
}

func ReActive(playerNum string) {
    player := GetPlayerByNum(playerNum)
    player.Active = true
    delete(zdfTimerMap, playerNum)
    conn,_:=myConnPool.get(zdfAddress[playerNum])
    conn.Write([]byte("ZDF" + playerNum + "=1\r\n"))
}
func ChangeJHDColor(jhdNum string, playerNum string) {
    jhd := GetJHDByNum(jhdNum)
    player :=GetPlayerByNum(playerNum)
    jhd.Color=player.Team
    if player.Team == "red" {
        blueWin := CheckBlueWin()
        if blueWin {
            //完成逻辑处理
            log.Println("blue win")
        }
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
func CheckRedWin() bool {
    for _, player := range players {
        if player.Team == "blue" {
            return false
        }
    }
    return true
}

