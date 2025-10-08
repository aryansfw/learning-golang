package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"spendime/internal/auth"
	"spendime/internal/config"
	"spendime/internal/db"
	"spendime/internal/middleware"
	"spendime/internal/response"
	"spendime/internal/transaction"
	"spendime/internal/user"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	gormDb, err := db.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("Error loading database: %v", err)
	}

	authHandler := auth.NewHandler(auth.NewService(user.NewRepository(gormDb)), cfg)
	transactionHandler := transaction.NewHandler(transaction.NewService(transaction.NewRepository(gormDb)))

	mux := http.NewServeMux()

	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/register", authHandler.Register)

	mux.Handle("/transactions", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			transactionHandler.Create(w, r)
		default:
			response.Error(w, "Invalid method", http.StatusMethodNotAllowed, "invalid method")
			return
		}
	})))

	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		log.Printf("Starting server on %s", cfg.Addr())

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error processing server: %v", err)
		}
	}()

	// gracefully shutdown
	<-ctx.Done()
	stop()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("Error shutting down server gracefully: %v", err)
	}

	if err := db.Close(gormDb); err != nil {
		log.Fatalf("Error closing database connection: %v", err)
	}
}
