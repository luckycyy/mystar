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
	log.Println("remove" + key + " from connPool")
}
func (connPool *ConnPool) sendBroadcast(msg string) {
	for key, conn := range connPool.Pool {
		go BroadcastThread(key, conn, msg)
	}
}
func BroadcastThread(key string, conn net.Conn, msg string) {
	if msg != "hi" {
		log.Println("write msg:" + msg + " to " + key)
	}
	if _, err := conn.Write([]byte(msg)); err != nil {
		log.Println(err.Error())
		myConnPool.remove(key)
		log.Println("write msg:" + msg + " to " + key + " err")
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
var serverIP = "192.168.1.21"
var jdfAddress = "192.168.1.72"
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
	log.Println("tcp server running.listen:5567...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
		myConnPool.add(clientIP, conn)
		log.Println(clientIP + ",added to connPool")
		go recvMsg(clientIP, conn)
	}
}
func runHttpApi() {

	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/setTeam", setTeamHandler)
	http.HandleFunc("/setJHD", setJHDHandler)
	http.HandleFunc("/reset", resetHandler)
	http.HandleFunc("/resetNotRebind", resetNotRebindHandler)
	http.HandleFunc("/start", startHandler)
	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Println("http api server running.listen:5568")
	http.ListenAndServe(serverIP+":5568", nil)
}
func queryHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	myConnPool.sendBroadcast("hi")

	if len(req.Form["v"]) > 0 {
		v := string(req.Form["v"][0])

		if v == "players" {
			jsonData, err := json.Marshal(players)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Fprint(w, string(jsonData))
		} else if v == "onlineZDF" {
			jsonData, err := json.Marshal(myConnPool.Pool)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Fprint(w, string(jsonData))
		}
		return
	}

	log.Println("players:")
	log.Println(players)

	log.Println("jhds:")
	log.Println(jhds)
	log.Println("jhdTimerMap:")
	log.Println(jhdTimerMap)
	log.Println("zdfTimerMap:")
	log.Println(zdfTimerMap)

	log.Println("Pool:")
	log.Println(myConnPool.Pool)

}
func sendHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if (len(req.Form["k"]) > 0 ) && ( len(req.Form["v"]) > 0) && ( len(req.Form["to"]) > 0) {
		to := string(req.Form["to"][0])
		k := string(req.Form["k"][0])
		v := string(req.Form["v"][0])
		msg := k + "=" + v + "\r\n"
		if to == "00" {
			myConnPool.sendBroadcast(msg)
		} else {
			conn, _ := myConnPool.get(zdfAddress[to])
			_, err := conn.Write([]byte(msg))
			log.Println("msg is:" + msg)
			if err != nil {
				log.Println("send msg failed!")
			}
		}

		if v == "1" {
			ChangeToRed(to)
		} else if v == "2" {
			ChangeToBlue(to)
		}
		fmt.Fprint(w, "setok")
	}
}

func setTeamHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if (len(req.Form["player"]) > 0 ) && ( len(req.Form["team"]) > 0) {
		playerNum := string(req.Form["player"][0])
		team := string(req.Form["team"][0])
		player := GetPlayerByNum(playerNum)
		player.Team = team
		fmt.Fprint(w, "setok")
	}
}
func setJHDHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if (len(req.Form["jhdNum"]) > 0 ) && ( len(req.Form["cmd"]) > 0) {
		jhdNum := string(req.Form["jhdNum"][0])
		cmd := string(req.Form["cmd"][0])

		conn, _ := myConnPool.get(jdfAddress)
		_, err := conn.Write([]byte("JHD" + jhdNum + "=" + cmd + "\r\n"))
		if err != nil {
			myConnPool.remove(jdfAddress)
		}
		jhd := GetJHDByNum(jhdNum)
		if cmd == "1" {
			jhd.Color = "red"
		} else if cmd == "2" {
			jhd.Color = "blue"
		}

		log.Println("JHD" + jhdNum + "=" + cmd)
	}
}

