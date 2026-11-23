package grpc

import (
	"errors"
	"fmt"

	pb "github.com/faizan2786/movie-backend/gen/ratingpb"

	"github.com/faizan2786/movie-backend/rating/pkg/model"
)

func recordTypeFromProto(pbRecordType pb.RecordType) (model.RecordType, error) {
	switch pbRecordType {
	case pb.RecordType_RECORD_TYPE_UNSPECIFIED:
		return "", errors.New("unspecified record_type")
	case pb.RecordType_RECORD_TYPE_MOVIE:
		return model.RecordTypeMovie, nil
	default:
		return "", fmt.Errorf("unknown record_type: %s; supported record type(s) are: [%s]", pbRecordType, pb.RecordType_RECORD_TYPE_MOVIE)
	}
}

// fromProto converts a [pb.PutRatingRequest] type to a [model.Rating] type.
func fromProto(msg *pb.PutRatingRequest) (model.Rating, error) {
	recordType, err := recordTypeFromProto(msg.GetRecordType())
	if err != nil {
		return model.Rating{}, err
	}

	return model.Rating{
		RecordType: recordType,
		RecordID:   model.RecordID(msg.GetRecordId()),
		UserID:     model.UserID(msg.GetUserId()),
		Value:      model.RatingValue(msg.GetRatingValue()),
	}, nil
}
