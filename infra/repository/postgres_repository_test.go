package repository

import (
	"errors"
	"testing"
	"time"

	"casino/boundary/repo_model"
	"casino/boundary/repository"
	"casino/domain/entity"
	"casino/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	cleanup := func() {
		mockDB.Close()
	}

	return db, mock, cleanup
}

func TestNewPostgresTransactionRepository(t *testing.T) {
	db, _, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)
	if repo == nil {
		t.Error("Expected repository to be created")
	}

	repoNil := NewPostgresTransactionRepository(nil)
	if repoNil == nil {
		t.Error("Expected repository to be created even with nil DB")
	}
}

func TestPostgresTransactionRepository_Save_Success(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	model := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO (.+)").WithArgs(
		model.ID,
		model.UserID,
		model.TransactionType,
		model.Amount,
		model.Timestamp,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Save(model)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_Save_Error(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	model := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "win",
		Amount:          200,
		Timestamp:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO (.+)").WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	err := repo.Save(model)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "failed to save transaction: database error" {
		t.Errorf("Expected specific error message, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_Save_NilTransaction(t *testing.T) {
	db, _, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	err := repo.Save(nil)
	if err == nil {
		t.Error("Expected error for nil transaction")
	}
}

func TestPostgresTransactionRepository_Save_NilDB(t *testing.T) {
	repo := NewPostgresTransactionRepository(nil)

	model := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	err := repo.Save(model)
	if err == nil {
		t.Error("Expected error for nil DB")
	}
}

func TestPostgresTransactionRepository_GetByUserID_Success(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	userID := utils.GenerateUUID()
	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "transaction_type", "amount", "timestamp"}).
		AddRow(utils.GenerateUUID(), userID, "bet", 100, time.Now()).
		AddRow(utils.GenerateUUID(), userID, "win", 200, time.Now())

	mock.ExpectQuery("SELECT (.+) FROM (.+) WHERE user_id = (.+) ORDER BY timestamp DESC").
		WithArgs(userID).
		WillReturnRows(expectedRows)

	models, err := repo.GetByUserID(userID, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(models))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_GetByUserID_WithTransactionType(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	userID := utils.GenerateUUID()
	transactionType := "bet"
	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "transaction_type", "amount", "timestamp"}).
		AddRow(utils.GenerateUUID(), userID, "bet", 100, time.Now())

	mock.ExpectQuery("SELECT (.+) FROM (.+) WHERE user_id = (.+) AND transaction_type = (.+) ORDER BY timestamp DESC").
		WithArgs(userID, "bet").
		WillReturnRows(expectedRows)

	models, err := repo.GetByUserID(userID, &transactionType)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(models))
	}

	if models[0].TransactionType != "bet" {
		t.Errorf("Expected transaction type bet, got %s", models[0].TransactionType)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_GetByUserID_EmptyResult(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	userID := utils.GenerateUUID()
	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "transaction_type", "amount", "timestamp"})

	mock.ExpectQuery("SELECT (.+) FROM (.+) WHERE user_id = (.+) ORDER BY timestamp DESC").
		WithArgs(userID).
		WillReturnRows(expectedRows)

	models, err := repo.GetByUserID(userID, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(models))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_GetByUserID_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	userID := utils.GenerateUUID()

	mock.ExpectQuery("SELECT (.+) FROM (.+) WHERE user_id = (.+) ORDER BY timestamp DESC").
		WithArgs(userID).
		WillReturnError(errors.New("database connection error"))

	models, err := repo.GetByUserID(userID, nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "failed to get transactions by user_id: database connection error" {
		t.Errorf("Expected specific error message, got %v", err)
	}

	if models != nil {
		t.Error("Expected nil models on error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_GetByUserID_NilDB(t *testing.T) {
	repo := NewPostgresTransactionRepository(nil)

	userID := utils.GenerateUUID()

	models, err := repo.GetByUserID(userID, nil)
	if err == nil {
		t.Error("Expected error for nil DB")
	}

	if models != nil {
		t.Error("Expected nil models on error")
	}
}

func TestPostgresTransactionRepository_GetAll_Success(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "transaction_type", "amount", "timestamp"}).
		AddRow(utils.GenerateUUID(), utils.GenerateUUID(), "bet", 100, time.Now()).
		AddRow(utils.GenerateUUID(), utils.GenerateUUID(), "win", 200, time.Now()).
		AddRow(utils.GenerateUUID(), utils.GenerateUUID(), "bet", 150, time.Now())

	mock.ExpectQuery("SELECT (.+) FROM (.+) ORDER BY timestamp DESC").
		WillReturnRows(expectedRows)

	models, err := repo.GetAll(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(models))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_GetAll_WithTransactionType(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	transactionType := "win"
	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "transaction_type", "amount", "timestamp"}).
		AddRow(utils.GenerateUUID(), utils.GenerateUUID(), "win", 200, time.Now()).
		AddRow(utils.GenerateUUID(), utils.GenerateUUID(), "win", 300, time.Now())

	mock.ExpectQuery("SELECT (.+) FROM (.+) WHERE transaction_type = (.+) ORDER BY timestamp DESC").
		WithArgs("win").
		WillReturnRows(expectedRows)

	models, err := repo.GetAll(&transactionType)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(models))
	}

	for _, model := range models {
		if model.TransactionType != "win" {
			t.Errorf("Expected transaction type win, got %s", model.TransactionType)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_GetAll_EmptyResult(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "transaction_type", "amount", "timestamp"})

	mock.ExpectQuery("SELECT (.+) FROM (.+) ORDER BY timestamp DESC").
		WillReturnRows(expectedRows)

	models, err := repo.GetAll(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(models) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(models))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_GetAll_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM (.+) ORDER BY timestamp DESC").
		WillReturnError(errors.New("database query error"))

	models, err := repo.GetAll(nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "failed to get all transactions: database query error" {
		t.Errorf("Expected specific error message, got %v", err)
	}

	if models != nil {
		t.Error("Expected nil models on error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_GetAll_NilDB(t *testing.T) {
	repo := NewPostgresTransactionRepository(nil)

	models, err := repo.GetAll(nil)
	if err == nil {
		t.Error("Expected error for nil DB")
	}

	if models != nil {
		t.Error("Expected nil models on error")
	}
}

func TestPostgresTransactionRepository_DataConversion(t *testing.T) {
	userID := utils.GenerateUUID()
	transactionType := entity.TransactionTypeBet
	amount := uint(100)
	timestamp := time.Now()

	transactionEntity := &entity.Transaction{
		ID:              utils.GenerateUUID(),
		UserID:          userID,
		TransactionType: transactionType,
		Amount:          amount,
		Timestamp:       timestamp,
	}

	model := &repo_model.TransactionModel{}
	model.FromEntity(transactionEntity)

	convertedEntity := model.ToEntity()

	if transactionEntity.ID != model.ID {
		t.Errorf("ID conversion failed: %s != %s", transactionEntity.ID, model.ID)
	}

	if transactionEntity.UserID != model.UserID {
		t.Errorf("UserID conversion failed: %s != %s", transactionEntity.UserID, model.UserID)
	}

	if string(transactionEntity.TransactionType) != model.TransactionType {
		t.Errorf("TransactionType conversion failed: %s != %s", transactionEntity.TransactionType, model.TransactionType)
	}

	if transactionEntity.Amount != model.Amount {
		t.Errorf("Amount conversion failed: %d != %d", transactionEntity.Amount, model.Amount)
	}

	if convertedEntity.ID != transactionEntity.ID {
		t.Errorf("Entity conversion failed: %s != %s", convertedEntity.ID, transactionEntity.ID)
	}

	if convertedEntity.UserID != transactionEntity.UserID {
		t.Errorf("UserID conversion failed: %s != %s", convertedEntity.UserID, transactionEntity.UserID)
	}

	if convertedEntity.TransactionType != transactionEntity.TransactionType {
		t.Errorf("TransactionType conversion failed: %s != %s", convertedEntity.TransactionType, transactionEntity.TransactionType)
	}

	if convertedEntity.Amount != transactionEntity.Amount {
		t.Errorf("Amount conversion failed: %d != %d", convertedEntity.Amount, transactionEntity.Amount)
	}
}

func TestPostgresTransactionRepository_TransactionTypeConversion(t *testing.T) {
	betType := entity.TransactionTypeBet
	winType := entity.TransactionTypeWin

	if string(betType) != "bet" {
		t.Errorf("Expected 'bet', got %s", string(betType))
	}

	if string(winType) != "win" {
		t.Errorf("Expected 'win', got %s", string(winType))
	}

	betFromString := entity.TransactionType("bet")
	winFromString := entity.TransactionType("win")

	if betFromString != entity.TransactionTypeBet {
		t.Errorf("Expected TransactionTypeBet, got %s", betFromString)
	}

	if winFromString != entity.TransactionTypeWin {
		t.Errorf("Expected TransactionTypeWin, got %s", winFromString)
	}
}

func TestPostgresTransactionRepository_EdgeCases(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	transaction := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          0,
		Timestamp:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO (.+)").WithArgs(
		transaction.ID,
		transaction.UserID,
		transaction.TransactionType,
		transaction.Amount,
		transaction.Timestamp,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Save(transaction)
	if err != nil {
		t.Errorf("Expected no error for zero amount, got %v", err)
	}

	largeTransaction := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "win",
		Amount:          999999999,
		Timestamp:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO (.+)").WithArgs(
		largeTransaction.ID,
		largeTransaction.UserID,
		largeTransaction.TransactionType,
		largeTransaction.Amount,
		largeTransaction.Timestamp,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Save(largeTransaction)
	if err != nil {
		t.Errorf("Expected no error for large amount, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_MethodSignatures(t *testing.T) {
	repo := NewPostgresTransactionRepository(nil)

	model := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	_ = repo.Save
	_ = repo.GetByUserID
	_ = repo.GetAll

	if repo == nil {
		t.Error("Repository should not be nil")
	}

	if model == nil {
		t.Error("Model should not be nil")
	}
}

func TestPostgresTransactionRepository_InterfaceCompliance(t *testing.T) {
	db, _, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	var _ repository.TransactionRepository = repo

	_ = repo.Save
	_ = repo.GetByUserID
	_ = repo.GetAll
}

func TestPostgresTransactionRepository_ModelStructure(t *testing.T) {
	model := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	if model.ID == "" {
		t.Error("Model ID should not be empty")
	}

	if model.UserID == "" {
		t.Error("Model UserID should not be empty")
	}

	if model.TransactionType != "bet" && model.TransactionType != "win" {
		t.Errorf("Model TransactionType should be 'bet' or 'win', got %s", model.TransactionType)
	}

	if model.Amount == 0 {
		t.Error("Model Amount should not be zero")
	}

	if model.Timestamp.IsZero() {
		t.Error("Model Timestamp should not be zero")
	}
}

func TestPostgresTransactionRepository_ErrorHandling(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	transaction := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          100,
		Timestamp:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO (.+)").WillReturnError(errors.New("constraint violation"))
	mock.ExpectRollback()

	err := repo.Save(transaction)
	if err == nil {
		t.Error("Expected error for constraint violation")
	}

	if err.Error() != "failed to save transaction: constraint violation" {
		t.Errorf("Expected specific error message, got %v", err)
	}

	userID := utils.GenerateUUID()
	mock.ExpectQuery("SELECT (.+) FROM (.+) WHERE user_id = (.+) ORDER BY timestamp DESC").
		WithArgs(userID).
		WillReturnError(errors.New("connection timeout"))

	_, err = repo.GetByUserID(userID, nil)
	if err == nil {
		t.Error("Expected error for connection timeout")
	}

	if err.Error() != "failed to get transactions by user_id: connection timeout" {
		t.Errorf("Expected specific error message, got %v", err)
	}

	mock.ExpectQuery("SELECT (.+) FROM (.+) ORDER BY timestamp DESC").
		WillReturnError(errors.New("table not found"))

	_, err = repo.GetAll(nil)
	if err == nil {
		t.Error("Expected error for table not found")
	}

	if err.Error() != "failed to get all transactions: table not found" {
		t.Errorf("Expected specific error message, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_BoundaryValues(t *testing.T) {
	db, mock, cleanup := setupMockTestDB(t)
	defer cleanup()

	repo := NewPostgresTransactionRepository(db)

	minTransaction := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "bet",
		Amount:          1,
		Timestamp:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO (.+)").WithArgs(
		minTransaction.ID,
		minTransaction.UserID,
		minTransaction.TransactionType,
		minTransaction.Amount,
		minTransaction.Timestamp,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Save(minTransaction)
	if err != nil {
		t.Errorf("Expected no error for minimum values, got %v", err)
	}

	maxTransaction := &repo_model.TransactionModel{
		ID:              utils.GenerateUUID(),
		UserID:          utils.GenerateUUID(),
		TransactionType: "win",
		Amount:          2147483647,
		Timestamp:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO (.+)").WithArgs(
		maxTransaction.ID,
		maxTransaction.UserID,
		maxTransaction.TransactionType,
		maxTransaction.Amount,
		maxTransaction.Timestamp,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Save(maxTransaction)
	if err != nil {
		t.Errorf("Expected no error for maximum values, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
