package main

import (
	"log"
	"net/http"

	"github.com/faizan2786/movie-backend/rating/internal/controller"
	"github.com/faizan2786/movie-backend/rating/internal/repository/memory"
)

const port = "8082"

func main() {

	var _ controller.Repository = (*memory.DB)(nil)
	// db := memory.NewDB()
	// controller := controller.New(db)
	// handler := handler.New(controller)

	mux := http.NewServeMux()

	// TO-DO: register handlers...

	// start the web server
	log.Printf("starting rating service on port: %s\n", port)
	err := http.ListenAndServe("localhost:"+port, mux)
	if err != nil {
		log.Fatalf("rating service stopped")
	}
}
