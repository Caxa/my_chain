package handlers

import (
	"encoding/json"
	"net/http"

	"mychain/wallet"
)

func handleCreateLocalWallet(w http.ResponseWriter, r *http.Request) {
	newWallet := wallet.CreateWallet()
	json.NewEncoder(w).Encode(map[string]string{
		"type":       "local",
		"address":    newWallet.Address,
		"privateKey": newWallet.PrivateKey.D.Text(16),
	})
}

type CreateWalletRequest struct {
	Currency string `json:"currency"`
}

func handleCreateNetworkWallet(w http.ResponseWriter, r *http.Request) {
	var req CreateWalletRequest
	json.NewDecoder(r.Body).Decode(&req)

	switch req.Currency {
	case "BTC":
		address, privateKey := wallet.CreateBitcoinWallet()
		json.NewEncoder(w).Encode(map[string]string{
			"currency":   "BTC-TESTNET",
			"address":    address,
			"privateKey": privateKey,
		})
	case "ETH":
		address, privateKey := wallet.CreateEthereumWallet()
		json.NewEncoder(w).Encode(map[string]string{
			"currency":   "ETH-TESTNET",
			"address":    address,
			"privateKey": privateKey,
		})
	default:
		http.Error(w, "Unsupported currency", http.StatusBadRequest)
	}
}
