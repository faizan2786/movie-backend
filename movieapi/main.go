package main

import (
	"context"
	"log"
	"net/http"

	"github.com/faizan2786/movie-backend/movieapi/internal/controller"
	"github.com/faizan2786/movie-backend/movieapi/internal/gateway/metadata"
	"github.com/faizan2786/movie-backend/movieapi/internal/gateway/rating"
	"github.com/faizan2786/movie-backend/movieapi/internal/handler"
	"github.com/faizan2786/movie-backend/pkg/discovery"
	registry "github.com/faizan2786/movie-backend/pkg/discovery/memory"
)

const port = "8083"

func main() {

	// register an instance of metadata and rating service
	reg := registry.New()
	err := reg.Register(context.Background(), discovery.MetadataServiceName, discovery.GenerateInstanceID(), "localhost:8081")
	if err != nil {
		panic(err)
	}
	err = reg.Register(context.Background(), discovery.RatingServiceName, discovery.GenerateInstanceID(), "localhost:8082")
	if err != nil {
		panic(err)
	}

	metadataGw := metadata.NewGateway(reg)
	ratingGw := rating.NewGateway(reg)
	controller := controller.New(metadataGw, ratingGw)

	handler := handler.New(controller)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies/{id}", handler.GetMovie)
	mux.HandleFunc("PUT /movies/{id}/ratings/{userID}", handler.PutMovieRating)

	// start the web server
	log.Printf("starting movieapi service on port: %s\n", port)
	err = http.ListenAndServe("localhost:"+port, mux)
	if err != nil {
		log.Fatalf("movieapi service stopped")
	}
}
