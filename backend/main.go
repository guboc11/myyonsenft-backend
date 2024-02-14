package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Payload struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

func viewBalanceOf() {
	// ethereum client 생성
	client, err := ethclient.Dial("https://rpc.holesky.ethpandaops.io")
	if err != nil {
		log.Fatal(err)
	}

	// 호출할 함수와 인자 데이터 ABI 인코딩
	const definition = `
	[{"inputs":[{"internalType":"address","name":"account","type":"address"}],
	"name":"balanceOf",
	"outputs":[{"internalType":"uint256","name":"","type":"uint256"}],
	"stateMutability":"view",
	"type":"function"}]
	`

	abi, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		log.Fatal(err)
	}

	from := common.HexToAddress("0xF8c847Fc824B441f0b4D9641371e6eD3f56CF145")
	to := common.HexToAddress("0xd16d41635c7ece3c13b2c7eae094a92adf41bb2a")

	data, err := abi.Pack("balanceOf", from)
	if err != nil {
		log.Fatal(err)
	}

	msg := ethereum.CallMsg{
		From:  from,
		To:    &to,
		Gas:   200_000_000,
		Data:  data,
		Value: big.NewInt(0),
	}

	resp, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Call Contract error : %v", err)
	}

	// uint256 반환값 추출
	returnValue, err := abi.Unpack("balanceOf", resp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("result", returnValue)
}

func mint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mint success"))

}

func getBalanceOf(w http.ResponseWriter, r *http.Request) {
	viewBalanceOf()
	w.Write([]byte("get balance of"))
}

func main() {

	http.HandleFunc("/mint", mint)
	http.HandleFunc("/balanceOf", getBalanceOf)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
