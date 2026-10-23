package main

import (
	"log"
	"net/http"
)

const port = "8083"

func main() {

	// db := memory.NewDB()
	// controller := controller.New(db)
	// handler := handler.New(controller)

	mux := http.NewServeMux()

	// mux.HandleFunc("GET /{type}/{id}", handler.GetRating)
	// mux.HandleFunc("PUT /{type}/{id}", handler.PutRating)

	// start the web server
	log.Printf("starting the apiserver on port: %s\n", port)
	err := http.ListenAndServe("localhost:"+port, mux)
	if err != nil {
		log.Fatalf("apiserver stopped")
	}
}
