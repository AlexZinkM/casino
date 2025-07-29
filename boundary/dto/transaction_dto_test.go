package dto

import (
	"testing"
	"time"

	"casino/domain/entity"
	"casino/utils"
)

func TestTransactionDTO_FromEntity(t *testing.T) {
	transactionEntity := &entity.Transaction{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: entity.TransactionTypeBet,
		Amount:          100,
		Timestamp:       time.Now(),
	}

	dto := &TransactionDTO{}
	dto.FromEntity(transactionEntity)

	if dto.ID != transactionEntity.ID {
		t.Errorf("Expected ID %s, got %s", transactionEntity.ID, dto.ID)
	}

	if dto.UserID != transactionEntity.UserID {
		t.Errorf("Expected UserID %s, got %s", transactionEntity.UserID, dto.UserID)
	}

	if dto.TransactionType != string(transactionEntity.TransactionType) {
		t.Errorf("Expected TransactionType %s, got %s", transactionEntity.TransactionType, dto.TransactionType)
	}

	if dto.Amount != transactionEntity.Amount {
		t.Errorf("Expected Amount %d, got %d", transactionEntity.Amount, dto.Amount)
	}

	if dto.Timestamp != transactionEntity.Timestamp {
		t.Errorf("Expected Timestamp %v, got %v", transactionEntity.Timestamp, dto.Timestamp)
	}
}

func TestTransactionDTO_ToEntity(t *testing.T) {
	dto := &TransactionDTO{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "win",
		Amount:          200,
		Timestamp:       time.Now(),
	}

	transactionEntity := dto.ToEntity()

	if transactionEntity.ID != dto.ID {
		t.Errorf("Expected ID %s, got %s", dto.ID, transactionEntity.ID)
	}

	if transactionEntity.UserID != dto.UserID {
		t.Errorf("Expected UserID %s, got %s", dto.UserID, transactionEntity.UserID)
	}

	if transactionEntity.TransactionType != entity.TransactionType(dto.TransactionType) {
		t.Errorf("Expected TransactionType %s, got %s", dto.TransactionType, transactionEntity.TransactionType)
	}

	if transactionEntity.Amount != dto.Amount {
		t.Errorf("Expected Amount %d, got %d", dto.Amount, transactionEntity.Amount)
	}

	if transactionEntity.Timestamp != dto.Timestamp {
		t.Errorf("Expected Timestamp %v, got %v", dto.Timestamp, transactionEntity.Timestamp)
	}
}

func TestCreateTransactionDTO_FromEntity(t *testing.T) {
	transactionEntity := &entity.Transaction{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: entity.TransactionTypeWin,
		Amount:          150,
		Timestamp:       time.Now(),
	}

	dto := &CreateTransactionDTO{}
	dto.FromEntity(transactionEntity)

	if dto.UserID != transactionEntity.UserID {
		t.Errorf("Expected UserID %s, got %s", transactionEntity.UserID, dto.UserID)
	}

	if dto.TransactionType != string(transactionEntity.TransactionType) {
		t.Errorf("Expected TransactionType %s, got %s", transactionEntity.TransactionType, dto.TransactionType)
	}

	if dto.Amount != transactionEntity.Amount {
		t.Errorf("Expected Amount %d, got %d", transactionEntity.Amount, dto.Amount)
	}
}

func TestCreateTransactionDTO_ToEntity(t *testing.T) {
	dto := &CreateTransactionDTO{
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          75,
	}

	transactionEntity := dto.ToEntity()

	if transactionEntity.UserID != dto.UserID {
		t.Errorf("Expected UserID %s, got %s", dto.UserID, transactionEntity.UserID)
	}

	if transactionEntity.TransactionType != entity.TransactionType(dto.TransactionType) {
		t.Errorf("Expected TransactionType %s, got %s", dto.TransactionType, transactionEntity.TransactionType)
	}

	if transactionEntity.Amount != dto.Amount {
		t.Errorf("Expected Amount %d, got %d", dto.Amount, transactionEntity.Amount)
	}

	if transactionEntity.Timestamp.IsZero() {
		t.Error("Expected Timestamp to be set")
	}
}

func TestTransactionFilterDTO_ToEntity(t *testing.T) {
	transactionType := "bet"
	filter := &TransactionFilterDTO{
		UserID:          stringPtr("user123"),
		TransactionType: &transactionType,
	}

	entityType := filter.ToEntity()

	if entityType == nil {
		t.Error("Expected entity to be created")
	}

	if *entityType != entity.TransactionTypeBet {
		t.Errorf("Expected TransactionType 'bet', got %s", *entityType)
	}
}

func TestTransactionFilterDTO_ToEntity_Nil(t *testing.T) {
	filter := &TransactionFilterDTO{
		UserID:          stringPtr("user123"),
		TransactionType: nil,
	}

	entity := filter.ToEntity()

	if entity != nil {
		t.Error("Expected entity to be nil")
	}
}

func stringPtr(s string) *string {
	return &s
}
