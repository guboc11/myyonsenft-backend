package api

import (
	"context"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Balance struct {
	// Balance big.Int `json:"balance"`
	Balance string `json:"balance"`
}

func GetBalanceOf(client *ethclient.Client, address string) Balance {
	// 호출할 함수와 인자 데이터 ABI 인코딩
	abi, err := abi.JSON(strings.NewReader(os.Getenv("CONTRACT_ABI")))
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := common.HexToAddress(address)
	toAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	data, err := abi.Pack("balanceOf", fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	msg := ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Gas:   200_000_000,
		Data:  data,
		Value: big.NewInt(0),
	}

	response, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Call Contract error : %v", err)
	}

	// uint256 반환값 추출
	var returnValue interface{}
	err = abi.UnpackIntoInterface(&returnValue, "balanceOf", response)
	if err != nil {
		log.Fatal(err)
	}
	if returnValue == nil {
		log.Fatal("return value is empty")
	}
	balance, ok := returnValue.(*big.Int)
	if !ok {
		log.Println(returnValue)
		log.Fatal("Convert Error")
	}

	return Balance{Balance: balance.String()}
}
