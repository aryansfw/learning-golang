package category

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) List(ctx context.Context) (*[]Category, error) {
	categories, err := gorm.G[Category](r.db).Find(ctx)
	if err != nil {
		return nil, err
	}

	return &categories, nil
}
