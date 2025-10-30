package domain

import "context"

type ShoppingCart struct {
	ID int64

	Quantity int

	UserID    int64 // this user is buying
	ProductID int64 // this
}

type ShoppingCartRepository interface {
	Create(ctx context.Context, sc *ShoppingCart) error
}
