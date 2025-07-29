package usecase

import (
	"casino/boundary/dto"
)

type TransactionUseCase interface {
	ProcessTransaction(dto *dto.CreateTransactionDTO) error
	GetUserTransactions(userID string, filter *dto.TransactionFilterDTO) ([]*dto.TransactionDTO, error)
	GetAllTransactions(filter *dto.TransactionFilterDTO) ([]*dto.TransactionDTO, error)
}
