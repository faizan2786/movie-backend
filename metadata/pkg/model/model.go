package model

type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseYear"`
	Description string `json:"description,omitempty"`
}
