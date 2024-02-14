package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"guboc11.com/m/api"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// ethereum client 생성
	client, err := ethclient.Dial("https://rpc.holesky.ethpandaops.io")
	if err != nil {
		log.Fatal(err)
	}

	// set log flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// /mint endpoint
	http.HandleFunc("/mint", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		api.Mint(client, address)
		w.Write([]byte("mint success"))
	})

	// /balanceOf endpoint
	http.HandleFunc("/balanceOf", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		balance := api.GetBalanceOf(client, address)

		// JSON으로 변환하여 응답
		w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(result))
		json.NewEncoder(w).Encode(balance)

	})

	// api.Mint("0xF8c847Fc824B441f0b4D9641371e6eD3f56CF145")

	// for i := 0; i < 10; i++ {
	// 	func() {
	// 		api.Mint("0xF8c847Fc824B441f0b4D9641371e6eD3f56CF145")
	// 	}()
	// }

	log.Fatal(http.ListenAndServe(":8080", nil))
}
