package transactions

import (
	"gorm.io/datatypes"
)

type TransactionRequest struct {
	Amount    float64 `json:"amount" validate:"required,gt=0"`
	Type      string  `json:"type" validate:"required,oneof=income expense"`
	Categorie string  `json:"categorie" validate:"required,alphaunicode,min=2,max=50"`
}

type TransactionResponse struct {
	Id        uint           `json:"id"`
	Amount    float64        `json:"amount"`
	Type      string         `json:"type"`
	Categorie string         `json:"categorie"`
	Date      datatypes.Date `json:"date"`
}

type GetAllResponse struct {
	Transns *[]Transaction `json:"transactions"`
	Count   int64          `json:"count"`
}
