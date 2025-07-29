package entity

import (
	"time"

	"casino/utils"
)

type TransactionType string

const (
	TransactionTypeBet TransactionType = "bet"
	TransactionTypeWin TransactionType = "win"
)

type Transaction struct {
	ID              string
	UserID          string
	TransactionType TransactionType
	Amount          uint
	Timestamp       time.Time
}

func NewTransaction(userID string, transactionType TransactionType, amount uint) *Transaction {
	return &Transaction{
		ID:              utils.GenerateUUID(),
		UserID:          userID,
		TransactionType: transactionType,
		Amount:          amount,
		Timestamp:       time.Now(),
	}
}
