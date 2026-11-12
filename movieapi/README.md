# MovieAPI Service

This Go module provides **the public HTTP API** for the movie backend. It combines data from the `metadata` and `rating` services, so clients can fetch a movie details with its average rating and submit user ratings for movies.

Its design uses a **layered architecture**, separating the HTTP handling (i.e. API layer), business logic (i.e. controller layer), and downstream service calls (i.e. gateway layer) concerns. It has no storage of its own.

## Endpoints

| Method | Path                            | Description                                |
| ------ | ------------------------------- | ------------------------------------------ |
| `GET`  | `/movies/{id}`                  | Get a movie's metadata and average rating  |
| `PUT`  | `/movies/{id}/ratings/{userID}` | Submit a user's rating for a movie         |

`GET` returns the movie's metadata with its average rating:

```json
{"id": "the-matrix", "title": "The Matrix", "releaseYear": 1999, "rating": 8.5}
```

`PUT` takes the rating value as JSON; a user's new rating for a movie replaces their previous one:

```json
{"value": 8}
```

## Behavior

- Movie and user IDs may contain only letters, digits, and single hyphens between them (e.g. `the-matrix`); other IDs return 400
- Both endpoints return 404 if the movie doesn't exist in the metadata service
- A movie with no ratings yet is returned without the `rating` field
- `PUT` returns 400 if the rating service rejects the value (it must be a whole number from 1 to 10)
- Finds `metadata` and `rating` instances through Consul on each request, picking a healthy instance at random, and calls it with a 5 second timeout and no redirects
- Returns 500 if no healthy instance of a downstream service is registered

## Startup & graceful shutdown

- On startup, the service registers itself with Consul (at `localhost:8500`) as `movieapi` with a unique instance ID, then sends a heartbeat every second to keep the instance healthy
- On `SIGINT` (Ctrl+C) or `SIGTERM`, it first deregisters the instance, so callers stop discovering it, and then shuts down gracefully, giving in-flight requests up to 5 seconds to finish before closing the remaining connections
- A second Ctrl+C during shutdown forces the service to quit
- If the server fails to start (e.g. the port is already in use), the instance is deregistered before the service exits

## Running

Start the Consul server (see the [root README](../README.md#service-discovery)) and the `metadata` and `rating` services first, then:

```sh
go run ./movieapi                # listens on localhost:8083
go run ./movieapi -port 9083     # listens on localhost:9083
go test ./movieapi/...
```

Controller tests use [gomock](https://github.com/uber-go/mock) mocks. After changing the gateway interfaces in `controller.go`, regenerate them with `go generate ./movieapi/...` (needs `go install go.uber.org/mock/mockgen@v0.6.0`).

## Project structure

```text
movieapi/
├── main.go                         # Wires the service layers and starts the HTTP server
├── internal/
│   ├── handler/                    # API layer for HTTP request handling
│   ├── controller/                 # Business logic layer combining metadata and ratings
│   │   └── mock_controller/        # Generated gateway mocks for controller tests
│   └── gateway/                    # Gateway layer for calling downstream services
│       ├── metadata/               # Gateway to the metadata service
│       └── rating/                 # Gateway to the rating service
└── pkg/
    └── model/                      # Shared data model for the service
```
