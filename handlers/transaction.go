package handlers

import (
	"encoding/json"
	"net/http"

	"mychain/blockchain"

	"github.com/gorilla/mux"
)

type SendTxRequest struct {
	Currency       string  `json:"currency"`
	FromPrivateKey string  `json:"fromPrivateKey"`
	ToAddress      string  `json:"toAddress"`
	Amount         float64 `json:"amount"`
}

func handleSendTransaction(w http.ResponseWriter, r *http.Request) {
	var tx blockchain.Transaction
	json.NewDecoder(r.Body).Decode(&tx)

	err := blockchain.AddTransaction(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "transaction added"})
}

func handleSendNetworkTransaction(w http.ResponseWriter, r *http.Request) {
	var req SendTxRequest
	json.NewDecoder(r.Body).Decode(&req)

	var txHash string
	var err error

	switch req.Currency {
	case "BTC":
		txHash, err = blockchain.SendBitcoinTransaction(req.FromPrivateKey, req.ToAddress, req.Amount)
	case "ETH":
		txHash, err = blockchain.SendEthereumTransaction(req.FromPrivateKey, req.ToAddress, req.Amount)
	default:
		http.Error(w, "Unsupported currency", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"currency": req.Currency,
		"txHash":   txHash,
	})
}

func RegisterHandlers(r *mux.Router) {
	blockchain.InitBlockchain()

	r.HandleFunc("/wallet/create", handleCreateLocalWallet).Methods("POST")
	r.HandleFunc("/wallet/network/create", handleCreateNetworkWallet).Methods("POST")
	r.HandleFunc("/wallet/balance/{address}", handleGetBalance).Methods("GET")
	r.HandleFunc("/transaction/send", handleSendTransaction).Methods("POST")
	r.HandleFunc("/transaction/network/send", handleSendNetworkTransaction).Methods("POST")
	r.HandleFunc("/mine", handleMineBlock).Methods("POST")
	r.HandleFunc("/blockchain", handleGetBlockchain).Methods("GET")
}
