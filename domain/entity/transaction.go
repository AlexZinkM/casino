package entity

import (
	"time"
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


