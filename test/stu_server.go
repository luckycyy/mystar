package main
import (
	"net/http"
	"log"
	"database/sql"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/stu/add", addHandler)
	http.HandleFunc("/stu/del", delHandler)
	http.HandleFunc("/stu/list", listHandler)
	http.Handle("/", http.FileServer(http.Dir("/opt/www")))
	log.Print("www server 8888 running.")
	http.ListenAndServe(":8888", nil)
}
func addHandler(w http.ResponseWriter, req *http.Request) {

}
func listHandler(w http.ResponseWriter, req *http.Request) {

}
func delHandler(w http.ResponseWriter, req *http.Request) {

}





func init() {
	// 设置默认数据库
	orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
	//注册定义的model
	orm.RegisterModel(new(User))

	// 创建table
	orm.RunSyncdb("default", false, true)
}