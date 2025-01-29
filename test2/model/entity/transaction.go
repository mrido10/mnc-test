package entity

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID          uuid.UUID
	TransactionType string
	Amount          float64
	Remarks         string
	BalanceBefore   float64
	BalanceAfter    float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
