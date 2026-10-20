package memory_test

import (
	"errors"
	"testing"

	"github.com/faizan2786/movie-backend/rating/internal/repository"
	"github.com/faizan2786/movie-backend/rating/internal/repository/memory"
	"github.com/faizan2786/movie-backend/rating/pkg/model"
)

func TestMemory_EmptyDB(t *testing.T) {

	db := memory.NewDB()

	recordType := model.RecordType("movies")
	recordID := model.RecordID("1")

	// verify correct error is thrown (instead of crashing)
	s, c, err := db.GetRating(t.Context(), recordType, recordID)
	if err == nil {
		t.Fatalf("got no error; want error = %s", repository.ErrNotFound)
	}
	if s != model.RatingSum(0) {
		t.Fatalf("got sum = %d; want sum = %d", s, 0)
	}
	if c != model.RatingCount(0) {
		t.Fatalf("got count = %d; want count = %d", c, 0)
	}
}

func TestMemory_InsertAndFetch(t *testing.T) {

	db := memory.NewDB()

	recordType := model.RecordType("movies")
	recordID := model.RecordID("1")
	userID := model.UserID("123")
	ratingValue := model.RatingValue(8)

	// insert a rating
	r := model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err := db.UpsertRating(t.Context(), r)
	if err != nil {
		t.Fatalf("UpsertRating: got error = %s; want no errors", err)
	}

	// verify the stored rating
	s, c, err := db.GetRating(t.Context(), recordType, recordID)
	if err != nil {
		t.Fatalf("GetRating: got error = %s; want no errors", err)
	}
	if s != model.RatingSum(ratingValue) {
		t.Fatalf("got sum = %d; want sum = %d", s, ratingValue)
	}
	if c != model.RatingCount(1) {
		t.Fatalf("got count = %d; want count = %d", c, 1)
	}

	// insert a new rating (i.e. rating by a different user)
	userID = model.UserID("456")
	ratingValue = model.RatingValue(5)

	r = model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err = db.UpsertRating(t.Context(), r)
	if err != nil {
		t.Fatalf("UpsertRating: got error = %s; want no errors", err)
	}

	// verify stored results
	s, c, err = db.GetRating(t.Context(), recordType, recordID)
	if err != nil {
		t.Fatalf("GetRating: got error = %s; want no errors", err)
	}
	if s != model.RatingSum(13) {
		t.Fatalf("got sum = %d; want sum = %d", s, 13)
	}
	if c != model.RatingCount(2) {
		t.Fatalf("got count = %d; want count = %d", c, 2)
	}

}

func TestMemory_UpsertAndFetch(t *testing.T) {

	db := memory.NewDB()

	// seed some data
	recordType := model.RecordType("movies")
	recordID := model.RecordID("1")
	userID := model.UserID("123")
	ratingValue := model.RatingValue(8)

	r := model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err := db.UpsertRating(t.Context(), r)
	if err != nil {
		t.Fatalf("UpsertRating: got error = %s; want no errors", err)
	}

	userID = model.UserID("456")
	ratingValue = model.RatingValue(6)

	r = model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err = db.UpsertRating(t.Context(), r)
	if err != nil {
		t.Fatalf("UpsertRating: got error = %s; want no errors", err)
	}

	// verify the stored rating
	s, c, err := db.GetRating(t.Context(), recordType, recordID)
	if err != nil {
		t.Fatalf("GetRating: got error = %s; want no errors", err)
	}
	if s != model.RatingSum(14) {
		t.Fatalf("got sum = %d; want sum = %d", s, 14)
	}
	if c != model.RatingCount(2) {
		t.Fatalf("got count = %d; want count = %d", c, 2)
	}

	// update an existing rating record
	userID = model.UserID("123")
	ratingValue = model.RatingValue(5)

	r = model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err = db.UpsertRating(t.Context(), r)
	if err != nil {
		t.Fatalf("UpsertRating: got error = %s; want no errors", err)
	}

	s, c, err = db.GetRating(t.Context(), recordType, recordID)
	if err != nil {
		t.Fatalf("GetRating: got error = %s; want no errors", err)
	}
	if s != model.RatingSum(11) {
		t.Fatalf("got sum = %d; want sum = %d", s, 11)
	}
	if c != model.RatingCount(2) {
		t.Fatalf("got count = %d; want count = %d", c, 2)
	}
}

func TestMemory_UpsertInvalidInputs(t *testing.T) {

	db := memory.NewDB()

	// out of range rating value
	recordType := model.RecordType("movies")
	recordID := model.RecordID("1")
	userID := model.UserID("123")
	ratingValue := model.RatingValue(15)

	r := model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err := db.UpsertRating(t.Context(), r)
	if !errors.Is(err, repository.ErrInvalidValue) {
		t.Fatalf("got error = %s; want error = %s", err, repository.ErrInvalidValue)
	}

	// empty record type
	recordType = model.RecordType("movies")
	recordID = model.RecordID("")
	userID = model.UserID("123")
	ratingValue = model.RatingValue(8)

	r = model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err = db.UpsertRating(t.Context(), r)
	if !errors.Is(err, repository.ErrEmptyID) {
		t.Fatalf("got error = %s; want error = %s", err, repository.ErrEmptyID)
	}

	// empty recordID
	recordType = model.RecordType("movies")
	recordID = model.RecordID("")
	userID = model.UserID("123")

	r = model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err = db.UpsertRating(t.Context(), r)
	if !errors.Is(err, repository.ErrEmptyID) {
		t.Fatalf("got error = %s; want error = %s", err, repository.ErrEmptyID)
	}

	// empty userID
	recordType = model.RecordType("movies")
	recordID = model.RecordID("1")
	userID = model.UserID("")

	r = model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err = db.UpsertRating(t.Context(), r)
	if !errors.Is(err, repository.ErrEmptyID) {
		t.Fatalf("got error = %s; want error = %s", err, repository.ErrEmptyID)
	}
}

func TestMemory_GetNotFound(t *testing.T) {

	db := memory.NewDB()

	// seed some data
	recordType := model.RecordType("movies")
	recordID := model.RecordID("1")
	userID := model.UserID("123")
	ratingValue := model.RatingValue(8)

	r := model.Rating{
		RecordType: recordType,
		RecordID:   recordID,
		UserID:     userID,
		Value:      ratingValue,
	}

	err := db.UpsertRating(t.Context(), r)
	if err != nil {
		t.Fatalf("UpsertRating: got error = %s; want no errors", err)
	}

	// empty record type
	recordType = model.RecordType("")

	_, _, err = db.GetRating(t.Context(), recordType, recordID)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("got error = %s; want error = %s", err, repository.ErrNotFound)
	}

	// non-existent record type
	recordType = model.RecordType("actors")

	_, _, err = db.GetRating(t.Context(), recordType, recordID)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("got error = %s; want error = %s", err, repository.ErrNotFound)
	}

	// empty record ID
	recordType = model.RecordType("movies")
	recordID = model.RecordID("")

	_, _, err = db.GetRating(t.Context(), recordType, recordID)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("got error = %s; want error = %s", err, repository.ErrNotFound)
	}

	// non-existent record ID
	recordType = model.RecordType("movies")
	recordID = model.RecordID("456")

	_, _, err = db.GetRating(t.Context(), recordType, recordID)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("got error = %s; want error = %s", err, repository.ErrNotFound)
	}
}
