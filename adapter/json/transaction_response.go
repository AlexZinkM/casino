package json

import (
	"casino/boundary/dto"
	"time"
)

type TransactionResponse struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	TransactionType string `json:"transaction_type"`
	Amount          uint   `json:"amount"`
	Timestamp       string `json:"timestamp"`
}

func (r *TransactionResponse) FromDto(dto *dto.TransactionDTO) {
	r.ID = dto.ID
	r.UserID = dto.UserID
	r.TransactionType = dto.TransactionType
	r.Amount = dto.Amount
	r.Timestamp = dto.Timestamp.Format(time.RFC3339)
}

