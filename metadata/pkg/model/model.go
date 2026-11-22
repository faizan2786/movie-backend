package model

import (
	pb "github.com/faizan2786/movie-backend/gen/metadatapb"
)

type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseYear"`
	Description string `json:"description,omitempty"`
}

func ToProto(data Metadata) *pb.Metadata {
	return &pb.Metadata{
		MovieId:     data.ID,
		Title:       data.Title,
		ReleaseYear: int32(data.ReleaseYear),
		Description: data.Description,
	}
}

func FromProto(proto *pb.Metadata) Metadata {
	return Metadata{
		ID:          proto.MovieId,
		Title:       proto.Title,
		ReleaseYear: int(proto.ReleaseYear),
		Description: proto.Description,
	}
}
