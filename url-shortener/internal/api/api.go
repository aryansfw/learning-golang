package api

import (
	"context"
	"net/http"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type api struct {
	linkRepo domain.LinkRepository
}

func NewAPI(ctx context.Context, pool *pgxpool.Pool) *api {
	linkRepo := repository.NewPGLink(pool)
	return &api{
		linkRepo: linkRepo,
	}
}

func (a *api) Server() *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /r/{url}", a.redirectLinkHandler)

	mux.HandleFunc("POST /link", a.createLinkHandler)
	return mux
}
