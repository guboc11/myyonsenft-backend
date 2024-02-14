package api

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func Mint(address string) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client, err := ethclient.Dial("https://rpc.holesky.ethpandaops.io")
	if err != nil {
		log.Fatal(err)
	}

	// my privatekey
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddress := common.HexToAddress("0xd16d41635C7ECe3c13B2c7Eae094a92aDF41bB2a")

	// fmt.Println("fromAddress :", fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// var nonce uint64 = 58
	// 호출할 함수와 인자 데이터 ABI 인코딩
	const definition = `[{
    "inputs": [
      {
        "internalType": "address",
        "name": "_user",
        "type": "address"
      }
    ],
    "name": "mint",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  }],`

	abi, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		log.Fatal(err)
	}
	userAddress := common.HexToAddress(address)
	// fmt.Println("userAddress :", userAddress)
	data, err := abi.Pack("mint", userAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("nonce :", nonce)
	value := big.NewInt(0)         // in wei (1 eth)
	gasLimit := uint64(10_000_000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	// test 용 : transaction 을 성공시키기 위해
	gasPrice.Mul(gasPrice, big.NewInt(3))
	fmt.Println("gasPrice :", gasPrice)

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
	fmt.Printf("tx sent: %s, nonce :%d\n", signedTx.Hash().Hex(), nonce)

	fmt.Println("mint done!!!")
}
