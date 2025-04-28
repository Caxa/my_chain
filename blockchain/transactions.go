package blockchain

type Transaction struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    int    `json:"amount"`
	GasFee    int    `json:"gasFee"`
	Nonce     int    `json:"nonce"`
	Signature string `json:"signature"`
}
