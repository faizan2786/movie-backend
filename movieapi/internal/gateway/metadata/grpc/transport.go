package grpc

import (
	pb "github.com/faizan2786/movie-backend/gen/metadatapb"
	"github.com/faizan2786/movie-backend/movieapi/pkg/model"
)

func fromProto(msg *pb.Metadata) model.Metadata {
	return model.Metadata{
		ID:          model.ID(msg.GetMovieId()),
		Title:       msg.GetTitle(),
		ReleaseYear: msg.GetReleaseYear(),
		Description: msg.GetDescription(),
	}
}