//先开机再reset
func resetHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into ResetHandler")
	myConnPool.sendBroadcast("ZDF00=3\r\n")
	players = []Player{}
	InitJHD()
	InitTimerMap()
	InitZDFAddress()
	lockFlag = false
	doorFlag = false
	lockPlayer = ""
	fmt.Fprint(w, "reset ok")
	log.Println("reset ok")
}
func resetNotRebindHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into ResetNotRebindHandler")
	myConnPool.sendBroadcast("ZDF00=0\r\n")
	for index, p := range players {
		players[index].setTeam("blue")
		log.Println("player num:" + p.Num + " setTeam:" + "blue")
	}
	InitJHD()
	InitTimerMap()
	InitZDFAddress()
	lockFlag = false
	doorFlag = false
	lockPlayer = ""
	fmt.Fprint(w, "reset ok")
	log.Println("reset ok")
}
func startHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into StartHandler")
	log.Println("player is:")
	log.Println(players)
	for _, player := range players {
		conn, _ := myConnPool.get(zdfAddress[player.Num])
		if conn == nil {
			continue //如果连接池中找不到玩家的连接 就跳过这个玩家
		}
		log.Println("to:" + conn.RemoteAddr().String())
		if player.Team == "red" {
			_, err := conn.Write([]byte("ZDF" + player.Num + "=1\r\n"))
			if err != nil {
				myConnPool.remove(zdfAddress[player.Num])
			}
			log.Println("ZDF" + player.Num + "=1\r\n")
		} else if player.Team == "blue" {
			_, err := conn.Write([]byte("ZDF" + player.Num + "=2\r\n"))
			if err != nil {
				myConnPool.remove(zdfAddress[player.Num])
			}
			log.Println("ZDF" + player.Num + "=2\r\n")
		}
	}
	myConnPool.sendBroadcast("JHD00=1\r\n")
	fmt.Fprint(w, "start ok")
	log.Println("start ok")
}

