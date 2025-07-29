package json

import (
	"testing"
	"time"

	boundarydto "casino/boundary/dto"
	"casino/utils"
)

func TestTransactionResponse_FromDto(t *testing.T) {
	boundaryDto := &boundarydto.TransactionDTO{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	response := &TransactionResponse{}
	response.FromDto(boundaryDto)

	if response.ID != boundaryDto.ID {
		t.Errorf("Expected ID %s, got %s", boundaryDto.ID, response.ID)
	}

	if response.UserID != boundaryDto.UserID {
		t.Errorf("Expected UserID %s, got %s", boundaryDto.UserID, response.UserID)
	}

	if response.TransactionType != boundaryDto.TransactionType {
		t.Errorf("Expected TransactionType %s, got %s", boundaryDto.TransactionType, response.TransactionType)
	}

	if response.Amount != boundaryDto.Amount {
		t.Errorf("Expected Amount %d, got %d", boundaryDto.Amount, response.Amount)
	}

	if response.Timestamp == "" {
		t.Error("Expected Timestamp to be formatted")
	}
}

func TestTransactionResponse_ToDto(t *testing.T) {
	response := &TransactionResponse{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "win",
		Amount:          200,
		Timestamp:       time.Now().Format(time.RFC3339),
	}

	boundaryDto := response.ToDto()

	if boundaryDto.ID != response.ID {
		t.Errorf("Expected ID %s, got %s", response.ID, boundaryDto.ID)
	}

	if boundaryDto.UserID != response.UserID {
		t.Errorf("Expected UserID %s, got %s", response.UserID, boundaryDto.UserID)
	}

	if boundaryDto.TransactionType != response.TransactionType {
		t.Errorf("Expected TransactionType %s, got %s", response.TransactionType, boundaryDto.TransactionType)
	}

	if boundaryDto.Amount != response.Amount {
		t.Errorf("Expected Amount %d, got %d", response.Amount, boundaryDto.Amount)
	}
}
