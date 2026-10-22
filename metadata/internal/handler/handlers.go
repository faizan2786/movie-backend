package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/faizan2786/movie-backend/metadata/internal/controller"
	"github.com/faizan2786/movie-backend/metadata/pkg/model"
)

// Handler implements http handlers for the metadata service.
type Handler struct {
	c *controller.Controller
}

// New returns a new of [Handler].
func New(ctlr *controller.Controller) *Handler {
	return &Handler{c: ctlr}
}

// GetMetadata is an http GET handler to retrieve movie metadata
func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	// read `id` from the (named) path parameter
	id := r.PathValue("id")
	if len(id) == 0 {
		log.Println("GetMetadata: id can not be null")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("GetMetadata: request received for id =", id)

	// get metadata from the controller
	data, err := h.c.GetMetadata(r.Context(), id)
	if err != nil {
		log.Printf("GetMetadata: id = %s; error = %s", id, err)
		if errors.Is(err, controller.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// json encode the response and write it back
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("GetMetadata: json encoding error; id = %s; error = %s", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// PutMetadata is an http PUT handler to store movie metadata
func (h *Handler) PutMetadata(w http.ResponseWriter, r *http.Request) {
	// read `id` from the path parameter
	id := r.PathValue("id")
	if len(id) == 0 {
		log.Println("PutMetadata: id can not be null")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("PuttMetadata: request received for id =", id)

	// decode metadata from the body (assuming body is in json format)
	var data model.Metadata
	err := json.NewDecoder(r.Body).Decode(&data)
	data.ID = id // use the id given as the path parameter
	if err != nil {
		log.Printf("PutMetadata: json body decoding error; id = %s; error = %s", id, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// store the metadata
	err = h.c.PutMetadata(r.Context(), data)
	if err != nil {
		log.Printf("PuttMetadata: id = %s; error = %s", id, err)
		if errors.Is(err, controller.ErrAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
