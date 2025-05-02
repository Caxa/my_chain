package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

func CreateEthereumWallet() (address, privateKey string) {
	privKey, _ := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	privKeyBytes := crypto.FromECDSA(privKey)
	address = crypto.PubkeyToAddress(privKey.PublicKey).Hex()
	privateKey = hex.EncodeToString(privKeyBytes)
	return
}
