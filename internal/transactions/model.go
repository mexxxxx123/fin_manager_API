package transactions

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
	Categorie string  `json:"categorie"`
	UserId    uint    `json:"userId"`
}

func NewTransaction(amount float64, tyype string, categorie string, userId uint) *Transaction {
	Transaction := &Transaction{
		Amount:    amount,
		Type:      tyype,
		Categorie: categorie,
		UserId:    userId,
	}

	return Transaction
}
