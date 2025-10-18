package domain

import "context"

// Link represents url short and long form
type Link struct {
	ID       int64
	ShortURL string
	LongURL  string
}

type LinkRepository interface {
	GetByShortURL(ctx context.Context, url string) (*Link, error)

	Create(ctx context.Context, link *Link) error
}
