package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgProductRepository struct {
	conn *pgxpool.Pool
}

func NewPGProduct(conn *pgxpool.Pool) domain.ProductRepository {
	return &pgProductRepository{conn: conn}
}

func (p *pgProductRepository) GetByID(ctx context.Context, id int64) (domain.Product, error) {
	query := `
		SELECT id, name, price, stock
		FROM products
		WHERE id = $1`

	var prd domain.Product

	if err := p.conn.QueryRow(ctx, query, id).Scan(&prd.ID, &prd.Name, &prd.Price, &prd.Stock); err != nil {
		return domain.Product{}, nil
	}

	return prd, nil
}

func (p *pgProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products
			(name, price, stock, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	return p.conn.QueryRow(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.UserID,
	).Scan(&product.ID)
}

func (p *pgProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET name = $3,
			price = $4,
			stock = $5,
		WHERE id = $1
			AND user_id = $2`

	_, err := p.conn.Exec(
		ctx,
		query,
		product.ID,
		product.UserID,
		product.Name,
		product.Price,
		product.Stock,
	)

	return err
}

func (p *pgProductRepository) Delete(ctx context.Context, product domain.Product) error {
	query := `DELETE FROM products WHERE id = $1 AND user_id = $2`

	_, err := p.conn.Exec(ctx, query, product.ID, product.UserID)
	if err != nil {
		return err
	}

	return err
}
