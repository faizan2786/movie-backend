package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/faizan2786/movie-backend/apiserver/internal/controller"
	"github.com/faizan2786/movie-backend/apiserver/pkg/model"
)

// Handler implements http handlers for the apiserver.
type Handler struct {
	controller *controller.Controller
}

// New returns a new [Handler].
func New(c *controller.Controller) *Handler {
	return &Handler{controller: c}
}

// GetMovie is an http GET handler to retrieve movie details.
func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// get the movie details from the controller
	movie, err := h.controller.GetMovie(r.Context(), model.ID(id))
	if errors.Is(err, controller.ErrInvalidID) {
		log.Printf("GetMovie: error=%s; id=%s", err, id)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if errors.Is(err, controller.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("GetMovie: error=%s; id=%s", err, id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send the response in json format
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		log.Printf("GetMovie: json response encode error: error=%s; id=%s", err, id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// PutMovieRating is an http PUT handler to submit movie ratings.
func (h *Handler) PutMovieRating(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("id")

	// deserialise the request body to get rating details
	var rating model.Rating
	err := json.NewDecoder(r.Body).Decode(&rating)
	if err != nil {
		log.Printf("PutMovieRating: json request decode error: error=%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// submit the rating using controller
	status, err := h.controller.SubmitMovieRating(r.Context(), model.ID(movieID), rating)
	if err != nil {
		log.Printf("PutMovieRating: error=%s; movie_id=%s; user_id=%s", err, movieID, rating.UserID)

		if errors.Is(err, controller.ErrInvalidID) {
			w.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(err, controller.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// forward the status code received from the controller
	w.WriteHeader(int(status))
}
