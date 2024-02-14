package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func simpleAPICall() {
	data, err := json.Marshal(Payload{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "eth_call",
		Params: []interface{}{
			map[string]interface{}{
				"from": nil,
				"data": "0x70a08231000000000000000000000000f8c847fc824b441f0b4d9641371e6ed3f56cf145",
				"to":   "0xd16d41635c7ece3c13b2c7eae094a92adf41bb2a",
			},
			"latest",
		},
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	resp, err := http.Post(
		"https://rpc.holesky.ethpandaops.io",
		"application/json",
		strings.NewReader(string(data)),
	)
	if err != nil {
		log.Fatalf("Post: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}
	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Println("result", result)
}

func viewBalanceOf() {
	client, err := ethclient.Dial("https://rpc.holesky.ethpandaops.io")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("connect success")

	from := common.HexToAddress("0xF8c847Fc824B441f0b4D9641371e6eD3f56CF145")
	to := common.HexToAddress("0xd16d41635c7ece3c13b2c7eae094a92adf41bb2a")

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
	data, err := abi.Pack("balanceOf", from)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("data : %x\n", data)

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
	// fmt.Println("response :", resp)

	// uint256 반환값 추출
	returnValue, err := abi.Unpack("balanceOf", resp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("반환된 값:", returnValue)

	fmt.Println("called balanceOf")
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
