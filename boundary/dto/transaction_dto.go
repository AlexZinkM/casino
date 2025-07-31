package dto

import (
	"casino/domain/entity"
	"casino/utils"
	"time"
)

type TransactionDTO struct {
	ID              string   
	UserID          string    
	TransactionType string    
	Amount          uint      
	Timestamp       time.Time 
}

func (d *TransactionDTO) FromEntity(entity *entity.Transaction) {
	d.ID = entity.ID
	d.UserID = entity.UserID
	d.TransactionType = string(entity.TransactionType)
	d.Amount = entity.Amount
	d.Timestamp = entity.Timestamp
}

func (d *TransactionDTO) ToEntity() *entity.Transaction {
	return &entity.Transaction{
		ID:              d.ID,
		UserID:          d.UserID,
		TransactionType: entity.TransactionType(d.TransactionType),
		Amount:          d.Amount,
		Timestamp:       d.Timestamp,
	}
}

type CreateTransactionDTO struct {
	UserID          string 
	TransactionType string 
	Amount          uint   
}

func (d *CreateTransactionDTO) FromEntity(entity *entity.Transaction) {
	d.UserID = entity.UserID
	d.TransactionType = string(entity.TransactionType)
	d.Amount = entity.Amount
}

func (d *CreateTransactionDTO) ToEntity() *entity.Transaction {
	return &entity.Transaction{
		ID:              utils.GenerateUUID(),
		UserID:          d.UserID,
		TransactionType: entity.TransactionType(d.TransactionType),
		Amount:          d.Amount,
		Timestamp:       time.Now(),
	}
}

type TransactionFilterDTO struct {
	UserID          *string 
	TransactionType *string 
}

func (d *TransactionFilterDTO) ToEntity() *entity.TransactionType {
	if d.TransactionType == nil {
		return nil
	}
	tt := entity.TransactionType(*d.TransactionType)
	return &tt
}
