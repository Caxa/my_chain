package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/ripemd160"
)

// CreateBitcoinWallet создает биткойн-адрес на Testnet без использования btcutil
func CreateBitcoinWallet() (address string, privateKey string) {
	privKey, _ := btcec.NewPrivateKey() // <<<<< без аргументов!
	pubKey := privKey.PubKey()

	// 1. Хешируем публичный ключ через SHA-256
	shaHash := sha256.New()
	shaHash.Write(pubKey.SerializeCompressed())
	shaHashedPubKey := shaHash.Sum(nil)

	// 2. Далее RIPEMD-160
	ripemdHasher := ripemd160.New()
	ripemdHasher.Write(shaHashedPubKey)
	publicRIPEMD160 := ripemdHasher.Sum(nil)

	// 3. Добавляем версионный байт Testnet (0x6F для Bitcoin Testnet)
	versionedPayload := append([]byte{0x6F}, publicRIPEMD160...)

	// 4. Делаем двойной SHA256 для контрольной суммы
	firstSHA := sha256.Sum256(versionedPayload)
	secondSHA := sha256.Sum256(firstSHA[:])

	// 5. Добавляем 4 байта контрольной суммы к нашему payload
	fullPayload := append(versionedPayload, secondSHA[:4]...)

	// 6. Адрес — это Base58-кодировка полного payload
	address = base58Encode(fullPayload)

	privateKey = hex.EncodeToString(privKey.Serialize())
	return
}

// Простая реализация Base58 без зависимости от btcutil
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func base58Encode(input []byte) string {
	var result []byte
	x := new(big.Int).SetBytes(input)

	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	// Добавляем лидирующие нули
	for _, b := range input {
		if b == 0x00 {
			result = append(result, b58Alphabet[0])
		} else {
			break
		}
	}

	// Переворачиваем результат
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
