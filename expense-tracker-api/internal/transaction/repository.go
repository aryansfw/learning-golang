package transaction

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(transaction *Transaction) error {
	ctx := context.Background()
	if err := gorm.G[Transaction](r.db).Create(ctx, transaction); err != nil {
		return err
	}
	return nil
}

func (r *Repository) List(userID uuid.UUID) (*[]Transaction, error) {
	ctx := context.Background()
	transactions, err := gorm.G[Transaction](r.db).Where("id = ?", userID).Find(ctx)

	if err != nil {
		return nil, err
	}

	return &transactions, nil
}
