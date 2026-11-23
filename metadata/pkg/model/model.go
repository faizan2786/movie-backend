package model

import (
	pb "github.com/faizan2786/movie-backend/gen/metadatapb"
)

type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int32  `json:"releaseYear"`
	Description string `json:"description,omitempty"`
}

func ToProto(data Metadata) *pb.Metadata {
	return &pb.Metadata{
		MovieId:     data.ID,
		Title:       data.Title,
		ReleaseYear: data.ReleaseYear,
		Description: data.Description,
	}
}

func FromProto(proto *pb.Metadata) Metadata {
	return Metadata{
		ID:          proto.MovieId,
		Title:       proto.Title,
		ReleaseYear: proto.ReleaseYear,
		Description: proto.Description,
	}
}
