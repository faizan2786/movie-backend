package handler

import (
	"github.com/faizan2786/movie-backend/apiserver/internal/controller"
)

// Handler implements http handlers for the rating service.
type Handler struct {
	controller *controller.Controller
}

// New returns a new [Handler].
func New(c *controller.Controller) *Handler {
	return &Handler{controller: c}
}
