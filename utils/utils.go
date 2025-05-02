package utils

import (
	"crypto/sha256"
	"fmt"
)

func SerializeTransactions(txData string) string {
	return txData
}
func CalculateHash(index int, timestamp int64, prevHash, txData string, nonce int) string {
	data := fmt.Sprintf("%d%d%s%s%d", index, timestamp, prevHash, txData, nonce)
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
