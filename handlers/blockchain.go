package handlers

import (
	"encoding/json"
	"net/http"

	"mychain/blockchain"

	"github.com/gorilla/mux"
)

func handleGetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	balance := blockchain.Balances[address]
	json.NewEncoder(w).Encode(map[string]int{"balance": balance})
}

func handleMineBlock(w http.ResponseWriter, r *http.Request) {
	block := blockchain.MineBlock()
	json.NewEncoder(w).Encode(block)
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain.Blockchain)
}
