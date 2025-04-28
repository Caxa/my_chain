package blockchain

import "errors"

var (
	ErrInvalidNonce      = errors.New("invalid nonce")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
