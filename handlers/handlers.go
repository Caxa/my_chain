package handlers

import (
	"encoding/json"
	"net/http"

	"mychain/blockchain"
	"mychain/wallet"
)

func handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	newWallet := wallet.CreateWallet()
	blockchain.Balances[newWallet.Address] = 1000
	json.NewEncoder(w).Encode(map[string]string{
		"address":    newWallet.Address,
		"privateKey": newWallet.PrivateKey.D.Text(16),
	})
}
