package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxStatus struct {
	Tx        string    `json:"tx"`
	Nonce     uint64    `json:"nonce"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

var txHistory map[string][]TxStatus

func init() {
	txHistory = make(map[string][]TxStatus)

	// txHistory 변수에 기존 txHistory.json 파일 내용 불러오기
	existingData, err := os.ReadFile("txHistory.json")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	if len(existingData) > 0 {
		if err := json.Unmarshal(existingData, &txHistory); err != nil {
			log.Fatalf("failed to unmarshal existing JSON: %v", err)
		}
	} else {
		// txHistory.json 파일이 없으면 빈 값을 넣어 생성
		writeTxHistory()
	}
}

func writeTxHistory() {
	newJSONData, err := json.MarshalIndent(txHistory, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal JSON: %v", err)
	}
	err = os.WriteFile("txHistory.json", newJSONData, 0644)
	if err != nil {
		log.Fatalf("failed to write JSON file: %v", err)
	}
}

func Mint(client *ethclient.Client, address string, nonceQueue chan uint64, txStatusQueue chan TxStatus) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// 호출할 contract 주소
	toAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	// 호출할 함수와 인자 데이터 ABI 인코딩
	abi, err := abi.JSON(strings.NewReader(os.Getenv("CONTRACT_ABI")))
	if err != nil {
		log.Fatal(err)
	}
	userAddress := common.HexToAddress(address)
	data, err := abi.Pack("mint", userAddress)
	if err != nil {
		log.Fatal(err)
	}

	// transaction 생성 위한 값들 생성
	value := big.NewInt(0)
	gasLimit := uint64(10_000_000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasPrice.Mul(gasPrice, big.NewInt(4)) // test 용 : transaction 을 성공시키기 위해

	// nonceQueue에서 nonce를 불러 와서 1 추가한 값을 다시 보냄
	nonce := <-nonceQueue
	nonceQueue <- nonce + 1
	txStatus := TxStatus{
		Tx:        "",
		Nonce:     nonce,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}
	txHistory[address] = append(txHistory[address], txStatus)

	// transaction 생성 및 서명
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	txStatus.Tx = signedTx.Hash().Hex()
	txHistory[address] = append(txHistory[address], txStatus)

	// signed transaction 전송
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// error 없으면 전송 성공
	txStatus.Status = "SUCCESS"
	txHistory[address] = append(txHistory[address], txStatus)
	// txHistory 업데이트
	writeTxHistory()

	// 결과 출력 및 전송
	str := fmt.Sprintf("mint done!!! tx sent: %s, nonce :%d\n", signedTx.Hash().Hex(), nonce)
	fmt.Print(str)

	txStatusQueue <- txStatus
}
