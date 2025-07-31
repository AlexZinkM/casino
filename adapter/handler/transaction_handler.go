package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	adapterjson "casino/adapter/json"
	boundarydto "casino/boundary/dto"
	"casino/boundary/logging"
	"casino/boundary/usecase"
)

type TransactionHandler struct {
	transactionUseCase usecase.TransactionUseCase
	logger             logging.Logger
}

func NewTransactionHandler(transactionUseCase usecase.TransactionUseCase, logger logging.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
		logger:             logger,
	}
}

// GetUserTransactions godoc
// @Summary Get user transactions
// @Description Get transactions for a specific user with optional filtering
// @Tags transactions
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param transaction_type query string false "Transaction type filter (bet or win)"
// @Success 200 {object} json.TransactionsResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /transactions/user [get]
func (h *TransactionHandler) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.logger.Error(r.Context(), fmt.Errorf("user_id is required"))
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	transactionTypeStr := r.URL.Query().Get("transaction_type")
	var filter *boundarydto.TransactionFilterDTO
	if transactionTypeStr != "" {
		filter = &boundarydto.TransactionFilterDTO{
			TransactionType: &transactionTypeStr,
		}
	}

	dtos, err := h.transactionUseCase.GetUserTransactions(userID, filter)
	if err != nil {
		h.logger.Error(r.Context(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := adapterjson.TransactionsResponse{}
	response.FromDtos(dtos)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(r.Context(), err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetAllTransactions godoc
// @Summary Get all transactions
// @Description Get all transactions with optional filtering
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction_type query string false "Transaction type filter (bet or win)"
// @Success 200 {object} json.TransactionsResponse
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactionTypeStr := r.URL.Query().Get("transaction_type")
	var filter *boundarydto.TransactionFilterDTO
	if transactionTypeStr != "" {
		filter = &boundarydto.TransactionFilterDTO{
			TransactionType: &transactionTypeStr,
		}
	}

	dtos, err := h.transactionUseCase.GetAllTransactions(filter)
	if err != nil {
		h.logger.Error(r.Context(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := adapterjson.TransactionsResponse{}
	response.FromDtos(dtos)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(r.Context(), err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
