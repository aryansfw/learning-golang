package transaction

import (
	"spendime/internal/category"
	"time"

	"github.com/google/uuid"
)

type Type = string

const (
	TypeIncome  Type = "income"
	TypeExpense Type = "expense"
)

type Transaction struct {
	ID         uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid()" json:"id,omitempty"`
	Name       string             `gorm:"not null" json:"name,omitempty"`
	Amount     int64              `gorm:"not null" json:"amount,omitempty"`
	Date       *time.Time         `gorm:"not null" json:"date,omitempty"`
	Type       Type               `gorm:"type:varchar(20);not null" json:"type,omitempty"`
	CategoryID *uuid.UUID         `gorm:"type:uuid" json:"category_id,omitempty"`
	Category   *category.Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category,omitempty"`
	UserID     uuid.UUID          `gorm:"type:uuid" json:"-"`
	CreatedAt  time.Time          `json:"-"`
	UpdatedAt  time.Time          `json:"-"`
}

func IsValidType(t string) bool {
	switch Type(t) {
	case TypeExpense, TypeIncome:
		return true
	}
	return false
}
