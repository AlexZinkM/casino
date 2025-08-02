package usecase

import (
	"testing"
	"time"

	"casino/boundary/dto"
	"casino/boundary/repo_model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Save(transaction *repo_model.TransactionModel) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(id string) (*repo_model.TransactionModel, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repo_model.TransactionModel), args.Error(1)
}

func (m *MockTransactionRepository) GetByUserID(userID string, transactionType *string) ([]*repo_model.TransactionModel, error) {
	args := m.Called(userID, transactionType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repo_model.TransactionModel), args.Error(1)
}

func (m *MockTransactionRepository) GetAll(transactionType *string) ([]*repo_model.TransactionModel, error) {
	args := m.Called(transactionType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repo_model.TransactionModel), args.Error(1)
}

func TestProcessTransaction_Idempotency(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	transactionID := "550e8400-e29b-41d4-a716-446655440001"
	createDto := &dto.CreateTransactionDTO{
		ID:              transactionID,
		UserID:          "550e8400-e29b-41d4-a716-446655440010",
		TransactionType: "bet",
		Amount:          1000,
	}

	expectedEntity := createDto.ToEntity()
	expectedModel := &repo_model.TransactionModel{}
	expectedModel.FromEntity(expectedEntity)

	mockRepo.On("GetByID", transactionID).Return(nil, nil).Once()
	mockRepo.On("Save", mock.AnythingOfType("*repo_model.TransactionModel")).Return(nil).Once()

	err := useCase.ProcessTransaction(createDto)
	assert.NoError(t, err)

	mockRepo.On("GetByID", transactionID).Return(expectedModel, nil).Once()

	err = useCase.ProcessTransaction(createDto)
	assert.Error(t, err)

	mockRepo.AssertNumberOfCalls(t, "Save", 1)
	mockRepo.AssertExpectations(t)
}

func TestProcessTransaction_ErrorHandling(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	createDto := &dto.CreateTransactionDTO{
		ID:              "550e8400-e29b-41d4-a716-446655440001",
		UserID:          "550e8400-e29b-41d4-a716-446655440010",
		TransactionType: "bet",
		Amount:          1000,
	}

	mockRepo.On("GetByID", createDto.ID).Return(nil, assert.AnError).Once()

	err := useCase.ProcessTransaction(createDto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check for existing transaction")

	mockRepo.AssertExpectations(t)
}

func TestProcessTransaction_SaveError(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	createDto := &dto.CreateTransactionDTO{
		ID:              "550e8400-e29b-41d4-a716-446655440001",
		UserID:          "550e8400-e29b-41d4-a716-446655440010",
		TransactionType: "bet",
		Amount:          1000,
	}

	mockRepo.On("GetByID", createDto.ID).Return(nil, nil).Once()
	mockRepo.On("Save", mock.AnythingOfType("*repo_model.TransactionModel")).Return(assert.AnError).Once()

	err := useCase.ProcessTransaction(createDto)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)

	mockRepo.AssertExpectations(t)
}

