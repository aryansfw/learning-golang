package main

import (
	// "context"
	// "encoding/json"
	// "fmt"
	"context"
	"ecommerce/internal/api"
	"log"
	"os"
	"os/signal"
	"syscall"

	// "os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	// xendit "github.com/xendit/xendit-go/v7"
)

func main() {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}
	defer db.Close()

	api := api.NewAPI(ctx, db)
	srv := api.Server(8080)

	go func() {
		log.Println("Server started at :8080")
		_ = srv.ListenAndServe()
	}()

	<-ctx.Done()

	_ = srv.Shutdown(ctx)
	log.Println("Server shutdown successfully")
	// xnd := xendit.NewClient(os.Getenv("API_KEY"))

	// response from `GetAllPaymentRequests`: PaymentRequestListResponse
	// http.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
	// 	resp, hr, err := xnd.PaymentRequestApi.GetAllPaymentRequests(context.Background()).
	// 		Execute()

	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error when calling `PaymentRequestApi.GetAllPaymentRequests``: %v\n", err.Error())

	// 		b, _ := json.Marshal(err.FullError())
	// 		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

	// 		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", hr)
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	_ = json.NewEncoder(w).Encode(resp)
	// 	fmt.Fprintf(os.Stdout, "Response from `PaymentRequestApi.GetAllPaymentRequests`: %v\n", resp)
	// })
}
