package repository

import (
	"fmt"

	"casino/boundary/repo_model"
	"casino/boundary/repository"

	"gorm.io/gorm"
)

type PostgresTransactionRepository struct {
	db *gorm.DB
}

func NewPostgresTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) Save(transaction *repo_model.TransactionModel) error {
	if transaction == nil {
		return fmt.Errorf("transaction cannot be nil")
	}

	if r.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	if err := r.db.Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	return nil
}

func (r *PostgresTransactionRepository) GetByID(id string) (*repo_model.TransactionModel, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var model repo_model.TransactionModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get transaction by id: %w", err)
	}

	return &model, nil
}

func (r *PostgresTransactionRepository) GetByUserID(userID string, transactionType *string) ([]*repo_model.TransactionModel, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var models []*repo_model.TransactionModel
	query := r.db.Where("user_id = ?", userID)

	if transactionType != nil {
		query = query.Where("transaction_type = ?", *transactionType)
	}

	if err := query.Order("timestamp DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions by user_id: %w", err)
	}

	return models, nil
}

func (r *PostgresTransactionRepository) GetAll(transactionType *string) ([]*repo_model.TransactionModel, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var models []*repo_model.TransactionModel
	query := r.db

	if transactionType != nil {
		query = query.Where("transaction_type = ?", *transactionType)
	}

	if err := query.Order("timestamp DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}

	return models, nil
}
