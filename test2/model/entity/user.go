package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"unique"`
	Pin         string `gorm:"not null"`
	Address     string
	Balance     float64 `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
