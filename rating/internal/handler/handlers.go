package handler

import (
	"github.com/faizan2786/movie-backend/rating/internal/controller"
)

// Handler implements http handlers for the metadata service.
type Handler struct {
	c *controller.Controller
}

// New returns a new of [Handler].
func New(ctlr *controller.Controller) *Handler {
	return &Handler{c: ctlr}
}
