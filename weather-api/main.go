package main

import (
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

func main() {
	c := cache.New(5*time.Minute, 10*time.Minute)
	mux := http.NewServeMux()

	mux.HandleFunc("/test", WeatherHandler(c))

	log.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalln("Error starting server")
	}
}
