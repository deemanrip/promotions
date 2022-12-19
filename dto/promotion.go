package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Promotion struct {
	Id             uuid.UUID       `json:"id"`
	Price          decimal.Decimal `json:"price"`
	ExpirationDate string          `json:"expiration_date"`
}
