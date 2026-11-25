package grpc

import (
	pb "github.com/faizan2786/movie-backend/gen/metadatapb"
	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
)

func fromProto(msg *pb.GetMetadataResponse) *model.Metadata {
	if msg.Metadata == nil {
		return nil
	}

	return &model.Metadata{
		ID:          model.ID(msg.Metadata.GetMovieId()),
		Title:       msg.Metadata.GetTitle(),
		ReleaseYear: msg.Metadata.GetReleaseYear(),
		Description: msg.Metadata.GetDescription(),
	}
}
