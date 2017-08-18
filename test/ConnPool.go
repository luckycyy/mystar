package main

import (
    "net"
    "errors"
    "log"
    "strings"
    "time"
)

type ConnPool struct {
    Pool map[string]net.Conn
}
func(connPool *ConnPool)add(key string,conn net.Conn){
    connPool.Pool[key]=conn
}
func(connPool *ConnPool)remove(key string){
    if conn:=connPool.Pool[key];conn!=nil{
        conn.Close()
    }
    delete(connPool.Pool,key)
}
func(connPool *ConnPool)get(key string)(net.Conn,error){
    if conn:=connPool.Pool[key];conn!=nil{
        return connPool.Pool[key],nil
    }
    return nil,errors.New("connPool not found key:"+key)
}
func NewConnPool() *ConnPool{
    return &ConnPool{Pool:make(map[string]net.Conn)}
}
var connPool *ConnPool
func main(){
    connPool = NewConnPool()
    go sendBroadcast("hi")
    listener,_:=net.Listen("tcp","127.0.0.1:5566")
    for {
        conn,err:=listener.Accept()
        if err!=nil{
            log.Println(err.Error())
        }
        clientIP:=strings.Split(conn.RemoteAddr().String(),":")[0]
        connPool.add(clientIP,conn)
        log.Println(clientIP+",added to connPool")
        log.Println(connPool)
    }
}
func sendBroadcast(msg string){
    for {
        for key,conn:=range connPool.Pool{
            if _,err:=conn.Write([]byte(key+","+msg));err != nil {
                log.Println(err.Error())
                connPool.remove(key)
                log.Println("remove:"+key)
            }
        }
        time.Sleep(3*time.Second)
    }
}