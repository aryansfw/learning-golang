package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgUserRepository struct {
	conn *pgxpool.Pool
}

func NewPGUser(conn *pgxpool.Pool) domain.UserRepository {
	return &pgUserRepository{conn: conn}
}

func (p *pgUserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	query := `
		SELECT id, name, email
		FROM users
		WHERE id = $1`

	var user domain.User

	if err := p.conn.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (p *pgUserRepository) GetPasswordByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT password
		FROM users
		WHERE email = $1`

	var user domain.User

	if err := p.conn.QueryRow(ctx, query, email).Scan(&user.Password); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (p *pgUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users
			(name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id`

	return p.conn.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(&user.ID)
}
