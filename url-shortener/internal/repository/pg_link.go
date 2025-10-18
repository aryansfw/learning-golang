package repository

import (
	"context"
	"url-shortener/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgLinkRepository struct {
	conn *pgxpool.Pool
}

func NewPGLink(conn *pgxpool.Pool) domain.LinkRepository {
	return &pgLinkRepository{conn: conn}
}

func (p *pgLinkRepository) GetByShortURL(ctx context.Context, url string) (*domain.Link, error) {
	query := `
	SELECT id, short_url, long_url
	FROM links
	WHERE short_url = $1
	`

	var link domain.Link

	if err := p.conn.QueryRow(ctx, query, url).Scan(&link.ID, &link.ShortURL, &link.LongURL); err != nil {
		return nil, err
	}

	return &link, nil
}

func (p *pgLinkRepository) Create(ctx context.Context, link *domain.Link) error {
	query := `
	INSERT INTO links 
		(short_url, long_url)
	VALUES ($1, $2)
	RETURNING id`

	return p.conn.QueryRow(ctx, query, link.ShortURL, link.LongURL).Scan(&link.ID)
}
