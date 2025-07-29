package repo_model

import (
	"testing"
	"time"

	"casino/domain/entity"
	"casino/utils"
)

func TestTransactionModel_Structure(t *testing.T) {
	model := &TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	if model.ID == "" {
		t.Error("Expected ID to be set")
	}

	if model.UserID == "" {
		t.Error("Expected UserID to be set")
	}

	if model.TransactionType != "bet" {
		t.Errorf("Expected TransactionType 'bet', got %s", model.TransactionType)
	}

	if model.Amount != 100 {
		t.Errorf("Expected Amount 100, got %d", model.Amount)
	}

	if model.Timestamp.IsZero() {
		t.Error("Expected Timestamp to be set")
	}
}

func TestTransactionModel_TransactionTypes(t *testing.T) {
	betModel := &TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	winModel := &TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "win",
		Amount:          200,
		Timestamp:       time.Now(),
	}

	if betModel.TransactionType != "bet" {
		t.Errorf("Expected TransactionType 'bet', got %s", betModel.TransactionType)
	}

	if winModel.TransactionType != "win" {
		t.Errorf("Expected TransactionType 'win', got %s", winModel.TransactionType)
	}
}

func TestTransactionModel_AmountValidation(t *testing.T) {
	model := &TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          0,
		Timestamp:       time.Now(),
	}

	if model.Amount != 0 {
		t.Errorf("Expected Amount 0, got %d", model.Amount)
	}

	model.Amount = 100
	if model.Amount != 100 {
		t.Errorf("Expected Amount 100, got %d", model.Amount)
	}
}

func TestTransactionModel_ToEntity(t *testing.T) {
	model := &TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	transactionEntity := model.ToEntity()

	if transactionEntity.ID != model.ID {
		t.Errorf("Expected ID %s, got %s", model.ID, transactionEntity.ID)
	}

	if transactionEntity.UserID != model.UserID {
		t.Errorf("Expected UserID %s, got %s", model.UserID, transactionEntity.UserID)
	}

	if transactionEntity.TransactionType != entity.TransactionTypeBet {
		t.Errorf("Expected TransactionType %s, got %s", entity.TransactionTypeBet, transactionEntity.TransactionType)
	}

	if transactionEntity.Amount != model.Amount {
		t.Errorf("Expected Amount %d, got %d", model.Amount, transactionEntity.Amount)
	}

	if transactionEntity.Timestamp != model.Timestamp {
		t.Errorf("Expected Timestamp %v, got %v", model.Timestamp, transactionEntity.Timestamp)
	}
}

func TestTransactionModel_FromEntity(t *testing.T) {
	entity := &entity.Transaction{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: entity.TransactionTypeWin,
		Amount:          200,
		Timestamp:       time.Now(),
	}

	model := &TransactionModel{}
	model.FromEntity(entity)

	if model.ID != entity.ID {
		t.Errorf("Expected ID %s, got %s", entity.ID, model.ID)
	}

	if model.UserID != entity.UserID {
		t.Errorf("Expected UserID %s, got %s", entity.UserID, model.UserID)
	}

	if model.TransactionType != string(entity.TransactionType) {
		t.Errorf("Expected TransactionType %s, got %s", string(entity.TransactionType), model.TransactionType)
	}

	if model.Amount != entity.Amount {
		t.Errorf("Expected Amount %d, got %d", entity.Amount, model.Amount)
	}

	if model.Timestamp != entity.Timestamp {
		t.Errorf("Expected Timestamp %v, got %v", entity.Timestamp, model.Timestamp)
	}
}

func TestTransactionModel_ConversionRoundTrip(t *testing.T) {
	originalEntity := &entity.Transaction{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: entity.TransactionTypeBet,
		Amount:          150,
		Timestamp:       time.Now(),
	}

	model := &TransactionModel{}
	model.FromEntity(originalEntity)

	convertedEntity := model.ToEntity()

	if originalEntity.ID != convertedEntity.ID {
		t.Errorf("ID round-trip failed: %s != %s", originalEntity.ID, convertedEntity.ID)
	}

	if originalEntity.UserID != convertedEntity.UserID {
		t.Errorf("UserID round-trip failed: %s != %s", originalEntity.UserID, convertedEntity.UserID)
	}

	if originalEntity.TransactionType != convertedEntity.TransactionType {
		t.Errorf("TransactionType round-trip failed: %s != %s", originalEntity.TransactionType, convertedEntity.TransactionType)
	}

	if originalEntity.Amount != convertedEntity.Amount {
		t.Errorf("Amount round-trip failed: %d != %d", originalEntity.Amount, convertedEntity.Amount)
	}

	if originalEntity.Timestamp != convertedEntity.Timestamp {
		t.Errorf("Timestamp round-trip failed: %v != %v", originalEntity.Timestamp, convertedEntity.Timestamp)
	}
}
