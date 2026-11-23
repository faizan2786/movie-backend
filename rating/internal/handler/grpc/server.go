package grpc

import (
	"context"
	"errors"
	"log"

	pb "github.com/faizan2786/movie-backend/gen/ratingpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/faizan2786/movie-backend/rating/internal/controller"
	"github.com/faizan2786/movie-backend/rating/pkg/model"
)

type RatingController interface {

	// GetRating returns average ratings for a given record type and record ID.
	//
	// It returns [controller.ErrNotFound] when a rating record doesn't exist.
	GetRating(ctx context.Context, recordType model.RecordType, recordID model.RecordID) (float64, error)

	// PutRating stores user rating by its record type and ID. It replaces the existing rating.
	//
	// It returns [controller.ErrEmptyID] if the record type, record ID or user ID is empty and
	// returns [controller.ErrInvalidValue] if the rating value is out of range.
	PutRating(ctx context.Context, r model.Rating) error
}

type RatingServer struct {
	ctlr RatingController
	pb.UnimplementedRatingServiceServer
}

func NewServer(rc RatingController) *RatingServer {
	return &RatingServer{ctlr: rc}
}

// GetAggregatedRating returns aggregated rating for a record.
//
// If the record hasn't been rated yet, it returns OK with the rating unset.
func (s *RatingServer) GetAggregatedRating(ctx context.Context, req *pb.GetAggregatedRatingRequest) (*pb.GetAggregatedRatingResponse, error) {
	log.Printf("GetAggregatedRating: level=%s, record_type=%s, record_id=%s, msg=%s", "info", req.GetRecordType(), req.GetRecordId(), "request received")

	recordType, err := recordTypeFromProto(req.GetRecordType())
	if err != nil {
		log.Printf("GetAggregatedRating: level=%s, record_type=%s, record_id=%s, error=%s", "debug", recordType, req.GetRecordId(), err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	recordID := req.GetRecordId()
	if len(recordID) == 0 {
		log.Printf("GetAggregatedRating: level=%s, record_type=%s, record_id=%s, msg=%s", "debug", recordType, recordID, "record id is null")
		return nil, status.Error(codes.InvalidArgument, "record_id cannot be empty")
	}

	r, err := s.ctlr.GetRating(ctx, recordType, model.RecordID(recordID))

	// a non-existing rating for a record is a valid case
	if errors.Is(err, controller.ErrNotFound) {
		log.Printf("GetAggregatedRating: level=%s, record_type=%s, record_id=%s, error=%s", "debug", recordType, recordID, err)
		return &pb.GetAggregatedRatingResponse{}, nil
	}
	if err != nil {
		log.Printf("GetAggregatedRating: level=%s, record_type=%s, record_id=%s, error=%s", "error", recordType, recordID, err)
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &pb.GetAggregatedRatingResponse{
		Rating: new(r),
	}, nil
}

// PutRating upserts a user's rating for a record.
func (s *RatingServer) PutRating(ctx context.Context, req *pb.PutRatingRequest) (*pb.PutRatingResponse, error) {
	log.Printf("PutRating: level=%s, record_type=%s, record_id=%s, msg=%s", "info", req.GetRecordType(), req.GetRecordId(), "request received")

	rating, err := fromProto(req)
	if err != nil {
		log.Printf("PutRating: level=%s, record_type=%s, record_id=%s, error=%s", "debug", req.GetRecordType(), req.GetRecordId(), err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.ctlr.PutRating(ctx, rating)
	if err != nil {
		rType, _ := recordTypeFromProto(req.GetRecordType())

		switch {
		case errors.Is(err, controller.ErrEmptyID):
			log.Printf("PutRating: level=%s, record_type=%s, record_id=%s, error=%s", "debug", rType, req.GetRecordId(), err)
			return nil, status.Error(codes.InvalidArgument, "record_id and user_id cannot be empty")
		case errors.Is(err, controller.ErrInvalidValue):
			log.Printf("PutRating: level=%s, record_type=%s, record_id=%s, error=%s", "debug", rType, req.GetRecordId(), err)
			return nil, status.Errorf(codes.InvalidArgument, "rating_value must be between %d and %d (both inclusive)", model.MinRatingValue, model.MaxRatingValue)
		default:
			log.Printf("PutRating: level=%s, record_type=%s, record_id=%s, error=%s", "error", rType, req.GetRecordId(), err)
			return nil, status.Error(codes.Internal, "something went wrong")
		}
	}

	log.Printf("PutRating: level=%s, record_type=%s, record_id=%s, msg=%s", "debug", req.GetRecordType(), req.GetRecordId(), "successful rating upsert")
	return &pb.PutRatingResponse{}, nil
}
