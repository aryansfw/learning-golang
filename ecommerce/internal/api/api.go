package api

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/repository"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	xendit "github.com/xendit/xendit-go/v7"
)

type api struct {
	xenditClient *xendit.APIClient
	userRepo     domain.UserRepository
}

func NewAPI(ctx context.Context, pool *pgxpool.Pool) *api {
	xnd := xendit.NewClient(os.Getenv("API_KEY"))
	userRepo := repository.NewPGUser(pool)

	return &api{
		xenditClient: xnd,
		userRepo:     userRepo,
	}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("POST /login", a.loginHandler)
	r.HandleFunc("POST /register", a.registerHandler)
	return r
}
