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

// MetadataServer is a grpc server for metadata service.
type MetadataServer struct {
	controller *controller.Controller
	pb.UnimplementedMetadataServiceServer
}

func NewServer(ctlr *controller.Controller) *MetadataServer {
	return &MetadataServer{controller: ctlr}
}

func (s *MetadataServer) GetMetadata(ctx context.Context, req *pb.GetMetadataRequest) (*pb.GetMetadataResponse, error) {
	id := req.GetMovieId()
	if len(id) == 0 {
		return nil, status.Error(codes.InvalidArgument, "movie id cannot be empty")
	}

	log.Printf("GetMetadata request: movie_id=%s", id)

	data, err := s.controller.GetMetadata(ctx, id)
	if errors.Is(err, controller.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "movie not found")
	}
	if err != nil {
		log.Printf("GetMetadata: movie_id=%s; error=%s", id, err)
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &pb.GetMetadataResponse{Metadata: model.ToProto(data)}, nil
}

func (s *MetadataServer) PutMetadata(ctx context.Context, req *pb.PutMetadataRequest) (*pb.PutMetadataResponse, error) {
	if req.GetMetadata() == nil {
		return nil, status.Error(codes.InvalidArgument, "metadata cannot be null")
	}
	data := model.FromProto(req.GetMetadata())
	log.Printf("PutMetadata request: movie_id=%s", data.ID)

	_, err := s.controller.PutMetadata(ctx, data)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrNullId):
			return nil, status.Error(codes.InvalidArgument, "movie id cannot be empty")
		case errors.Is(err, controller.ErrNullField):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, controller.ErrImmutableField):
			return nil, status.Error(codes.InvalidArgument, "title and release_year cannot be changed")
		default:
			log.Printf("PutMetadata: data=%+v; error=%s", data, err)
			return nil, status.Error(codes.Internal, "something went wrong")
		}
	}
	return &pb.PutMetadataResponse{}, nil
}
