package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	config := InitConfig()

	db, err := sql.Open(config.DBDriver, config.DSN())
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	repository := NewBlogRepository(db)

	service := NewBlogService(repository)

	handler := NewBlogHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /blogs", handler.HandleGetBlogs)
	mux.HandleFunc("POST /blogs", handler.HandleCreateBlog)
	mux.HandleFunc("PATCH /blogs/{id}", handler.HandleUpdateBlog)
	mux.HandleFunc("PUT /blogs/{id}", handler.HandleUpdateBlog)
	mux.HandleFunc("DELETE /blogs/{id}", handler.HandleDeleteBlog)
	mux.HandleFunc("GET /blogs/{id}", handler.HandleGetBlogById)

	http.Handle("/api/", http.StripPrefix("/api", mux))

	log.Printf("Server started at %s:%d", config.Host, config.Port)
	if err := http.ListenAndServe(config.Address(), nil); err != nil {
		log.Println("Server failed to start")
	}
}
