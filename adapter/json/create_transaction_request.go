package json

import (
	"casino/boundary/dto"
)

type CreateTransactionRequest struct {
	UserID          string `json:"user_id"`
	TransactionType string `json:"transaction_type"`
	Amount          uint   `json:"amount"`
}

func (r *CreateTransactionRequest) ToDto() *dto.CreateTransactionDTO {
	return &dto.CreateTransactionDTO{
		UserID:          r.UserID,
		TransactionType: r.TransactionType,
		Amount:          r.Amount,
	}
}

func (r *CreateTransactionRequest) FromDto(dto *dto.CreateTransactionDTO) {
	r.UserID = dto.UserID
	r.TransactionType = dto.TransactionType
	r.Amount = dto.Amount
}
