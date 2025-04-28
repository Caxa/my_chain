package utils

import (
	"crypto/sha256"
	"fmt"
)

func CalculateHash(index int, timestamp int64, prevHash string, txData string, nonce int) string {
	data := fmt.Sprintf("%d%d%s%s%d", index, timestamp, prevHash, txData, nonce)
	hash := sha256.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func SerializeTransactions(txData string) string {
	return txData
}
