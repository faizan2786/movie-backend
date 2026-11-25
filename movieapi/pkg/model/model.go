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

// Metadata is movie metadata.
type Metadata struct {
	ID          ID     `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int32  `json:"releaseYear"`
	Description string `json:"description,omitempty"`
}

// Rating is a movie rating submitted by a user.
type Rating struct {
	UserID ID    `json:"userID"`
	Value  uint8 `json:"value"`
}

// Movie is the metadata and rating for a movie.
type Movie struct {
	Metadata
	Rating float64 `json:"rating,omitzero"` // in json/v2, we must use "omitzero" to omit empty numeric or bool fields
}
