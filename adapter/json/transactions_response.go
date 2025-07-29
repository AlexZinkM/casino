package json

import (
	"casino/boundary/dto"
)

type TransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
}

func (r *TransactionsResponse) FromDtos(dtos []*dto.TransactionDTO) {
	r.Transactions = make([]TransactionResponse, len(dtos))
	for i, dto := range dtos {
		r.Transactions[i].FromDto(dto)
	}
}

func (r *TransactionsResponse) ToDtos() []*dto.TransactionDTO {
	dtos := make([]*dto.TransactionDTO, len(r.Transactions))
	for i, response := range r.Transactions {
		dtos[i] = response.ToDto()
	}
	return dtos
}
