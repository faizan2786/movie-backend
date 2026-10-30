package model

import (
	"regexp"
)

var validID = regexp.MustCompile(`^[A-Za-z0-9]+(-[A-Za-z0-9]+)*$`)

// ID is a unique ID for a movie
type ID string

// Valid reports whether id contains only letters, digits, and hyphens.
func (id ID) Valid() bool {
	return validID.MatchString(string(id))
}

// status is an HTTP response status
type Status int

// Metadata is movie metadata.
type Metadata struct {
	ID          ID     `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseYear,omitempty"`
	Description string `json:"description,omitempty"`
}

// RatingValue is aggregated rating for a movie.
type RatingValue float64

// Rating is a movie rating submitted by a user.
type Rating struct {
	UserID ID          `json:"userID"`
	Value  RatingValue `json:"value"`
}

// Movie is the metadata and rating for a movie.
type Movie struct {
	Metadata
	Rating RatingValue `json:"rating,omitempty"`
}
