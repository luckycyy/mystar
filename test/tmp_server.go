package main
import (
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/dati", DatiHandler)
	http.Handle("/", http.FileServer(http.Dir("c:\\www")))
	log.Print("www server running.")
	http.ListenAndServe(":7777", nil)
}
func DatiHandler(w http.ResponseWriter, req *http.Request) {



}