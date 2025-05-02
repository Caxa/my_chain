package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/feb25711f0e44b0fac2c578789fd418c")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("83dab171ec53b5a31c6ddc6f5464e7af0d570290f7c1ec0bc516e91d8daaca0a")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0)
	gasLimit := uint64(21000)
	gasPriceGwei := 50.0

	gasPrice := new(big.Int)
	gasPrice.SetString(fmt.Sprintf("%.0f", gasPriceGwei*math.Pow10(9)), 10)

	toAddress := common.HexToAddress("0xRecipientAddressHere")
	var data []byte

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID := big.NewInt(11155111) //Chain ID для Sepolia: 11155111
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🚀 Transaction sent: %s\n", signedTx.Hash().Hex())
}
