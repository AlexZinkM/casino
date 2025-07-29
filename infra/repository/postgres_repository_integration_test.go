package repository

import (
	"testing"
	"time"

	"casino/boundary/repo_model"
	"casino/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&repo_model.TransactionModel{}); err != nil {
		t.Skipf("Skipping test: failed to migrate database: %v", err)
	}

	return db
}

func TestPostgresTransactionRepository_Integration_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresTransactionRepository(db)

	model := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	err := repo.Save(model)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var savedModel repo_model.TransactionModel
	if err := db.First(&savedModel, "id = ?", model.ID).Error; err != nil {
		t.Errorf("Expected transaction to be saved, got error: %v", err)
	}

	if savedModel.ID != model.ID {
		t.Errorf("Expected ID %s, got %s", model.ID, savedModel.ID)
	}

	if savedModel.UserID != model.UserID {
		t.Errorf("Expected UserID %s, got %s", model.UserID, savedModel.UserID)
	}

	if savedModel.TransactionType != model.TransactionType {
		t.Errorf("Expected TransactionType %s, got %s", model.TransactionType, savedModel.TransactionType)
	}

	if savedModel.Amount != model.Amount {
		t.Errorf("Expected Amount %d, got %d", model.Amount, savedModel.Amount)
	}
}

func TestPostgresTransactionRepository_Integration_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresTransactionRepository(db)

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

	repo.Save(model1)
	repo.Save(model2)

	models, err := repo.GetByUserID(userID, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(models))
	}

	transactionType := "bet"
	models, err = repo.GetByUserID(userID, &transactionType)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(models))
	}

	if models[0].TransactionType != "bet" {
		t.Errorf("Expected transaction type 'bet', got %s", models[0].TransactionType)
	}
}

func TestPostgresTransactionRepository_Integration_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresTransactionRepository(db)

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
	model3 := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          150,
		Timestamp:       time.Now(),
	}

	repo.Save(model1)
	repo.Save(model2)
	repo.Save(model3)

	models, err := repo.GetAll(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(models))
	}

	transactionType := "win"
	models, err = repo.GetAll(&transactionType)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(models))
	}

	if models[0].TransactionType != "win" {
		t.Errorf("Expected transaction type 'win', got %s", models[0].TransactionType)
	}
}
