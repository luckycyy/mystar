package main
import (
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/bohao", BohaoHandler)

	http.Handle("/", http.FileServer(http.Dir("/opt/project/go_server/www")))
	log.Print("server running.")
	http.ListenAndServe(":5572", nil)
}

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
