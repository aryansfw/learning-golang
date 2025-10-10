package transaction

import (
	"context"
	"spendime/internal/contextkey"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(ctx context.Context, transaction *Transaction) error {
	if err := gorm.G[Transaction](r.db).Create(ctx, transaction); err != nil {
		return err
	}
	return nil
}

func (r *Repository) List(ctx context.Context, filters TransactionFilters) (*[]Transaction, error) {
	q := gorm.G[Transaction](r.db).Where("user_id = ?", ctx.Value(contextkey.UserID))

	if filters.Category != nil && *filters.Category {
		q = q.Preload("Category", nil)
	}
	if filters.DateFrom != nil {
		q = q.Where("date >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		q = q.Where("date <= ?", *filters.DateTo)
	}
	transactions, err := q.Find(ctx)

	if err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[Transaction](r.db).Where("id = ?", id).Delete(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, transaction *Transaction) error {
	_, err := gorm.G[Transaction](r.db).Where("id = ?", transaction.ID).Updates(ctx, *transaction)

	if err != nil {
		return err
	}

	return nil
}
