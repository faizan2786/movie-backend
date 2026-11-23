package grpc

import (
	"context"
	"errors"
	"log"

	pb "github.com/faizan2786/movie-backend/gen/metadatapb"
	"github.com/faizan2786/movie-backend/metadata/internal/controller"
	"github.com/faizan2786/movie-backend/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MetadataController is a controller layer for metadata service.
type MetadataController interface {
	// GetMetadata returns [model.Metadata] for a movie with given id.
	//
	// It returns [controller.ErrNotFound] if a movie doesn't exist.
	GetMetadata(ctx context.Context, id string) (model.Metadata, error)

	// PutMetadata creates or replaces metadata for a movie with given id,
	// and reports whether a new record was created.
	//
	// It returns [controller.ErrNullId] if the id is empty, [controller.ErrNullField] if the title or release year is empty,
	// and [controller.ErrImmutableField] if the title or release year differs from the stored metadata.
	PutMetadata(ctx context.Context, data model.Metadata) (created bool, err error)
}

// MetadataServer is a grpc server for metadata service.
type MetadataServer struct {
	ctlr MetadataController
	pb.UnimplementedMetadataServiceServer
}

func NewServer(mc MetadataController) *MetadataServer {
	return &MetadataServer{ctlr: mc}
}

func (s *MetadataServer) GetMetadata(ctx context.Context, req *pb.GetMetadataRequest) (*pb.GetMetadataResponse, error) {
	id := req.GetMovieId()
	if len(id) == 0 {
		return nil, status.Error(codes.InvalidArgument, "movie id cannot be empty")
	}

	log.Printf("GetMetadata request: movie_id=%s", id)

	data, err := s.ctlr.GetMetadata(ctx, id)
	if errors.Is(err, controller.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "movie not found")
	}
	if err != nil {
		log.Printf("GetMetadata: movie_id=%s; error=%s", id, err)
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &pb.GetMetadataResponse{Metadata: toProto(data)}, nil
}

func (s *MetadataServer) PutMetadata(ctx context.Context, req *pb.PutMetadataRequest) (*pb.PutMetadataResponse, error) {
	if req.GetMetadata() == nil {
		return nil, status.Error(codes.InvalidArgument, "metadata cannot be null")
	}
	data := fromProto(req.GetMetadata())
	log.Printf("PutMetadata request: movie_id=%s", data.ID)

	_, err := s.ctlr.PutMetadata(ctx, data)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrNullId):
			return nil, status.Error(codes.InvalidArgument, "movie_id cannot be empty")
		case errors.Is(err, controller.ErrNullField):
			return nil, status.Error(codes.InvalidArgument, "title and release_year cannot be null")
		case errors.Is(err, controller.ErrImmutableField):
			return nil, status.Error(codes.InvalidArgument, "title and release_year cannot be changed")
		default:
			log.Printf("PutMetadata: data=%+v; error=%s", data, err)
			return nil, status.Error(codes.Internal, "something went wrong")
		}
	}
	return &pb.PutMetadataResponse{}, nil
}
