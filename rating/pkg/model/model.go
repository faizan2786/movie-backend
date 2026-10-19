package model

// RecordType is the type of a record being rated.
// Together with [RecordId], it uniquely identifies the record being rated.
//
// Example: Movie, Actor, etc.
type RecordType string

// RecordId is the ID of a record being rated.
// Together with [RecordType], it uniquely identifies the record being rated.
type RecordId string

// UserId is the ID of the user to whom the rating belongs.
type UserID string

// RatingValue represents a value for a rating record.
type RatingValue int

// Rating is a user's rating for a given record.
// Its RecordType and RecordId together identify the record being rated.
type Rating struct {
	RecordType RecordType  `json:"recordType"`
	RecordId   RecordId    `json:"recordId"`
	UserID     UserID      `json:"userID"`
	Value      RatingValue `json:"value"`
}
