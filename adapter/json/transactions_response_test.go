package json

import (
	"testing"

	boundarydto "casino/boundary/dto"
	"casino/utils"
)

func TestTransactionsResponse_FromDtos(t *testing.T) {
	boundaryDtos := []*boundarydto.TransactionDTO{
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
	}

	response := &TransactionsResponse{}
	response.FromDtos(boundaryDtos)

	if len(response.Transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(response.Transactions))
	}

	if response.Transactions[0].UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got %s", response.Transactions[0].UserID)
	}

	if response.Transactions[1].UserID != "user456" {
		t.Errorf("Expected UserID 'user456', got %s", response.Transactions[1].UserID)
	}
}

func TestTransactionsResponse_ToDtos(t *testing.T) {
	response := &TransactionsResponse{
		Transactions: []TransactionResponse{
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user123",
				TransactionType: "bet",
				Amount:          100,
				Timestamp:       "2023-01-01T00:00:00Z",
			},
			{
				ID:              utils.GenerateUUID(),
				UserID:          "user456",
				TransactionType: "win",
				Amount:          200,
				Timestamp:       "2023-01-02T00:00:00Z",
			},
		},
	}

	boundaryDtos := response.ToDtos()

	if len(boundaryDtos) != 2 {
		t.Errorf("Expected 2 DTOs, got %d", len(boundaryDtos))
	}

	if boundaryDtos[0].UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got %s", boundaryDtos[0].UserID)
	}

	if boundaryDtos[1].UserID != "user456" {
		t.Errorf("Expected UserID 'user456', got %s", boundaryDtos[1].UserID)
	}
}

func TestTransactionsResponse_FromDtos_Empty(t *testing.T) {
	response := &TransactionsResponse{}
	response.FromDtos([]*boundarydto.TransactionDTO{})

	if len(response.Transactions) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(response.Transactions))
	}
}

func TestTransactionsResponse_ToDtos_Empty(t *testing.T) {
	response := &TransactionsResponse{
		Transactions: []TransactionResponse{},
	}

	boundaryDtos := response.ToDtos()

	if len(boundaryDtos) != 0 {
		t.Errorf("Expected 0 DTOs, got %d", len(boundaryDtos))
	}
}
