package memory_test

// These tests are meant to be run with the race detector:
//
//	go test -race ./internal/repository/memory/

import (
	"strconv"
	"sync"
	"testing"

	"github.com/faizan2786/movie-backend/rating/internal/repository/memory"
	"github.com/faizan2786/movie-backend/rating/pkg/model"
)

// TestMemory_ConcurrentInsertsAndGets inserts ratings from many different users
// while concurrently reading the same record, then verifies the final sum and count.
func TestMemory_ConcurrentInsertsAndGets(t *testing.T) {
	db := memory.NewDB()

	recordType := model.RecordType("movies")
	recordID := model.RecordID("1")
	const users = 200

	var wg sync.WaitGroup
	var wantSum model.RatingSum
	for i := range users {
		ratingValue := model.RatingValue(i%10 + 1) // values 1..10
		wantSum += model.RatingSum(ratingValue)

		// writer: a new rating by a different user
		wg.Go(func() {
			r := model.Rating{
				RecordType: recordType,
				RecordID:   recordID,
				UserID:     model.UserID(strconv.Itoa(i)),
				Value:      ratingValue,
			}

			// use t.Errorf (not t.Fatalf) outside the test goroutine
			if err := db.UpsertRating(t.Context(), r); err != nil {
				t.Errorf("UpsertRating: user = %d; got error = %s; want no errors", i, err)
			}
		})

		// reader: runs alongside the writers (the result depends on timing, so it's not checked)
		wg.Go(func() {
			db.GetRating(t.Context(), recordType, recordID)
		})
	}
	wg.Wait()

	// verify the final state once all goroutines have finished
	s, c, err := db.GetRating(t.Context(), recordType, recordID)
	if err != nil {
		t.Fatalf("GetRating: got error = %s; want no errors", err)
	}
	if s != wantSum {
		t.Fatalf("got sum = %d; want sum = %d", s, wantSum)
	}
	if c != model.RatingCount(users) {
		t.Fatalf("got count = %d; want count = %d", c, users)
	}
}

// TestMemory_ConcurrentUpdatesSameUser updates the same user's rating from many goroutines,
// then verifies that only one rating is stored and the sum matches a single written value.
func TestMemory_ConcurrentUpdatesSameUser(t *testing.T) {
	db := memory.NewDB()

	recordType := model.RecordType("movies")
	recordID := model.RecordID("1")
	userID := model.UserID("123")
	const updates = 200

	// create a pool of rating values to write
	ratingPool := []model.RatingValue{4, 5, 7, 9}
	n := len(ratingPool)
	ratingSet := map[model.RatingValue]bool{}
	for _, v := range ratingPool {
		ratingSet[v] = true
	}

	var wg sync.WaitGroup
	for i := range updates {
		ratingValue := ratingPool[i%n]

		wg.Go(func() {
			r := model.Rating{
				RecordType: recordType,
				RecordID:   recordID,
				UserID:     userID,
				Value:      ratingValue,
			}

			if err := db.UpsertRating(t.Context(), r); err != nil {
				t.Errorf("UpsertRating: value = %d; got error = %s; want no errors", ratingValue, err)
			}
		})
	}
	wg.Wait()

	// the last write wins, so the sum must equal one of the written values
	s, c, err := db.GetRating(t.Context(), recordType, recordID)
	if err != nil {
		t.Fatalf("GetRating: got error = %s; want no errors", err)
	}
	if c != model.RatingCount(1) {
		t.Fatalf("got count = %d; want count = %d", c, 1)
	}
	if !ratingSet[model.RatingValue(s)] {
		t.Fatalf("got sum = %d; want sum = one of %v", s, ratingPool)
	}
}
