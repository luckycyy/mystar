package main
import (
	"net/http"
	"log"
	"strconv"
)

func main() {
	http.HandleFunc("/bohao", BohaoHandler)

	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Print("server running.")
	http.ListenAndServe(":5572", nil)
}
var bohao_start bool
func BohaoHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("into BohaoHandler")
	req.ParseForm()
	if len(req.Form["v"]) > 0 {
		val:=string(req.Form["v"][0])
		if val=="ok" {
			log.Println("  BohaoHandler v is ok")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%8B%A8%E9%80%9A110%E7%94%B5%E8%AF%9D&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if val=="audio"{
			log.Println(" BohaoHandler v is audio")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%89%8B%E6%9C%BA%E6%8C%89%E9%94%AE%E9%9F%B3%E6%95%88&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}else if val=="bohao_start"{
			log.Println(" BohaoHandler v is bohao_start")
			bohao_start=true
		}else if val=="bohao_end"{
			log.Println(" BohaoHandler v is bohao_end")
			bohao_start=false
		}else if val=="bohao_status"{
			log.Println(" BohaoHandler v is bohao_status")
			w.Write([]byte(strconv.FormatBool(bohao_start)))
		}else{
			log.Println(" BohaoHandler v is not ok")
			resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.55&group=action_st&st=%E6%9F%AF%E5%8D%97-%E6%89%8B%E6%9C%BA%E6%8B%A8%E6%89%93%E6%97%A0%E6%95%88&user_action=true")
			if err != nil {
				print(err)
			}
			resp.Body.Close()
		}


	}

}
