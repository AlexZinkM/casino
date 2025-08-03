package utils

import (
	"fmt"
	"strings"
)

type TransactionAlreadyExistsError struct {
	TransactionID string
}

func (e *TransactionAlreadyExistsError) Error() string {
	return fmt.Sprintf("transaction with id %s already exists", e.TransactionID)
}

func IsTransactionAlreadyExists(err error) bool {
	_, ok := err.(*TransactionAlreadyExistsError)
	return ok
}

func IsDatabaseConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "failed to connect") ||
		strings.Contains(errStr, "dial error") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "No connection could be made")
}
