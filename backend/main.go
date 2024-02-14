package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

func mint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mint success"))

}

func getBalanceOf(w http.ResponseWriter, r *http.Request) {
	simpleAPICall()
	w.Write([]byte("get balance of"))
}

func main() {
	fmt.Println("hello Planetarium")

	http.HandleFunc("/mint", mint)
	http.HandleFunc("/balanceOf", getBalanceOf)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
