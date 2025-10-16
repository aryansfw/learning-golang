package main

import (
	"context"
	"log"
	"md-note-api/db"
	"md-note-api/handler"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	defer dbConn.Close()

	noteHandler := handler.NewNote(dbConn)
	http.HandleFunc("POST /notes", noteHandler.Upload)
	http.HandleFunc("GET /notes/{id}", noteHandler.Download)
	http.HandleFunc("GET /notes/{id}/html", noteHandler.GetHTML)

	srv := &http.Server{
		Addr: ":8080",
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		log.Println("Started server at :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err.Error())
		}
	}()

	<-ctx.Done()
	stop()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// stop server
	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("Error shutting down server: %v", err.Error())
	}
	log.Println("Successfully shutdown")
}
