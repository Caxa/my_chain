package blockchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetBitcoinBalance(address string) (int64, error) {
	resp, err := http.Get(fmt.Sprintf("https://blockstream.info/testnet/api/address/%s", address))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		ChainStats struct {
			FundedTxoSum int64 `json:"funded_txo_sum"`
			SpentTxoSum  int64 `json:"spent_txo_sum"`
		} `json:"chain_stats"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	balance := data.ChainStats.FundedTxoSum - data.ChainStats.SpentTxoSum
	return balance, nil
}

func SendBitcoinTransaction(privKeyHex string, toAddress string, amount float64) (string, error) {
	return "btc-testnet-txhash-mock", nil
}
