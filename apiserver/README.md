# The API Server

This Go module provides the HTTP API server for the movie backend. It exposes a GET endpoint to retrieve a movie's metadata and its aggregated user ratings by its ID. It communicates with `metadata` and `rating` services to fetch the necessary movie data.