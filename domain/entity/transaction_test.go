package entity

import (
	"testing"
	"time"

	"casino/utils"
)

func TestNewTransaction(t *testing.T) {
	userID := utils.GenerateUUID()
	transactionType := TransactionTypeBet
	amount := uint(100)

	transaction := NewTransaction(userID, transactionType, amount)

	if transaction.ID == "" {
		t.Error("Expected transaction ID to be generated")
	}

	if transaction.UserID != userID {
		t.Errorf("Expected UserID to be %s, got %s", userID, transaction.UserID)
	}

	if transaction.TransactionType != transactionType {
		t.Errorf("Expected TransactionType to be %s, got %s", transactionType, transaction.TransactionType)
	}

	if transaction.Amount != amount {
		t.Errorf("Expected Amount to be %d, got %d", amount, transaction.Amount)
	}

	if transaction.Timestamp.IsZero() {
		t.Error("Expected Timestamp to be set")
	}

	if time.Since(transaction.Timestamp) > time.Second {
		t.Error("Expected Timestamp to be recent")
	}
}

func TestTransactionTypeConstants(t *testing.T) {
	if TransactionTypeBet != "bet" {
		t.Errorf("Expected TransactionTypeBet to be 'bet', got %s", TransactionTypeBet)
	}

	if TransactionTypeWin != "win" {
		t.Errorf("Expected TransactionTypeWin to be 'win', got %s", TransactionTypeWin)
	}
}
