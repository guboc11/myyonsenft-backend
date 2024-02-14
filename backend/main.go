package main

import (
	"fmt"
	"log"
	"net/http"
)

type Payload struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

func mint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mint success"))

}

func balanceOf(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get balance of"))
}

func main() {
	fmt.Println("hello Planetarium")

	http.HandleFunc("/mint", mint)
	http.HandleFunc("/balanceOf", balanceOf)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
