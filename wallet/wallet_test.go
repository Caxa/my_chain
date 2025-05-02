package wallet

import (
	"encoding/hex"
	"testing"
)

func TestCreateEthereumWallet(t *testing.T) {
	address, privateKey := CreateEthereumWallet()

	if len(address) == 0 || len(privateKey) == 0 {
		t.Error("Expected non-empty address and privateKey")
	}
	if _, err := hex.DecodeString(privateKey); err != nil {
		t.Errorf("Private key is not valid hex: %v", err)
	}
}
