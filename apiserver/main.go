package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/faizan2786/movie-backend/apiserver/internal/gateway/rating"
	"github.com/faizan2786/movie-backend/apiserver/pkg/model"
)

const port = "8083"

func main() {

	// db := memory.NewDB()
	// controller := controller.New(db)
	// handler := handler.New(controller)

	//mux := http.NewServeMux()

	// mux.HandleFunc("GET /{type}/{id}", handler.GetRating)
	// mux.HandleFunc("PUT /{type}/{id}", handler.PutRating)

	// // start the web server
	// log.Printf("starting the apiserver on port: %s\n", port)
	// err := http.ListenAndServe("localhost:"+port, mux)
	// if err != nil {
	// 	log.Fatalf("apiserver stopped")
	// }

	addr := "http://localhost:8082"
	gw, err := rating.NewGateway(addr)
	if err != nil {
		log.Fatal(err)
	}

	// submit ratings...

	err = gw.SubmitMovieRating(
		context.Background(),
		"123",
		model.Rating{
			UserID: model.ID("1"),
			Value:  model.RatingValue(6),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = gw.SubmitMovieRating(
		context.Background(),
		"123",
		model.Rating{
			UserID: model.ID("2"),
			Value:  model.RatingValue(9),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = gw.SubmitMovieRating(
		context.Background(),
		"456",
		model.Rating{
			UserID: model.ID("1"),
			Value:  model.RatingValue(8),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// GET ratings...
	data, err := gw.GetMovieRating(context.Background(), "123")
	if err != nil {
		log.Fatal(err)
	}
	jsonData := strings.Builder{}
	json.NewEncoder(&jsonData).Encode(data)
	fmt.Print(jsonData.String())

	data, err = gw.GetMovieRating(context.Background(), "456")
	if err != nil {
		log.Fatal(err)
	}
	jsonData = strings.Builder{}
	json.NewEncoder(&jsonData).Encode(data)
	fmt.Print(jsonData.String())

	data, err = gw.GetMovieRating(context.Background(), "789")
	if err != nil {
		log.Fatal(err)
	}
	jsonData = strings.Builder{}
	json.NewEncoder(&jsonData).Encode(data)
	fmt.Print(jsonData.String())
}
