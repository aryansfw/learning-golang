package transaction

import (
	"spendime/internal/user"
	"time"

	"github.com/google/uuid"
)

type Type = string

const (
	TypeIncome  Type = "income"
	TypeExpense Type = "expense"
)

type Transaction struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid()"`
	Name       string     `gorm:"not null"`
	Amount     int64      `gorm:"not null"`
	Type       Type       `gorm:"type:varchar(20);not null"`
	CategoryID *uuid.UUID `gorm:"type:uuid"`
	UserID     *uuid.UUID `gorm:"type:uuid"`
	User       user.User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
