package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

func CreateEthereumWallet() (address string, privateKey string) {
	privateKeyECDSA, _ := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	publicAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()
	return publicAddress, privateKeyHex
}
