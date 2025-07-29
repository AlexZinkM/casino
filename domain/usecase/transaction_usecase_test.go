package usecase

import (
	"testing"
	"time"

	"casino/boundary/dto"
	"casino/boundary/repo_model"
	"casino/utils"
)

type MockTransactionRepository struct {
	models    []*repo_model.TransactionModel
	saveError error
	getError  error
}

func (m *MockTransactionRepository) Save(transaction *repo_model.TransactionModel) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.models = append(m.models, transaction)
	return nil
}

func (m *MockTransactionRepository) GetByUserID(userID string, transactionType *string) ([]*repo_model.TransactionModel, error) {
	if m.getError != nil {
		return nil, m.getError
	}

	var filtered []*repo_model.TransactionModel
	for _, t := range m.models {
		if t.UserID == userID {
			if transactionType == nil || t.TransactionType == *transactionType {
				filtered = append(filtered, t)
			}
		}
	}
	return filtered, nil
}

func (m *MockTransactionRepository) GetAll(transactionType *string) ([]*repo_model.TransactionModel, error) {
	if m.getError != nil {
		return nil, m.getError
	}

	var filtered []*repo_model.TransactionModel
	for _, t := range m.models {
		if transactionType == nil || t.TransactionType == *transactionType {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

func TestProcessTransaction(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	createDto := &dto.CreateTransactionDTO{
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
	}

	err := useCase.ProcessTransaction(createDto)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockRepo.models) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(mockRepo.models))
	}

	model := mockRepo.models[0]
	if model.UserID != createDto.UserID {
		t.Errorf("Expected UserID %s, got %s", createDto.UserID, model.UserID)
	}

	if model.TransactionType != createDto.TransactionType {
		t.Errorf("Expected TransactionType %s, got %s", createDto.TransactionType, model.TransactionType)
	}

	if model.Amount != createDto.Amount {
		t.Errorf("Expected Amount %d, got %d", createDto.Amount, model.Amount)
	}
}

func TestGetUserTransactions(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	userID := utils.GenerateUUID()
	model1 := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          userID,
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}
	model2 := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          userID,
		TransactionType: "win",
		Amount:          200,
		Timestamp:       time.Now(),
	}
	mockRepo.models = []*repo_model.TransactionModel{model1, model2}

	dtos, err := useCase.GetUserTransactions(userID, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(dtos) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(dtos))
	}

	filter := &dto.TransactionFilterDTO{
		TransactionType: stringPtr("bet"),
	}

	dtos, err = useCase.GetUserTransactions(userID, filter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(dtos) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(dtos))
	}

	if dtos[0].TransactionType != "bet" {
		t.Errorf("Expected transaction type 'bet', got %s", dtos[0].TransactionType)
	}
}

func TestGetAllTransactions(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	model1 := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}
	model2 := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "win",
		Amount:          200,
		Timestamp:       time.Now(),
	}
	mockRepo.models = []*repo_model.TransactionModel{model1, model2}

	dtos, err := useCase.GetAllTransactions(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(dtos) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(dtos))
	}

	filter := &dto.TransactionFilterDTO{
		TransactionType: stringPtr("win"),
	}

	dtos, err = useCase.GetAllTransactions(filter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(dtos) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(dtos))
	}

	if dtos[0].TransactionType != "win" {
		t.Errorf("Expected transaction type 'win', got %s", dtos[0].TransactionType)
	}
}

func stringPtr(s string) *string {
	return &s
}
