package blockchain

import (
	"fmt"
	"time"

	"mychain/utils"
)

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
}

func CreateGenesisBlock() Block {
	block := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: nil,
		PrevHash:     "0",
		Nonce:        0,
	}
	block.Hash = calculateBlockHash(block)
	return block
}

func serializeTransactions(transactions []Transaction) string {
	var txStr string
	for _, tx := range transactions {
		txStr += tx.From + tx.To + fmt.Sprintf("%d%d", tx.Amount, tx.GasFee)
	}
	return txStr
}

func calculateBlockHash(block Block) string {
	txData := serializeTransactions(block.Transactions)
	return utils.CalculateHash(
		block.Index,
		block.Timestamp,
		block.PrevHash,
		txData,
		block.Nonce,
	)
}
