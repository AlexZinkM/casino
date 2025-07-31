package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	getUserError     error
	getAllError      error
}

func (m *MockTransactionUseCase) ProcessTransaction(dto *boundarydto.CreateTransactionDTO) error {
	return m.processError
}

func (m *MockTransactionUseCase) GetUserTransactions(userID string, filter *boundarydto.TransactionFilterDTO) ([]*boundarydto.TransactionDTO, error) {
	if m.getUserError != nil {
		return nil, m.getUserError
	}
	return m.userTransactions, nil
}

func (m *MockTransactionUseCase) GetAllTransactions(filter *boundarydto.TransactionFilterDTO) ([]*boundarydto.TransactionDTO, error) {
	if m.getAllError != nil {
		return nil, m.getAllError
	}
	return m.allTransactions, nil
}

type MockLogger struct {
	errorCalled bool
	infoCalled  bool
	lastError   error
	lastMessage string
}

func (m *MockLogger) Error(ctx context.Context, errs ...error) {
	m.errorCalled = true
	if len(errs) > 0 {
		m.lastError = errs[0]
	}
}

func (m *MockLogger) Info(ctx context.Context, messages ...string) {
	m.infoCalled = true
	if len(messages) > 0 {
		m.lastMessage = messages[0]
	}
}

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

func TestTransactionHandler_GetUserTransactions_UseCaseError(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		getUserError: errors.New("database error"),
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}

	if !mockLogger.errorCalled {
		t.Error("Expected logger.Error to be called")
	}

	if mockLogger.lastError == nil {
		t.Error("Expected error to be logged")
	}
}

func TestTransactionHandler_GetAllTransactions_UseCaseError(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		getAllError: errors.New("database error"),
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}

	if !mockLogger.errorCalled {
		t.Error("Expected logger.Error to be called")
	}

	if mockLogger.lastError == nil {
		t.Error("Expected error to be logged")
	}
}

func TestTransactionHandler_GetUserTransactions_LoggerError(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		getUserError: errors.New("database error"),
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if !mockLogger.errorCalled {
		t.Error("Expected logger.Error to be called")
	}

	if mockLogger.lastError == nil {
		t.Error("Expected error to be logged")
	}
}

func TestTransactionHandler_GetAllTransactions_LoggerError(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		getAllError: errors.New("database error"),
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if !mockLogger.errorCalled {
		t.Error("Expected logger.Error to be called")
	}

	if mockLogger.lastError == nil {
		t.Error("Expected error to be logged")
	}
}

func TestTransactionHandler_GetUserTransactions_MissingUserID_LoggerError(t *testing.T) {
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

	if !mockLogger.errorCalled {
		t.Error("Expected logger.Error to be called")
	}
}

func TestTransactionHandler_GetUserTransactions_EmptyUserID_LoggerError(t *testing.T) {
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

	if !mockLogger.errorCalled {
		t.Error("Expected logger.Error to be called")
	}
}

func TestTransactionHandler_GetUserTransactions_WithAllFilterTypes(t *testing.T) {
	testCases := []struct {
		name            string
		transactionType string
		expectedStatus  int
	}{
		{"bet filter", "bet", http.StatusOK},
		{"win filter", "win", http.StatusOK},
		{"invalid filter", "invalid", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUseCase := &MockTransactionUseCase{
				userTransactions: []*boundarydto.TransactionDTO{
					{
						ID:              utils.GenerateUUID(),
						UserID:          "user123",
						TransactionType: tc.transactionType,
						Amount:          100,
					},
				},
			}
			mockLogger := &MockLogger{}
			handler := NewTransactionHandler(mockUseCase, mockLogger)

			req, err := http.NewRequest("GET", "/transactions/user?user_id=user123&transaction_type="+tc.transactionType, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.GetUserTransactions(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, status)
			}
		})
	}
}

func TestTransactionHandler_GetAllTransactions_WithAllFilterTypes(t *testing.T) {
	testCases := []struct {
		name            string
		transactionType string
		expectedStatus  int
	}{
		{"bet filter", "bet", http.StatusOK},
		{"win filter", "win", http.StatusOK},
		{"invalid filter", "invalid", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUseCase := &MockTransactionUseCase{
				allTransactions: []*boundarydto.TransactionDTO{
					{
						ID:              utils.GenerateUUID(),
						UserID:          "user123",
						TransactionType: tc.transactionType,
						Amount:          100,
					},
				},
			}
			mockLogger := &MockLogger{}
			handler := NewTransactionHandler(mockUseCase, mockLogger)

			req, err := http.NewRequest("GET", "/transactions?transaction_type="+tc.transactionType, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.GetAllTransactions(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, status)
			}
		})
	}
}

