package db

import (
	"context"
	"database/sql"
	"time"
	"todo/internal/config"

	_ "github.com/lib/pq"
)

func Connect(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DSN())

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
