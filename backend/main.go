package main

import (
	"fmt"
	"log"
	"net/http"

	"guboc11.com/m/contractcall"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/mint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("mint success"))
	})

	http.HandleFunc("/balanceOf", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		returnValue := contractcall.GetBalanceOf(address)
		result := fmt.Sprintf("get balance of %s : %s", address, returnValue)
		w.Write([]byte(result))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
