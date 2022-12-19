package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Promotion struct {
	Id             uuid.UUID       `json:"id" example:"d9433531-5b0a-431d-82d4-b413dc34253f"`
	Price          decimal.Decimal `json:"price" example:"32.180885"`
	ExpirationDate string          `json:"expiration_date" example:"2018-08-10 12:47:53 +0200 CEST"`
}
