package contractcall

import (
	"context"
	"log"
	"math/big"
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

func GetBalanceOf(address string) Balance {
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

	from := common.HexToAddress(address)
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
	bal, ok := returnValue.(*big.Int)
	if !ok {
		log.Println(returnValue)
		log.Fatal("Convert Error")
	}
	// bal, ok := returnValue.(string)

	balance := Balance{Balance: bal.String()}
	return balance
}
