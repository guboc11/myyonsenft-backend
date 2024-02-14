package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"guboc11.com/m/api"
)

var currentNonce uint64
var nonceQueue chan uint64
var txStatusQueue chan api.TxStatus
var client *ethclient.Client

func init() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// ethereum client 생성
	client, err = ethclient.Dial("https://rpc.holesky.ethpandaops.io")
	if err != nil {
		log.Fatal(err)
	}

	// set log flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// initialize nonce
	currentNonce, err = client.PendingNonceAt(context.Background(), common.HexToAddress(os.Getenv("DELIGATOR_ADDRESS")))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("currentNonce", currentNonce)

	// make queue
	nonceQueue = make(chan uint64)
	txStatusQueue = make(chan api.TxStatus)
}

func main() {
	// /mint endpoint
	go http.HandleFunc("/mint", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		go api.Mint(client, address, nonceQueue, txStatusQueue)

		nonceQueue <- currentNonce
		currentNonce++

		// JSON으로 변환하여 응답
		txStatus := <-txStatusQueue
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(txStatus)
	})

	// /balanceOf endpoint
	go http.HandleFunc("/balanceOf", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		balance := api.GetBalanceOf(client, address)

		// JSON으로 변환하여 응답
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(balance)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
