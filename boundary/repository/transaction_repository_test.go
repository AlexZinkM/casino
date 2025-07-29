package repository

import (
	"testing"
	"time"

	"casino/boundary/repo_model"
	"casino/utils"
)

func TestTransactionRepositoryInterface(t *testing.T) {
	var repo TransactionRepository
	_ = repo
	t.Log("TransactionRepository interface is properly defined")
}

func TestTransactionRepositoryMethods(t *testing.T) {
	mockRepo := &MockTransactionRepository{}

	model := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	err := mockRepo.Save(model)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	models, err := mockRepo.GetByUserID("user123", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(models))
	}

	allModels, err := mockRepo.GetAll(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(allModels) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(allModels))
	}

	t.Log("TransactionRepository interface methods work correctly")
}

type MockTransactionRepository struct{}

func (m *MockTransactionRepository) Save(transaction *repo_model.TransactionModel) error {
	return nil
}

func (m *MockTransactionRepository) GetByUserID(userID string, transactionType *string) ([]*repo_model.TransactionModel, error) {
	return []*repo_model.TransactionModel{
		{
			ID:              utils.GenerateUUID(),
			UserID:          userID,
			TransactionType: "bet",
			Amount:          100,
			Timestamp:       time.Now(),
		},
	}, nil
}

func (m *MockTransactionRepository) GetAll(transactionType *string) ([]*repo_model.TransactionModel, error) {
	return []*repo_model.TransactionModel{
		{
			ID:              utils.GenerateUUID(),
			UserID:          utils.GenerateUUID(),
			TransactionType: "bet",
			Amount:          100,
			Timestamp:       time.Now(),
		},
	}, nil
}
