package blockchain

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"mychain/utils"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Transaction struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"` // В ETH
	GasFee    float64 `json:"gasFee"` // В Gwei
	Nonce     uint64  `json:"nonce"`
	Signature string  `json:"signature"`
}

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
}

var (
	Chain    []Block
	Pending  []Transaction
	Balances = make(map[string]float64)
	NonceMap = make(map[string]uint64)
	mutex    = &sync.Mutex{}
)

func InitBlockchain() {
	genesis := Block{Index: 0, Timestamp: time.Now().Unix(), PrevHash: "0", Nonce: 0}
	genesis.Hash = utils.CalculateHash(genesis.Index, genesis.Timestamp, genesis.PrevHash, "", genesis.Nonce)
	Chain = append(Chain, genesis)
}

func AddTransaction(tx Transaction) error {
	mutex.Lock()
	defer mutex.Unlock()
	if tx.Nonce != NonceMap[tx.From] {
		return ErrInvalidNonce
	}
	if Balances[tx.From] < tx.Amount+(tx.GasFee/1e9) { // Преобразуем Gwei в ETH
		return ErrInsufficientFunds
	}
	NonceMap[tx.From]++
	Pending = append(Pending, tx)
	return nil
}

func MineBlock() Block {
	mutex.Lock()
	defer mutex.Unlock()

	last := Chain[len(Chain)-1]
	block := Block{
		Index:        last.Index + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: Pending,
		PrevHash:     last.Hash,
		Nonce:        0,
	}

	for {
		hash := utils.CalculateHash(block.Index, block.Timestamp, block.PrevHash, serialize(block.Transactions), block.Nonce)
		if hash[:4] == "0000" {
			block.Hash = hash
			break
		}
		block.Nonce++
	}

	for _, tx := range Pending {
		Balances[tx.From] -= tx.Amount + (tx.GasFee / 1e9)
		Balances[tx.To] += tx.Amount
	}

	Pending = nil
	Chain = append(Chain, block)
	return block
}

func serialize(txs []Transaction) string {
	var out string
	for _, tx := range txs {
		out += tx.From + tx.To + utils.Itoa(int(tx.Amount*1e18+tx.GasFee*1e9))
	}
	return out
}

// SendEthereumTransaction - отправка транзакции в Ethereum testnet (например, Goerli)
func SendEthereumTransaction(client *ethclient.Client, fromPrivateKey string, toAddress string, amount float64, gasFeeGwei float64, nonce uint64) (string, error) {
	privateKey, err := crypto.HexToECDSA(fromPrivateKey)
	if err != nil {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Проверка баланса отправителя (по желанию)
	_, err = client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return "", err
	}

	// Преобразуем amount ETH в Wei
	amountWei := new(big.Int)
	amountWei.SetString(fmt.Sprintf("%.0f", amount*math.Pow10(18)), 10)

	gasPriceWei := new(big.Int)
	gasPriceWei.SetString(fmt.Sprintf("%.0f", gasFeeGwei*math.Pow10(9)), 10)

	gasLimit := uint64(21000)
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), amountWei, gasLimit, gasPriceWei, nil)

	chainID := big.NewInt(11155111)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
