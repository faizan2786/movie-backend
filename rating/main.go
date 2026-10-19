package main

import (
	"log"
	"net/http"
)

const port = "8082"

func main() {

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
