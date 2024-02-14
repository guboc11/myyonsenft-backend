package main

import (
	"encoding/json"
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
		balance := contractcall.GetBalanceOf(address)

		// JSON으로 변환하여 응답
		w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(result))
		json.NewEncoder(w).Encode(balance)

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
