package main

import (
	"log"
	"net/http"

	"github.com/faizan2786/movie-backend/apiserver/internal/controller"
	"github.com/faizan2786/movie-backend/apiserver/internal/gateway/metadata"
	"github.com/faizan2786/movie-backend/apiserver/internal/gateway/rating"
	"github.com/faizan2786/movie-backend/apiserver/internal/handler"
)

const port = "8083"

func main() {

	metadataGw, err := metadata.NewGateway("http://localhost:8081")
	if err != nil {
		panic(err)
	}
	ratingGw, err := rating.NewGateway("http://localhost:8082")
	if err != nil {
		panic(err)
	}
	controller := controller.New(metadataGw, ratingGw)
	handler := handler.New(controller)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies/{id}", handler.GetMovie)
	mux.HandleFunc("PUT /ratings/movies/{id}", handler.PutMovieRating)

	// start the web server
	log.Printf("starting movie api server on port: %s\n", port)
	err = http.ListenAndServe("localhost:"+port, mux)
	if err != nil {
		log.Fatalf("api server stopped")
	}
}
