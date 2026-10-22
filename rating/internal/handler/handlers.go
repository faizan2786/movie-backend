package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/faizan2786/movie-backend/rating/internal/controller"
	"github.com/faizan2786/movie-backend/rating/pkg/model"
)

// Handler implements http handlers for the metadata service.
type Handler struct {
	controller *controller.Controller
}

// New returns a new of [Handler].
func New(c *controller.Controller) *Handler {
	return &Handler{controller: c}
}

func (h *Handler) GetRating(w http.ResponseWriter, r *http.Request) {
	recordType := r.PathValue("type")
	recordID := r.PathValue("id")
	rating, err := h.controller.GetRating(r.Context(), model.RecordType(recordType), model.RecordID(recordID))
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		log.Printf("GetRating: type=%q, id=%q, error=rating record doesn't exist", recordType, recordID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("GetRating: type=%q, id=%q, error=%s", recordType, recordID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rating)
	if err != nil {
		log.Printf("GetRating: type=%q, id=%q, rating=%f, error=json response encode error: %s", recordType, recordID, rating, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) PutRating(w http.ResponseWriter, r *http.Request) {
	recordType := r.PathValue("type")
	recordID := r.PathValue("id")

	var rating model.Rating
	err := json.NewDecoder(r.Body).Decode(&rating)
	if err != nil {
		log.Printf("GetRating: type=%q, id=%q, error=json body decode error", recordType, recordID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// use the record type & ID values from the path parameters
	rating.RecordType = model.RecordType(recordType)
	rating.RecordID = model.RecordID(recordID)

	if ok := rating.Value.Validate(); !ok {
		log.Printf("PutRating: type=%q, id=%q, user_id=%q, error=provided rating is not a valid", recordType, recordID, rating.UserID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.controller.PutRating(r.Context(), rating)

	if err != nil {
		log.Printf("PutRating: type=%q, id=%q, user_id=%q, error=%s", recordType, recordID, rating.UserID, err)
		if errors.Is(err, controller.ErrEmptyID) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}
