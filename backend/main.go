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

var door chan uint64

func init() {
	// set log flags
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// ethereum client 생성
	// client, err = ethclient.Dial("https://rpc.holesky.ethpandaops.io")
	// client, err = ethclient.Dial("https://eth-sepolia.api.onfinality.io/public")
	// client, err = ethclient.Dial("https://ethereum-sepolia-rpc.allthatnode.com/6O9ObmNFU0WHoqRLekyfFR9pRb6IqBbi")
	// client, err = ethclient.Dial("https://polygon-mainnet.infura.io/v3/a4d72135ac2a4366892eba1ec2e5aaef")
	client, err = ethclient.Dial("https://sepolia.infura.io/v3/a4d72135ac2a4366892eba1ec2e5aaef")
	if err != nil {
		log.Fatal(err)
	}

	// initialize nonce
	currentNonce, err = client.PendingNonceAt(context.Background(), common.HexToAddress(os.Getenv("SENDER_ADDRESS")))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("currentNonce", currentNonce)

	// make queues
	nonceQueue = make(chan uint64)
	// nonceQueue = make(chan uint64, 10)
	txStatusQueue = make(chan api.TxStatus)

	door = make(chan uint64, 1)
}

func main() {
	log.Println(api.DebuggingNumber, "top of main")
	// go func() {
	// 	for {
	// 		fmt.Println()
	// 		time.Sleep(500 * time.Millisecond)
	// 	}
	// }()
	// /mint endpoint
	go http.HandleFunc("/mint", func(w http.ResponseWriter, r *http.Request) {
		api.DebuggingNumber++
		log.Println(api.DebuggingNumber, "/mint endpoint 시작 점")

		log.Println(api.DebuggingNumber, "get parameters")
		address := r.URL.Query().Get("address")
		tokenUri := r.URL.Query().Get("tokenuri")
		// ethereum address 유효성 검사
		if !isValidEthereumAddress(address) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid address format")
			return
		}

		log.Println(api.DebuggingNumber, "switch method")
		switch r.Method {
		// case http.MethodPost:
		case http.MethodGet:
			log.Println(api.DebuggingNumber, "post method 시작")

			log.Println(api.DebuggingNumber, "Mint() 진입 전")
			go api.Mint(client, address, tokenUri, nonceQueue, txStatusQueue)
			// go api.Mint(client, address, nonceQueue, txStatusQueue, door)
			log.Println(api.DebuggingNumber, "Mint() 함수 실행 완료")

			door <- 1

			log.Println(api.DebuggingNumber, "nonceQueue <- currentNonce 진입 전")
			nonceQueue <- currentNonce
			log.Println(api.DebuggingNumber, "nonceQueue <- currentNonce 실행 완료")

			// log.Println(api.DebuggingNumber, "0.5s time sleep")
			// time.Sleep(500 * time.Millisecond)

			log.Println(api.DebuggingNumber, "currentNonce = <-nonceQueue 진입 전")
			currentNonce = <-nonceQueue
			log.Println(api.DebuggingNumber, "currentNonce = <-nonceQueue 실행 완료")

			var _ uint64 = <-door

			// log.Println(api.DebuggingNumber, "0.5s time sleep")
			// time.Sleep(500 * time.Millisecond)

			// JSON으로 변환하여 응답
			w.Header().Set("Content-Type", "application/json")
			log.Println(api.DebuggingNumber, "txStatus := <-txStatusQueue 진입 전")
			txStatus := <-txStatusQueue
			log.Println(api.DebuggingNumber, "txStatus := <-txStatusQueue 실행 완료")
			json.NewEncoder(w).Encode(txStatus)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		log.Println(api.DebuggingNumber, "/mint endpoint 끝점")
	})

	// /balanceOf endpoint
	// go http.HandleFunc("/balanceOf", func(w http.ResponseWriter, r *http.Request) {
	// 	address := r.URL.Query().Get("address")
	// 	// ethereum address 유효성 검사
	// 	if !isValidEthereumAddress(address) {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode("Invalid address format")
	// 		return
	// 	}
	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		// JSON으로 변환하여 응답
	// 		w.Header().Set("Content-Type", "application/json")
	// 		balance := api.GetBalanceOf(client, address)
	// 		json.NewEncoder(w).Encode(balance)
	// 	default:
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	}

	// })

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

	log.Fatal(http.ListenAndServe(":5555", nil))
}

func isValidEthereumAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
