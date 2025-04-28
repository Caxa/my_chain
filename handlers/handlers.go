package handlers

import (
	"encoding/json"
	"net/http"

	"mychain/blockchain"
	"mychain/wallet"

	"github.com/gorilla/mux"
)

func RegisterHandlers(r *mux.Router) {
	blockchain.InitBlockchain()

	r.HandleFunc("/wallet/create", handleCreateWallet).Methods("POST")
	r.HandleFunc("/wallet/balance/{address}", handleGetBalance).Methods("GET")
	r.HandleFunc("/transaction/send", handleSendTransaction).Methods("POST")
	r.HandleFunc("/mine", handleMineBlock).Methods("POST")
	r.HandleFunc("/blockchain", handleGetBlockchain).Methods("GET")
}

func handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	newWallet := wallet.CreateWallet()
	blockchain.Balances[newWallet.Address] = 1000
	json.NewEncoder(w).Encode(map[string]string{
		"address":    newWallet.Address,
		"privateKey": newWallet.PrivateKey.D.Text(16),
	})
}

func handleGetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	balance := blockchain.Balances[address]
	json.NewEncoder(w).Encode(map[string]int{"balance": balance})
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

func handleMineBlock(w http.ResponseWriter, r *http.Request) {
	block := blockchain.MineBlock()
	json.NewEncoder(w).Encode(block)
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain.Blockchain)
}