func TestGetUserTransactions_Success(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	userID := "550e8400-e29b-41d4-a716-446655440010"
	models := []*repo_model.TransactionModel{
		{
			ID:              "550e8400-e29b-41d4-a716-446655440001",
			UserID:          userID,
			TransactionType: "bet",
			Amount:          1000,
			Timestamp:       time.Now(),
		},
		{
			ID:              "550e8400-e29b-41d4-a716-446655440002",
			UserID:          userID,
			TransactionType: "win",
			Amount:          2000,
			Timestamp:       time.Now(),
		},
	}

	mockRepo.On("GetByUserID", userID, (*string)(nil)).Return(models, nil).Once()

	dtos, err := useCase.GetUserTransactions(userID, nil)
	assert.NoError(t, err)
	assert.Len(t, dtos, 2)
	assert.Equal(t, models[0].ID, dtos[0].ID)
	assert.Equal(t, models[1].ID, dtos[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestGetUserTransactions_WithFilter(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	userID := "550e8400-e29b-41d4-a716-446655440010"
	transactionType := "bet"
	filter := &dto.TransactionFilterDTO{
		TransactionType: &transactionType,
	}

	models := []*repo_model.TransactionModel{
		{
			ID:              "550e8400-e29b-41d4-a716-446655440001",
			UserID:          userID,
			TransactionType: "bet",
			Amount:          1000,
			Timestamp:       time.Now(),
		},
	}

	mockRepo.On("GetByUserID", userID, &transactionType).Return(models, nil).Once()

	dtos, err := useCase.GetUserTransactions(userID, filter)
	assert.NoError(t, err)
	assert.Len(t, dtos, 1)
	assert.Equal(t, models[0].ID, dtos[0].ID)
	assert.Equal(t, "bet", dtos[0].TransactionType)

	mockRepo.AssertExpectations(t)
}

func TestGetUserTransactions_EmptyResult(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	userID := "550e8400-e29b-41d4-a716-446655440010"

	mockRepo.On("GetByUserID", userID, (*string)(nil)).Return([]*repo_model.TransactionModel{}, nil).Once()

	dtos, err := useCase.GetUserTransactions(userID, nil)
	assert.NoError(t, err)
	assert.Len(t, dtos, 0)

	mockRepo.AssertExpectations(t)
}

func TestGetUserTransactions_RepositoryError(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	userID := "550e8400-e29b-41d4-a716-446655440010"

	mockRepo.On("GetByUserID", userID, (*string)(nil)).Return(nil, assert.AnError).Once()

	dtos, err := useCase.GetUserTransactions(userID, nil)
	assert.Error(t, err)
	assert.Nil(t, dtos)

	mockRepo.AssertExpectations(t)
}

func TestGetAllTransactions_Success(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	models := []*repo_model.TransactionModel{
		{
			ID:              "550e8400-e29b-41d4-a716-446655440001",
			UserID:          "550e8400-e29b-41d4-a716-446655440010",
			TransactionType: "bet",
			Amount:          1000,
			Timestamp:       time.Now(),
		},
		{
			ID:              "550e8400-e29b-41d4-a716-446655440002",
			UserID:          "550e8400-e29b-41d4-a716-446655440011",
			TransactionType: "win",
			Amount:          2000,
			Timestamp:       time.Now(),
		},
	}

	mockRepo.On("GetAll", (*string)(nil)).Return(models, nil).Once()

	dtos, err := useCase.GetAllTransactions(nil)
	assert.NoError(t, err)
	assert.Len(t, dtos, 2)
	assert.Equal(t, models[0].ID, dtos[0].ID)
	assert.Equal(t, models[1].ID, dtos[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestGetAllTransactions_WithFilter(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	transactionType := "win"
	filter := &dto.TransactionFilterDTO{
		TransactionType: &transactionType,
	}

	models := []*repo_model.TransactionModel{
		{
			ID:              "550e8400-e29b-41d4-a716-446655440002",
			UserID:          "550e8400-e29b-41d4-a716-446655440011",
			TransactionType: "win",
			Amount:          2000,
			Timestamp:       time.Now(),
		},
	}

	mockRepo.On("GetAll", &transactionType).Return(models, nil).Once()

	dtos, err := useCase.GetAllTransactions(filter)
	assert.NoError(t, err)
	assert.Len(t, dtos, 1)
	assert.Equal(t, models[0].ID, dtos[0].ID)
	assert.Equal(t, "win", dtos[0].TransactionType)

	mockRepo.AssertExpectations(t)
}

func TestGetAllTransactions_EmptyResult(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	mockRepo.On("GetAll", (*string)(nil)).Return([]*repo_model.TransactionModel{}, nil).Once()

	dtos, err := useCase.GetAllTransactions(nil)
	assert.NoError(t, err)
	assert.Len(t, dtos, 0)

	mockRepo.AssertExpectations(t)
}

func TestGetAllTransactions_RepositoryError(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	mockRepo.On("GetAll", (*string)(nil)).Return(nil, assert.AnError).Once()

	dtos, err := useCase.GetAllTransactions(nil)
	assert.Error(t, err)
	assert.Nil(t, dtos)

	mockRepo.AssertExpectations(t)
}

func TestNewTransactionUseCaseImpl(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	assert.NotNil(t, useCase)
	assert.Equal(t, mockRepo, useCase.transactionRepo)
}

func TestGetUserTransactions_DataConversion(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	userID := "550e8400-e29b-41d4-a716-446655440010"
	model := &repo_model.TransactionModel{
		ID:              "550e8400-e29b-41d4-a716-446655440001",
		UserID:          userID,
		TransactionType: "bet",
		Amount:          1000,
		Timestamp:       time.Now(),
	}

	models := []*repo_model.TransactionModel{model}

	mockRepo.On("GetByUserID", userID, (*string)(nil)).Return(models, nil).Once()

	dtos, err := useCase.GetUserTransactions(userID, nil)
	assert.NoError(t, err)
	assert.Len(t, dtos, 1)

	dto := dtos[0]
	assert.Equal(t, model.ID, dto.ID)
	assert.Equal(t, model.UserID, dto.UserID)
	assert.Equal(t, model.TransactionType, dto.TransactionType)
	assert.Equal(t, model.Amount, dto.Amount)
	assert.Equal(t, model.Timestamp, dto.Timestamp)

	mockRepo.AssertExpectations(t)
}

func TestGetAllTransactions_DataConversion(t *testing.T) {
	mockRepo := &MockTransactionRepository{}
	useCase := NewTransactionUseCaseImpl(mockRepo)

	model := &repo_model.TransactionModel{
		ID:              "550e8400-e29b-41d4-a716-446655440001",
		UserID:          "550e8400-e29b-41d4-a716-446655440010",
		TransactionType: "win",
		Amount:          2000,
		Timestamp:       time.Now(),
	}

	models := []*repo_model.TransactionModel{model}

	mockRepo.On("GetAll", (*string)(nil)).Return(models, nil).Once()

	dtos, err := useCase.GetAllTransactions(nil)
	assert.NoError(t, err)
	assert.Len(t, dtos, 1)

	dto := dtos[0]
	assert.Equal(t, model.ID, dto.ID)
	assert.Equal(t, model.UserID, dto.UserID)
	assert.Equal(t, model.TransactionType, dto.TransactionType)
	assert.Equal(t, model.Amount, dto.Amount)
	assert.Equal(t, model.Timestamp, dto.Timestamp)

	mockRepo.AssertExpectations(t)
}
