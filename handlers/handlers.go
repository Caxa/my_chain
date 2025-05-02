package handlers

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

const INFURA_URL = "https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID"

type WalletResponse struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

type SendTxRequest struct {
	PrivateKey string `json:"privateKey"`
	ToAddress  string `json:"toAddress"`
	Amount     string `json:"amount"` // in ETH
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/wallet/create", CreateWallet).Methods("GET")
	r.HandleFunc("/eth/balance/{address}", GetBalance).Methods("GET")
	r.HandleFunc("/eth/send", SendTransaction).Methods("POST")
}

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		http.Error(w, "Failed to generate key", http.StatusInternalServerError)
		return
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	json.NewEncoder(w).Encode(WalletResponse{
		Address:    address,
		PrivateKey: privateKeyHex,
	})
}

func GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	client, err := ethclient.Dial(INFURA_URL)
	if err != nil {
		http.Error(w, "Failed to connect to Ethereum", http.StatusInternalServerError)
		return
	}

	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	etherValue := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance": etherValue.String(),
	})
}

func SendTransaction(w http.ResponseWriter, r *http.Request) {
	var req SendTxRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	client, err := ethclient.Dial(INFURA_URL)
	if err != nil {
		http.Error(w, "Failed to connect to Ethereum", http.StatusInternalServerError)
		return
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(req.PrivateKey, "0x"))
	if err != nil {
		http.Error(w, "Invalid private key", http.StatusBadRequest)
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		http.Error(w, "Cannot parse public key", http.StatusInternalServerError)
		return
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		http.Error(w, "Failed to get nonce", http.StatusInternalServerError)
		return
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		http.Error(w, "Failed to get gas price", http.StatusInternalServerError)
		return
	}

	value, ok := new(big.Float).SetString(req.Amount)
	if !ok {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	ethToWei := new(big.Float).Mul(value, big.NewFloat(1e18))
	wei := new(big.Int)
	ethToWei.Int(wei)

	toAddress := common.HexToAddress(req.ToAddress)
	tx := types.NewTransaction(nonce, toAddress, wei, uint64(21000), gasPrice, nil)

	chainID := big.NewInt(11155111) // Sepolia
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		http.Error(w, "Failed to sign transaction", http.StatusInternalServerError)
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		http.Error(w, "Failed to send transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"txHash": signedTx.Hash().Hex(),
	})
}
