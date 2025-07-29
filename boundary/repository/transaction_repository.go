package repository

import (
	"casino/boundary/repo_model"
)

type TransactionRepository interface {
	Save(transaction *repo_model.TransactionModel) error
	GetByUserID(userID string, transactionType *string) ([]*repo_model.TransactionModel, error)
	GetAll(transactionType *string) ([]*repo_model.TransactionModel, error)
}
