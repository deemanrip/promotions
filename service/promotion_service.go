package service

import (
	"context"
	"fmt"
	"github.com/deemanrip/promotions/dto"
	"github.com/deemanrip/promotions/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

const promotionQuery = "SELECT * FROM promotions.promotions where id = '%v'"

func GetPromotionById(promotionId *string) (*dto.Promotion, error) {
	clickhouseConn := repository.ClickhouseConnection
	query := fmt.Sprintf(promotionQuery, *promotionId)
	rows, err := clickhouseConn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	if hasRows := rows.Next(); hasRows {
		var (
			id             uuid.UUID
			price          decimal.Decimal
			expirationDate string
		)
		if err := rows.Scan(&id, &price, &expirationDate); err != nil {
			return nil, err
		}
		return &dto.Promotion{Id: id, Price: price, ExpirationDate: expirationDate}, nil
	}
	if closeErr := rows.Close(); closeErr != nil {
		log.Error(closeErr)
	}

	return nil, nil
}
