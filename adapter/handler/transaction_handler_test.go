package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	boundarydto "casino/boundary/dto"
	"casino/utils"
)

type MockTransactionUseCase struct {
	processError     error
	userTransactions []*boundarydto.TransactionDTO
	allTransactions  []*boundarydto.TransactionDTO
}

func (m *MockTransactionUseCase) ProcessTransaction(dto *boundarydto.CreateTransactionDTO) error {
	return m.processError
}

func (m *MockTransactionUseCase) GetUserTransactions(userID string, filter *boundarydto.TransactionFilterDTO) ([]*boundarydto.TransactionDTO, error) {
	return m.userTransactions, nil
}

func (m *MockTransactionUseCase) GetAllTransactions(filter *boundarydto.TransactionFilterDTO) ([]*boundarydto.TransactionDTO, error) {
	return m.allTransactions, nil
}

type MockLogger struct{}

func (m *MockLogger) Error(ctx context.Context, errs ...error)     {}
func (m *MockLogger) Info(ctx context.Context, messages ...string) {}

func TestTransactionHandler_GetUserTransactions(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		userTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user123",
				TransactionType: "bet",
				Amount:          100,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	var response struct {
		Transactions []struct {
			ID              string `json:"id"`
			UserID          string `json:"user_id"`
			TransactionType string `json:"transaction_type"`
			Amount          uint   `json:"amount"`
			Timestamp       string `json:"timestamp"`
		} `json:"transactions"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response.Transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(response.Transactions))
	}

	if response.Transactions[0].UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got %s", response.Transactions[0].UserID)
	}
}

func TestTransactionHandler_GetUserTransactions_MissingUserID(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestTransactionHandler_GetUserTransactions_WithFilter(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		userTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user123",
				TransactionType: "bet",
				Amount:          100,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123&transaction_type=bet", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetUserTransactions_EmptyUserID(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestTransactionHandler_GetUserTransactions_WithWinFilter(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		userTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user123",
				TransactionType: "win",
				Amount:          200,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123&transaction_type=win", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetAllTransactions(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		allTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user123",
				TransactionType: "bet",
				Amount:          100,
			},
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user456",
				TransactionType: "win",
				Amount:          200,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	var response struct {
		Transactions []struct {
			ID              string `json:"id"`
			UserID          string `json:"user_id"`
			TransactionType string `json:"transaction_type"`
			Amount          uint   `json:"amount"`
			Timestamp       string `json:"timestamp"`
		} `json:"transactions"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response.Transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(response.Transactions))
	}
}

func TestTransactionHandler_GetAllTransactions_WithFilter(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		allTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user123",
				TransactionType: "win",
				Amount:          200,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions?transaction_type=win", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetAllTransactions_Empty(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		allTransactions: []*boundarydto.TransactionDTO{},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	var response struct {
		Transactions []struct{} `json:"transactions"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response.Transactions) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(response.Transactions))
	}
}

func TestTransactionHandler_GetAllTransactions_WithBetFilter(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		allTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user123",
				TransactionType: "bet",
				Amount:          100,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions?transaction_type=bet", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_Constructor(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}
	mockLogger := &MockLogger{}

	handler := NewTransactionHandler(mockUseCase, mockLogger)

	if handler == nil {
		t.Error("Expected handler to be created")
	}

	if handler.transactionUseCase != mockUseCase {
		t.Error("Expected transaction use case to be set")
	}

	if handler.logger != mockLogger {
		t.Error("Expected logger to be set")
	}
}

func TestTransactionHandler_Constructor_NilUseCase(t *testing.T) {
	mockLogger := &MockLogger{}

	handler := NewTransactionHandler(nil, mockLogger)

	if handler == nil {
		t.Error("Expected handler to be created with nil use case")
	}

	if handler.transactionUseCase != nil {
		t.Error("Expected transaction use case to be nil")
	}
}

func TestTransactionHandler_Constructor_NilLogger(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{}

	handler := NewTransactionHandler(mockUseCase, nil)

	if handler == nil {
		t.Error("Expected handler to be created with nil logger")
	}

	if handler.logger != nil {
		t.Error("Expected logger to be nil")
	}
}
