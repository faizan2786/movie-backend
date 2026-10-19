package memory

import (
	"context"
	"sync"

	"github.com/faizan2786/movie-backend/rating/internal/repository"
	"github.com/faizan2786/movie-backend/rating/pkg/model"
)

// DB is an in-memory database for the rating service.
// Its zero value is not ready for use, use [NewDB] to create a DB.
type DB struct {
	m sync.RWMutex

	// ratings stores user ratings by record type and ID.
	ratings map[model.RecordType]map[model.RecordID]map[model.UserID]model.Rating

	// ratingSums stores a running sum of user ratings
	// per record type and record ID
	ratingSums map[model.RecordType]map[model.RecordID]model.RatingSum
}

// NewDB returns a new [DB].
func NewDB() *DB {
	return &DB{
		ratings:    map[model.RecordType]map[model.RecordID]map[model.UserID]model.Rating{},
		ratingSums: map[model.RecordType]map[model.RecordID]model.RatingSum{},
	}
}

// GetRating returns the running sum and total count of user ratings
// for a given record type and record ID.
// It returns [repository.ErrNotFound] when a rating doesn't exist.
func (db *DB) GetRating(_ context.Context, rType model.RecordType, rID model.RecordID) (model.RatingSum, model.RatingCount, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	var sum model.RatingSum
	var count model.RatingCount
	var ok bool

	// return error if either record type or record ID don't exist
	if sum, ok = db.ratingSums[rType][rID]; !ok {
		return sum, count, repository.ErrNotFound
	}
	count = model.RatingCount(len(db.ratings[rType][rID]))

	return sum, count, nil
}

// UpsertRating inserts or updates user rating by its record type and ID.
//
// It returns [repository.ErrEmptyID] if the record type, record ID or user ID is empty and
// returns [repository.ErrInvalidValue] if the rating value is out of range.
func (db *DB) UpsertRating(_ context.Context, r model.Rating) error {
	if len(r.RecordType) == 0 || len(r.RecordID) == 0 || len(r.UserID) == 0 {
		return repository.ErrEmptyID
	}

	if r.Value < model.MinRatingValue || r.Value > model.MaxRatingValue {
		return repository.ErrInvalidValue
	}

	db.m.Lock()
	defer db.m.Unlock()

	// update - if record already exists
	if _, ok := db.ratings[r.RecordType][r.RecordID][r.UserID]; ok {
		db.updateRating(r)
		return nil
	}

	// insert - if it doesn't exist yet
	db.insertRating(r)
	return nil
}

// helper functions...

func (db *DB) updateRating(r model.Rating) {
	currValue := db.ratings[r.RecordType][r.RecordID][r.UserID].Value

	// update the user rating
	db.ratings[r.RecordType][r.RecordID][r.UserID] = r

	// update the sum for given record type and ID
	db.ratingSums[r.RecordType][r.RecordID] =
		db.ratingSums[r.RecordType][r.RecordID] - model.RatingSum(currValue) + model.RatingSum(r.Value)
}

func (db *DB) insertRating(r model.Rating) {
	if db.ratings[r.RecordType] == nil {
		db.ratings[r.RecordType] = map[model.RecordID]map[model.UserID]model.Rating{}
		db.ratingSums[r.RecordType] = map[model.RecordID]model.RatingSum{}
	}

	if db.ratings[r.RecordType][r.RecordID] == nil {
		db.ratings[r.RecordType][r.RecordID] = map[model.UserID]model.Rating{}
	}

	// store the user rating
	db.ratings[r.RecordType][r.RecordID][r.UserID] = r

	// update the sum for given record type and ID
	db.ratingSums[r.RecordType][r.RecordID] += model.RatingSum(r.Value)
}
