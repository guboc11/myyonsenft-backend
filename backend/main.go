package main

import (
	"fmt"
	"log"
	"net/http"
)

func mint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mint success"))

}

func getBalanceOf(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	returnValue := viewBalanceOf(address)
	result := fmt.Sprintf("get balance of %s : %s", address, returnValue)
	w.Write([]byte(result))
}

func main() {

	http.HandleFunc("/mint", mint)
	http.HandleFunc("/balanceOf", getBalanceOf)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
