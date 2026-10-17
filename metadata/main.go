package main

import (
	"log"
	"net/http"

	"github.com/faizan2786/movie-backend/metadata/internal/controller"
	"github.com/faizan2786/movie-backend/metadata/internal/handler"
	"github.com/faizan2786/movie-backend/metadata/internal/repository/memory"
)

const port = "8081"

func main() {

	db := memory.NewDB()
	controller := controller.New(db)
	handler := handler.New(*controller)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{id}", handler.GetMetadata)
	mux.HandleFunc("PUT /{id}", handler.PutMetadata)

	// start the web server
	log.Printf("starting metadata service on port: %s\n", port)
	err := http.ListenAndServe("localhost:"+port, mux)
	if err != nil {
		log.Fatalf("metadata service stopped")
	}
}
