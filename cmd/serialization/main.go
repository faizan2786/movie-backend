package main

import (
	"encoding/json/v2"
	"encoding/xml"
	"fmt"

	"github.com/faizan2786/movie-backend/gen/movieapi"
	"google.golang.org/protobuf/proto"
)

type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int32  `json:"releaseYear,omitempty"`
	Description string `json:"description,omitempty"`
}

type Movie struct {
	Metadata
	RatingValue float64 `json:"rating,omitempty"`
}

func main() {

	movie := &Movie{
		Metadata{
			ID:          "123-456-789",
			Title:       "The Test",
			ReleaseYear: 2020,
			Description: "A test movie to benchmark various serialisation format",
		},
		7.5,
	}

	msgJson := serializeToJson(movie)
	msgXML := serializeToXML(movie)
	msgPb := serializeToProtobuf(toProto(movie))

	fmt.Printf("Size of JSON message: %d bytes\n", len(msgJson))
	fmt.Printf("Size of XML message:  %d bytes\n", len(msgXML))
	fmt.Printf("Size of Protobuf message:  %d bytes\n", len(msgPb))

	// fmt.Println()
	// fmt.Printf("Json msg:\n%s\n\n", msgJson)
	// fmt.Printf("XML msg:\n%s\n\n", msgXML)
	// fmt.Printf("Protobuf msg:\n%s\n", msgPb)
}

func serializeToJson(movie *Movie) []byte {
	out, err := json.Marshal(movie)
	if err != nil {
		panic(err)
	}
	return out
}

func serializeToXML(movie *Movie) []byte {
	out, err := xml.Marshal(movie)
	if err != nil {
		panic(err)
	}
	return out
}

func serializeToProtobuf(pb *movieapi.Movie) []byte {
	out, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return out
}

// convert Movie struct into a protobuf type
func toProto(movie *Movie) *movieapi.Movie {
	return &movieapi.Movie{
		Metadata: &movieapi.Metadata{
			Id:          movie.ID,
			Title:       movie.Title,
			ReleaseYear: new(movie.ReleaseYear), // pointer proto fields must store a copy
			Description: new(movie.Description),
		},
		Rating: new(movie.RatingValue),
	}
}
