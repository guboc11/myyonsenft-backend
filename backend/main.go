package main

import (
	"fmt"
	"log"
	"net/http"

	"guboc11.com/m/contractcall"
)

func mint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mint success"))

}

func getBalanceOf(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	returnValue := contractcall.ViewBalanceOf(address)
	result := fmt.Sprintf("get balance of %s : %s", address, returnValue)
	w.Write([]byte(result))
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/mint", mint)
	http.HandleFunc("/balanceOf", getBalanceOf)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