func TestTransactionHandler_GetUserTransactions_ResponseHeaders(t *testing.T) {
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

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
	}
}

func TestTransactionHandler_GetAllTransactions_ResponseHeaders(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
	}
}

func TestTransactionHandler_GetUserTransactions_ComplexFilter(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		userTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user123",
				TransactionType: "bet",
				Amount:          100,
			},
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

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123&transaction_type=bet", nil)
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

	if len(response.Transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(response.Transactions))
	}
}

func TestTransactionHandler_GetAllTransactions_ComplexFilter(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions?transaction_type=bet", nil)
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

func TestTransactionHandler_GetUserTransactions_WithMultipleFilters(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123&transaction_type=bet&other_param=value", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetAllTransactions_WithMultipleFilters(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions?transaction_type=win&other_param=value", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetUserTransactions_WithSpecialCharacters(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		userTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user-123",
				TransactionType: "bet",
				Amount:          100,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user-123&transaction_type=bet", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetAllTransactions_WithSpecialCharacters(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		allTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user_456",
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

func TestTransactionHandler_GetUserTransactions_WithEmptyFilter(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123&transaction_type=", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetAllTransactions_WithEmptyFilter(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions?transaction_type=", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetUserTransactions_WithWhitespaceUserID(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		userTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "   ",
				TransactionType: "bet",
				Amount:          100,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=   ", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetUserTransactions_WithWhitespaceFilter(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123&transaction_type=   ", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetAllTransactions_WithWhitespaceFilter(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions?transaction_type=   ", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetUserTransactions_WithURLEncodedCharacters(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		userTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user%20123",
				TransactionType: "bet",
				Amount:          100,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user%20123&transaction_type=bet", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetAllTransactions_WithURLEncodedCharacters(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		allTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user%20456",
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

func TestTransactionHandler_GetUserTransactions_WithLongUserID(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		userTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "very_long_user_id_that_exceeds_normal_length_but_should_still_work",
				TransactionType: "bet",
				Amount:          100,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions/user?user_id=very_long_user_id_that_exceeds_normal_length_but_should_still_work&transaction_type=bet", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetAllTransactions_WithLongFilter(t *testing.T) {
	mockUseCase := &MockTransactionUseCase{
		allTransactions: []*boundarydto.TransactionDTO{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user123",
				TransactionType: "very_long_transaction_type_that_exceeds_normal_length",
				Amount:          100,
			},
		},
	}
	mockLogger := &MockLogger{}
	handler := NewTransactionHandler(mockUseCase, mockLogger)

	req, err := http.NewRequest("GET", "/transactions?transaction_type=very_long_transaction_type_that_exceeds_normal_length", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetUserTransactions_WithMultipleQueryParams(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions/user?user_id=user123&transaction_type=bet&limit=10&offset=0&sort=desc", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetUserTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetAllTransactions_WithMultipleQueryParams(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions?transaction_type=win&limit=10&offset=0&sort=desc", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAllTransactions(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestTransactionHandler_GetUserTransactions_JSONEncodingError(t *testing.T) {
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

	w := &failingResponseWriter{}

	handler.GetUserTransactions(w, req)

	if status := w.statusCode; status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}

	if !mockLogger.errorCalled {
		t.Error("Expected logger.Error to be called")
	}
}

func TestTransactionHandler_GetAllTransactions_JSONEncodingError(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := &failingResponseWriter{}

	handler.GetAllTransactions(w, req)

	if status := w.statusCode; status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}

	if !mockLogger.errorCalled {
		t.Error("Expected logger.Error to be called")
	}
}

type failingResponseWriter struct {
	statusCode int
	headers    http.Header
}

func (fw *failingResponseWriter) Header() http.Header {
	if fw.headers == nil {
		fw.headers = make(http.Header)
	}
	return fw.headers
}

func (fw *failingResponseWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write failed")
}

func (fw *failingResponseWriter) WriteHeader(statusCode int) {
	fw.statusCode = statusCode
}
