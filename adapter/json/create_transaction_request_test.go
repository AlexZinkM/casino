package json

import (
	"testing"

	boundarydto "casino/boundary/dto"
	"casino/utils"
)

func TestCreateTransactionRequest_ToDto(t *testing.T) {
	request := &CreateTransactionRequest{
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
	}

	boundaryDto := request.ToDto()

	if boundaryDto.UserID != request.UserID {
		t.Errorf("Expected UserID %s, got %s", request.UserID, boundaryDto.UserID)
	}

	if boundaryDto.TransactionType != request.TransactionType {
		t.Errorf("Expected TransactionType %s, got %s", request.TransactionType, boundaryDto.TransactionType)
	}

	if boundaryDto.Amount != request.Amount {
		t.Errorf("Expected Amount %d, got %d", request.Amount, boundaryDto.Amount)
	}
}

func TestCreateTransactionRequest_FromDto(t *testing.T) {
	boundaryDto := &boundarydto.CreateTransactionDTO{
		UserID:          utils.GenerateUUID(),
		TransactionType: "win",
		Amount:          200,
	}

	request := &CreateTransactionRequest{}
	request.FromDto(boundaryDto)

	if request.UserID != boundaryDto.UserID {
		t.Errorf("Expected UserID %s, got %s", boundaryDto.UserID, request.UserID)
	}

	if request.TransactionType != boundaryDto.TransactionType {
		t.Errorf("Expected TransactionType %s, got %s", boundaryDto.TransactionType, request.TransactionType)
	}

	if request.Amount != boundaryDto.Amount {
		t.Errorf("Expected Amount %d, got %d", boundaryDto.Amount, request.Amount)
	}
}
