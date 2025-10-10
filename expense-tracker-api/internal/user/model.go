package user

import (
	"spendime/internal/transaction"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID                 `gorm:"type:uuid;default:gen_random_uuid()" json:"id,omitempty"`
	Name         string                    `gorm:"not null" json:"name,omitempty"`
	Email        string                    `gorm:"unique;not null" json:"email,omitempty"`
	Password     string                    `gorm:"not null" json:"-"`
	Transactions []transaction.Transaction `gorm:"constraint:OnDelete:CASCADE" json:"transactions,omitempty"`
	CreatedAt    time.Time                 `json:"-"`
	UpdatedAt    time.Time                 `json:"-"`
}
