package repo_model

import (
	"time"

	"casino/domain/entity"
)

type TransactionModel struct {
	ID              string    `gorm:"primaryKey;type:uuid"`
	UserID          string    `gorm:"type:uuid;not null"`
	TransactionType string    `gorm:"type:varchar(10);not null;check:transaction_type IN ('bet', 'win')"`
	Amount          uint      `gorm:"type:integer;not null;check:amount > 0"`
	Timestamp       time.Time `gorm:"type:timestamp;not null"`
}

func (TransactionModel) TableName() string {
	return "transactions"
}

func (m *TransactionModel) ToEntity() *entity.Transaction {
	return &entity.Transaction{
		ID:              m.ID,
		UserID:          m.UserID,
		TransactionType: entity.TransactionType(m.TransactionType),
		Amount:          m.Amount,
		Timestamp:       m.Timestamp,
	}
}

func (m *TransactionModel) FromEntity(entity *entity.Transaction) {
	m.ID = entity.ID
	m.UserID = entity.UserID
	m.TransactionType = string(entity.TransactionType)
	m.Amount = entity.Amount
	m.Timestamp = entity.Timestamp
}
