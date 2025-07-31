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

