package controller

// Repository provides a database access layer for the metadata service.
type Repository interface {
}

// Controller implements business logic to retrieve movie metadata.
type Controller struct {
	repo Repository
}

// New returns a new [Controller].
func New(r Repository) *Controller {
	return &Controller{repo: r}
}
