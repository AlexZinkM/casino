package usecase

import (
	"casino/boundary/dto"
	"casino/boundary/repo_model"
	"casino/boundary/repository"
	"casino/utils"
	"fmt"
)

type TransactionUseCaseImpl struct {
	transactionRepo repository.TransactionRepository
}

func NewTransactionUseCaseImpl(transactionRepo repository.TransactionRepository) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		transactionRepo: transactionRepo,
	}
}

func (uc *TransactionUseCaseImpl) ProcessTransaction(dto *dto.CreateTransactionDTO) error {
	existingTransaction, err := uc.transactionRepo.GetByID(dto.ID)
	if err != nil {
		return fmt.Errorf("failed to check for existing transaction: %w", err)
	}

	if existingTransaction != nil {
		return &utils.TransactionAlreadyExistsError{TransactionID: dto.ID}
	}

	entity := dto.ToEntity()
	model := &repo_model.TransactionModel{}
	model.FromEntity(entity)
	return uc.transactionRepo.Save(model)
}

func (uc *TransactionUseCaseImpl) GetUserTransactions(userID string, filter *dto.TransactionFilterDTO) ([]*dto.TransactionDTO, error) {
	var transactionType *string
	if filter != nil && filter.TransactionType != nil {
		transactionType = filter.TransactionType
	}

	models, err := uc.transactionRepo.GetByUserID(userID, transactionType)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.TransactionDTO, len(models))
	for i, model := range models {
		entity := model.ToEntity()
		dtos[i] = &dto.TransactionDTO{}
		dtos[i].FromEntity(entity)
	}

	return dtos, nil
}

func (uc *TransactionUseCaseImpl) GetAllTransactions(filter *dto.TransactionFilterDTO) ([]*dto.TransactionDTO, error) {
	var transactionType *string
	if filter != nil && filter.TransactionType != nil {
		transactionType = filter.TransactionType
	}

	models, err := uc.transactionRepo.GetAll(transactionType)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.TransactionDTO, len(models))
	for i, model := range models {
		entity := model.ToEntity()
		dtos[i] = &dto.TransactionDTO{}
		dtos[i].FromEntity(entity)
	}

	return dtos, nil
}
