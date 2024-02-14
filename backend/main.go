package main

import (
	"encoding/json"
	"log"
	"net/http"

	"guboc11.com/m/api"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/mint", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		api.Mint(address)
		w.Write([]byte("mint success"))
	})

	http.HandleFunc("/balanceOf", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		balance := api.GetBalanceOf(address)

		// JSON으로 변환하여 응답
		w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(result))
		json.NewEncoder(w).Encode(balance)

	})

	api.Mint("0xF8c847Fc824B441f0b4D9641371e6eD3f56CF145")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
