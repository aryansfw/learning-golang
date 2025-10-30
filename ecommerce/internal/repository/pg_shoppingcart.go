package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgShoppingCartRepository struct {
	conn *pgxpool.Pool
}

func (p *pgShoppingCartRepository) Create(ctx context.Context, sc *domain.ShoppingCart) error {
	query := `
		INSERT INTO shopping_carts
			(user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		RETURNING id`

	return p.conn.QueryRow(
		ctx,
		query,
		sc.UserID,
		sc.ProductID,
		sc.Quantity,
	).Scan(&sc.ID)
}
