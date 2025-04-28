package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var infuraURL = "https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID"

func GetEthereumBalance(address string) (int64, error) {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return 0, err
	}

	return balance.Int64(), nil
}

func SendEthereumTransaction(privKeyHex string, toAddress string, amount float64) (string, error) {
	return "eth-testnet-txhash-mock", nil
}
