package blockchain

import (
	"mychain/utils"
	"sync"
	"time"
)

var Blockchain []Block
var PendingTransactions []Transaction
var Balances = make(map[string]int)
var NonceMap = make(map[string]int)
var mutex = &sync.Mutex{}

func AddTransaction(tx Transaction) error {
	mutex.Lock()
	defer mutex.Unlock()

	expectedNonce := NonceMap[tx.From]
	if tx.Nonce != expectedNonce {
		return ErrInvalidNonce
	}

	if Balances[tx.From] < tx.Amount+tx.GasFee {
		return ErrInsufficientFunds
	}

	NonceMap[tx.From]++

	PendingTransactions = append(PendingTransactions, tx)
	return nil
}

func MineBlock() Block {
	mutex.Lock()
	defer mutex.Unlock()

	lastBlock := Blockchain[len(Blockchain)-1]

	newBlock := Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: PendingTransactions,
		PrevHash:     lastBlock.Hash,
		Nonce:        0,
	}

	for {

		newBlock.Hash = utils.CalculateHash(
			newBlock.Index,
			newBlock.Timestamp,
			newBlock.PrevHash,
			serializeTransactions(newBlock.Transactions),
			newBlock.Nonce,
		)

		if newBlock.Hash[:4] == "0000" {
			break
		}
		newBlock.Nonce++
	}

	for _, tx := range PendingTransactions {
		Balances[tx.From] -= tx.Amount + tx.GasFee
		Balances[tx.To] += tx.Amount
	}

	PendingTransactions = []Transaction{}

	Blockchain = append(Blockchain, newBlock)
	return newBlock
}
func InitBlockchain() {
	genesis := CreateGenesisBlock()
	Blockchain = append(Blockchain, genesis)
}
