package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Address    string
}

func CreateWallet() Wallet {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	address := sha256.Sum256(pubKey)
	return Wallet{
		PrivateKey: privateKey,
		Address:    hex.EncodeToString(address[:]),
	}
}
