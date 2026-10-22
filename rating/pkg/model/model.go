package model

// RecordType is the type of a record being rated.
// Together with [RecordID], it uniquely identifies the record being rated.
//
// Example: Movie, Actor, etc.
type RecordType string

// RecordID is the ID of a record being rated.
// Together with [RecordType], it uniquely identifies the record being rated.
type RecordID string

// UserID is the ID of the user to whom the rating belongs.
type UserID string

// RatingValue represents a value for a rating record.
type RatingValue int

const (
	// MinRatingValue is the minimum allowed [RatingValue].
	MinRatingValue RatingValue = 1

	// MaxRatingValue is the maximum allowed [RatingValue].
	MaxRatingValue RatingValue = 10
)

// Validate returns true if the value is valid.
func (rv RatingValue) Validate() bool {
	if rv < MinRatingValue || rv > MaxRatingValue {
		return false
	}
	return true
}

// Rating is a user's rating for a given record.
// Its RecordType and RecordID together uniquely identify the record being rated.
type Rating struct {
	RecordType RecordType  `json:"recordType"`
	RecordID   RecordID    `json:"recordID"`
	UserID     UserID      `json:"userID"`
	Value      RatingValue `json:"value"`
}
