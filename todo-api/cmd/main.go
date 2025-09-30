package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"todo/internal/config"
	"todo/internal/db"
	"todo/internal/handlers"
	"todo/internal/middleware"
	"todo/internal/repository"
	"todo/internal/services"

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

	db, err := db.Connect(cfg)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	authHandler := handlers.NewAuthHandler(services.NewAuthService(repository.NewUserRepository(db)), cfg)
	taskHandler := handlers.NewTaskHandler(services.NewTaskService(repository.NewTaskRepository(db)))

	mux := http.NewServeMux()

	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/register", authHandler.Register)

	mux.Handle("/tasks", middleware.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.ListTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		default:
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
	})))

	mux.Handle("/tasks/", middleware.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
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

	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	// <-timeoutCtx.Done()

	// if timeoutCtx.Err() == context.DeadlineExceeded {
	// 	log.Fatalln("Error timeout exceeded, shutting down")
	// }
}
