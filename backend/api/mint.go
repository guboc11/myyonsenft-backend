package api

import (
	"context"
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
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func Mint(client *ethclient.Client, address string, nonceQueue chan uint64, txStatusQueue chan TxStatus) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	value := big.NewInt(0)
	gasLimit := uint64(10_000_000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	gasPrice.Mul(gasPrice, big.NewInt(3)) // test 용 : transaction 을 성공시키기 위해

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

	nonce := <-nonceQueue

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	str := fmt.Sprintf("mint done!!! tx sent: %s, nonce :%d\n", signedTx.Hash().Hex(), nonce)
	fmt.Print(str)

	txStatusQueue <- TxStatus{
		Tx:        signedTx.Hash().Hex(),
		Status:    "SUCCESS",
		CreatedAt: time.Now(),
	}
}