func recvMsg(clientIP string, conn net.Conn) {
	buf := make([]byte, 100)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("read msg err")
			myConnPool.remove(clientIP)
			return
		}
		msg := string(buf[0:n])
		log.Println("recv msg:", msg)
		var newMsg Msg
		if err := json.Unmarshal([]byte(msg), &newMsg); err != nil {
			log.Println("recv msg convert to json err")
		}
		if newMsg.From == "bind" && newMsg.Status == "on" {
			p := Player{Num: newMsg.Num, GloveNum: newMsg.GloveNum, Team: "blue", Active: true}
			if HasPlayer(players, p) {
				for k, v := range players {
					if v.Num == p.Num {
						kk := k + 1
						players = append(players[:k], players[kk:]...)
					}
				}
				log.Println("delete old player info")
			}
			players = append(players, p)
			log.Println("append:" + p.Num + " to players")
			conn, _ := myConnPool.get(zdfAddress[p.Num])
			_, err := conn.Write([]byte("ZDF" + p.Num + "=0\r\n"))
			if err != nil {
				myConnPool.remove(zdfAddress[p.Num])
			}
			log.Println("bind ok,players is:")
			log.Println(players)
		} else if newMsg.From == "jhd" {
			if newMsg.Num == "09" || newMsg.Num == "10" {
				SetTeamByJHD(newMsg.GloveNum, "red")
				if !lockFlag {
					lockFlag = true
					lockPlayer = newMsg.GloveNum
					openLockEvent()
				} else if lockFlag && (!doorFlag) && (lockPlayer != newMsg.GloveNum) {
					doorFlag = true
					openDoorEvent()
				}

			} else {
				jhd := GetJHDByNum(newMsg.Num)
				player, exist := GetPlayerByGloveNum(newMsg.GloveNum)
				log.Println("jhd and player status:")
				log.Println(player)
				log.Println(jhd)
				if !exist {
					log.Println("GetPlayerByGloveNum err")
					return
				}
				if newMsg.Status == "on" {
					if player.Dying {
						log.Println("player.dying")
						player.Dying = false
						BlueReLiveEvent()
						if v, ok := zdfTimerMap[player.Num]; ok {
							v.Stop()
							delete(zdfTimerMap, player.Num)
						}
						conn, _ := myConnPool.get(zdfAddress[player.Num])
						_, err := conn.Write([]byte("ZDF" + player.Num + "=2\r\n"))
						if err != nil {
							myConnPool.remove(zdfAddress[player.Num])
						}
					} else if player.Active {
						touchJHDEvent()
						log.Println("player.Active is true:" + player.Team + ",color:" + jhd.Color)
						//开启timer
						t := time.AfterFunc(5*time.Second, func() { ChangeJHDColor(jhd.Num, player.Num) })
						jhdTimerMap[player.Num] = t
					}
				} else if newMsg.Status == "off" {
					if v, ok := jhdTimerMap[player.Num]; ok {
						v.Stop()
						delete(jhdTimerMap, player.Num)
					}
				}
			}
		} else if newMsg.From == "zdf" {
			attackedPlayer := GetPlayerByNum(newMsg.Num)
			attacker, exist := GetPlayerByGloveNum(newMsg.GloveNum)
			if !exist {
				log.Println("GetPlayerByGloveNum err")
				return
			}
			if (attackedPlayer.Team != attacker.Team) && attacker.Active {
				conn, _ := myConnPool.get(zdfAddress[attackedPlayer.Num])
				if attackedPlayer.Team == "red" {
					redAttackedEvent()
					attackedPlayer.Active = false
					_, err := conn.Write([]byte("ZDF" + attackedPlayer.Num + "=0\r\n"))
					if err != nil {
						myConnPool.remove(zdfAddress[attackedPlayer.Num])
					}
					if v, ok := zdfTimerMap[attackedPlayer.Num]; ok {
						v.Stop()
						delete(zdfTimerMap, attackedPlayer.Num)
					}
					t := time.AfterFunc(15*time.Second, func() { ReActive(attackedPlayer.Num) })
					zdfTimerMap[attackedPlayer.Num] = t
				} else if attackedPlayer.Team == "blue" {
					blueAttackedEvent()
					attackedPlayer.Dying = true
					_, err := conn.Write([]byte("ZDF" + attackedPlayer.Num + "=5\r\n"))
					if err != nil {
						myConnPool.remove(zdfAddress[attackedPlayer.Num])
					}
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
	zdfAddress["01"] = "192.168.1.201"
	zdfAddress["02"] = "192.168.1.202"
	zdfAddress["03"] = "192.168.1.203"
	zdfAddress["04"] = "192.168.1.204"
	zdfAddress["05"] = "192.168.1.205"
	zdfAddress["06"] = "192.168.1.206"
	zdfAddress["07"] = "192.168.1.207"
	zdfAddress["08"] = "192.168.1.208"
	zdfAddress["09"] = "192.168.1.209"
	zdfAddress["10"] = "192.168.1.210"
	zdfAddress["11"] = "192.168.1.211"
	zdfAddress["12"] = "192.168.1.212"
	zdfAddress["13"] = "192.168.1.213"
	zdfAddress["14"] = "192.168.1.214"
	zdfAddress["15"] = "192.168.1.215"
	zdfAddress["16"] = "192.168.1.216"
	zdfAddress["17"] = "192.168.1.217"
	zdfAddress["18"] = "192.168.1.218"
	zdfAddress["19"] = "192.168.1.219"
	zdfAddress["20"] = "192.168.1.220"
}

func HasPlayer(players []Player, player Player) bool {
	for _, p := range players {
		if p.Num == player.Num {
			log.Println("has player:" + p.Num)
			return true
		}
	}
	return false
}
func SetTeamByJHD(gloveNum string, team string) {
	for index, p := range players {
		if p.GloveNum == gloveNum {
			players[index].setTeam(team)
			log.Println("player num:" + p.Num + " setTeam:" + team)
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
func GetPlayerByGloveNum(gloveNum string) (*Player, bool) {
	for index, player := range players {
		if player.GloveNum == gloveNum {
			return &players[index], true
		}
	}
	return nil, false
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
	BianYiEvent()
	player := GetPlayerByNum(playerNum)
	player.Team = "red"
	if _, ok := zdfTimerMap[player.Num]; ok {
		delete(zdfTimerMap, playerNum)
	}
	conn, _ := myConnPool.get(zdfAddress[playerNum])
	conn.Write([]byte("ZDF" + playerNum + "=1\r\n"))
	redWin := CheckRedWin()
	if redWin {
		//完成逻辑处理
		log.Println("red win")
		redWinEvent()
		myConnPool.sendBroadcast("JHD00=0\r\n")
		myConnPool.sendBroadcast("ZDF00=6\r\n")
	}
}
func ChangeToBlue(playerNum string) {
	player := GetPlayerByNum(playerNum)
	player.Team = "blue"
	if _, ok := zdfTimerMap[player.Num]; ok {
		delete(zdfTimerMap, playerNum)
	}
	conn, _ := myConnPool.get(zdfAddress[playerNum])
	conn.Write([]byte("ZDF" + playerNum + "=2\r\n"))
}
func ReActive(playerNum string) {
	player := GetPlayerByNum(playerNum)
	player.Active = true
	if _, ok := zdfTimerMap[playerNum]; ok {
		delete(zdfTimerMap, playerNum)
	}
	RedReActiveEvent()
	conn, _ := myConnPool.get(zdfAddress[playerNum])
	_, err := conn.Write([]byte("ZDF" + playerNum + "=1\r\n"))
	if err != nil {
		myConnPool.remove(zdfAddress[playerNum])
	}
}
func ChangeJHDColor(jhdNum string, playerNum string) {
	log.Println("jhdNum:" + jhdNum + " playerNum:" + playerNum)
	jhd := GetJHDByNum(jhdNum)
	player := GetPlayerByNum(playerNum)

	conn, _ := myConnPool.get(jdfAddress)
	if player.Team == "red" {
		JHDChangeRedEvent()
		_, err := conn.Write([]byte("JHD" + jhd.Num + "=1\r\n"))
		if err != nil {
			log.Println("JHD" + jhd.Num + " send msg err")
			myConnPool.remove(jdfAddress)
			return
		}
		jhd.Color = player.Team
		log.Println("JHD" + jhd.Num + "=1")
	} else if player.Team == "blue" {
		JHDChangeBlueEvent()
		_, err := conn.Write([]byte("JHD" + jhd.Num + "=2\r\n"))
		if err != nil {
			log.Println("JHD" + jhd.Num + " send msg err")
			myConnPool.remove(jdfAddress)
			return
		}
		jhd.Color = player.Team
		log.Println("JHD" + jhd.Num + "=2")
		blueWin := CheckBlueWin()
		if blueWin {
			//完成逻辑处理
			log.Println("blue win")
			blueWinEvent()
			myConnPool.sendBroadcast("JHD00=0\r\n")
			myConnPool.sendBroadcast("ZDF00=6\r\n")
		}
	}
	if _, ok := jhdTimerMap[playerNum]; ok {
		delete(jhdTimerMap, playerNum)
	}
}
func CheckBlueWin() bool {
	log.Println("into checkbluewin")
	for _, jhd := range jhds {
		if jhd.Color == "red" {
			return false
		}
	}
	return true
}
func CheckRedWin() bool {
	log.Println("into checkredwin")
	for _, player := range players {
		if player.Team == "blue" {
			return false
		}
	}
	return true
}
func openLockEvent() {
	log.Println("openlock event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%8E%A9%E5%AE%B6%E4%B8%80%E5%88%B7%E5%8D%A1%E5%BC%80%E5%9C%86%E9%97%A8&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}

func openDoorEvent() {
	log.Println("opendoor event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%8E%A9%E5%AE%B6%E4%BA%8C%E5%88%B7%E5%8D%A1%E5%BC%80%E5%8D%B7%E5%B8%98%E9%97%A8&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}

func touchJHDEvent() {
	log.Println("touchJHD event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E5%87%80%E5%8C%96%E7%82%B9%E8%A2%AB%E8%A7%A6%E6%91%B8&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}

func JHDChangeBlueEvent() {
	log.Println("changeBlue event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E5%87%80%E5%8C%96%E7%82%B9%E5%87%80%E5%8C%96%E4%B8%BA%E8%93%9D%E8%89%B2&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}

func JHDChangeRedEvent() {
	log.Println("changeRed event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E5%87%80%E5%8C%96%E7%82%B9%E5%87%80%E5%8C%96%E4%B8%BA%E7%BA%A2%E8%89%B2&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func redAttackedEvent() {
	log.Println("redattackted event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%BA%A2%E6%96%B9%E7%8E%A9%E5%AE%B6%E8%A2%AB%E5%87%BB%E4%B8%AD&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func blueAttackedEvent() {
	log.Println("blueattackted event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E8%93%9D%E6%96%B9%E7%8E%A9%E5%AE%B6%E8%A2%AB%E5%87%BB%E4%B8%AD&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func BianYiEvent() {
	log.Println("bianyi event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E8%93%9D%E6%96%B9%E7%8E%A9%E5%AE%B6%E5%8F%98%E5%BC%82&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func redWinEvent() {
	log.Println("redwin event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%BA%A2%E6%96%B9%E8%83%9C%E5%88%A9&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func blueWinEvent() {
	log.Println("bluewin event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E8%93%9D%E6%96%B9%E8%83%9C%E5%88%A9&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func BlueReLiveEvent() {
	log.Println("BlueReLive event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E8%93%9D%E6%96%B9%E6%B2%BB%E7%96%97%E6%88%90%E5%8A%9F&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func RedReActiveEvent() {
	log.Println("RedReActive event")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.48&group=action_st&st=%E7%BA%A2%E6%96%B9%E9%87%8D%E7%94%9F&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
