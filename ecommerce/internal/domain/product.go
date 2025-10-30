package domain

import "context"

type Product struct {
	ID    int64
	Name  string
	Price int
	Stock int

	UserID int64 // Who this product belongs to
}

type ProductRepository interface {
	GetByID(ctx context.Context, productID int64) (Product, error)

	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error

	Delete(ctx context.Context, product Product) error
}
