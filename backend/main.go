package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

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
	// set log flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

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

	// initialize nonce
	currentNonce, err = client.PendingNonceAt(context.Background(), common.HexToAddress(os.Getenv("DELIGATOR_ADDRESS")))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("currentNonce", currentNonce)

	// make queues
	nonceQueue = make(chan uint64)
	txStatusQueue = make(chan api.TxStatus)
}

func main() {
	// /mint endpoint
	go http.HandleFunc("/mint", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		// ethereum address 유효성 검사
		if !isValidEthereumAddress(address) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid address format")
			return
		}

		switch r.Method {
		case http.MethodPost:
			go api.Mint(client, address, nonceQueue, txStatusQueue)
			nonceQueue <- currentNonce
			currentNonce = <-nonceQueue

			// JSON으로 변환하여 응답
			w.Header().Set("Content-Type", "application/json")
			txStatus := <-txStatusQueue
			json.NewEncoder(w).Encode(txStatus)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})

	// /balanceOf endpoint
	go http.HandleFunc("/balanceOf", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		// ethereum address 유효성 검사
		if !isValidEthereumAddress(address) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid address format")
			return
		}
		switch r.Method {
		case http.MethodGet:
			// JSON으로 변환하여 응답
			w.Header().Set("Content-Type", "application/json")
			balance := api.GetBalanceOf(client, address)
			json.NewEncoder(w).Encode(balance)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})

	// /history endpoint
	go http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		// ethereum address 유효성 검사
		if !isValidEthereumAddress(address) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid address format")
			return
		}
		switch r.Method {
		case http.MethodGet:
			// JSON으로 변환하여 응답
			w.Header().Set("Content-Type", "application/json")
			txHistory := api.GetTxHistory(client, address)
			json.NewEncoder(w).Encode(txHistory)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func isValidEthereumAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
