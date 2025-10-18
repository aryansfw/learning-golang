package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"url-shortener/internal/api"
	"url-shortener/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	_ = godotenv.Load()

	db, err := database.NewDatabasePool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}

	api := api.NewAPI(ctx, db)
	srv := api.Server()

	go func() { _ = srv.ListenAndServe() }()

	<-ctx.Done()

	_ = srv.Shutdown(ctx)
	log.Println("Server shutdown successfully")
}
