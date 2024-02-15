package api

import (
	"sort"

	"github.com/ethereum/go-ethereum/ethclient"
)

func GetTxHistory(client *ethclient.Client, address string) []TxStatus {
	_txHistory := []TxStatus{}
	sort.Slice(txHistory[address], func(i, j int) bool {
		// 먼저 CreatedAt 내림차순 정렬
		if txHistory[address][i].CreatedAt.After(txHistory[address][j].CreatedAt) {
			return true
		} else if txHistory[address][i].CreatedAt.Before(txHistory[address][j].CreatedAt) {
			return false
		}

		// CreatedAt이 같을 경우 Tx 이 "" 일 떄 순서가 먼저
		if txHistory[address][i].Tx == "" && txHistory[address][j].Tx != "" {
			return true
		} else if txHistory[address][i].Tx != "" && txHistory[address][j].Tx == "" {
			return false
		}

		// 둘 다 Tx != "" 일 경우 SUCCESS 순서가 나중
		if txHistory[address][i].Status == "SUCCESS" {
			return false
		} else {
			return true

		}
	})

	// 100개 까지 출력
	for i, txStatus := range txHistory[address] {
		if i >= 100 {
			break
		}
		_txHistory = append(_txHistory, txStatus)
		// fmt.Println(i, txStatus)
	}

	return _txHistory
}
