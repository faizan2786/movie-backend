package grpc

import (
	"github.com/faizan2786/movie-backend/metadata/pkg/model"

	pb "github.com/faizan2786/movie-backend/gen/metadatapb"
)

func toProto(data model.Metadata) *pb.Metadata {
	return &pb.Metadata{
		MovieId:     data.ID,
		Title:       data.Title,
		ReleaseYear: data.ReleaseYear,
		Description: data.Description,
	}
}

func fromProto(proto *pb.Metadata) model.Metadata {
	return model.Metadata{
		ID:          proto.MovieId,
		Title:       proto.Title,
		ReleaseYear: proto.ReleaseYear,
		Description: proto.Description,
	}
}
