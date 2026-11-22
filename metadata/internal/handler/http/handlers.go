package http

import (
	"encoding/json/v2"
	"errors"
	"log"
	"net/http"

	"github.com/faizan2786/movie-backend/metadata/internal/controller"
	"github.com/faizan2786/movie-backend/metadata/pkg/model"
)

// Handler implements http handlers for the metadata service.
type HttpHandler struct {
	c *controller.Controller
}

// New returns a new of [Handler].
func New(ctlr *controller.Controller) *HttpHandler {
	return &HttpHandler{c: ctlr}
}

// GetMetadata is an http GET handler to retrieve movie metadata
func (h *HttpHandler) GetMetadata(w http.ResponseWriter, r *http.Request) {
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
	err = json.MarshalWrite(w, data)
	if err != nil {
		log.Printf("GetMetadata: json encoding error; id = %s; error = %s", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// PutMetadata is an http PUT handler to store movie metadata
func (h *HttpHandler) PutMetadata(w http.ResponseWriter, r *http.Request) {
	// read `id` from the path parameter
	id := r.PathValue("id")
	if len(id) == 0 {
		log.Println("PutMetadata: id can not be null")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("PutMetadata: request received for id =", id)

	// decode metadata from the body (assuming body is in json format)
	var data model.Metadata
	err := json.UnmarshalRead(r.Body, &data)
	data.ID = id // use the id given as the path parameter
	if err != nil {
		log.Printf("PutMetadata: json body decoding error; id = %s; error = %s", id, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// store the metadata
	created, err := h.c.PutMetadata(r.Context(), data)
	if err != nil {
		log.Printf("PutMetadata: id = %s; error = %s", id, err)
		switch {
		case errors.Is(err, controller.ErrNullField):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, controller.ErrImmutableField):
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// return 201 if a new record was created
	// (in any other case, 200 will be set by default).
	if created {
		w.WriteHeader(http.StatusCreated)
	}
}
