package model

type Metadata struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_year,omitempty"`
	Description string `json:"description,omitempty"`
}
