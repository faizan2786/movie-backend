# MovieAPI Service

This Go module provides **the public HTTP API** for the movie backend. It combines data from the `metadata` and `rating` gRPC services, so clients can fetch a movie's details with its average rating and submit user ratings for movies.

Its design uses a **layered architecture**, separating the HTTP handling (i.e. API layer), business logic (i.e. controller layer), and downstream service calls (i.e. gateway layer) concerns. It has no storage of its own.

## HTTP Endpoints

| Method | Path                            | Description                                |
| ------ | ------------------------------- | ------------------------------------------ |
| `GET`  | `/movies/{id}`                  | Get a movie's metadata and average rating  |
| `PUT`  | `/movies/{id}/ratings/{userID}` | Submit a user's rating for a movie         |

`GET` returns the movie's metadata with its average rating:

```json
{"id": "the-matrix", "title": "The Matrix", "releaseYear": 1999, "description": "A hacker learns the truth about reality.", "rating": 8.5}
```

`PUT` takes the rating value as JSON; a user's new rating for a movie replaces their previous one:

```json
{"value": 8}
```

## Behavior

- Movie and user IDs may contain only letters, digits, and single hyphens between them (e.g. `the-matrix`, `123-abc4-567`); other IDs return 400
- Both endpoints return 404 if the movie doesn't exist in the metadata service
- A movie with no ratings yet is returned without the `rating` field.
- `PUT` returns 400 if the rating service rejects the value (it must be a whole number from 1 to 10)
- Returns 500 if no healthy instance of a downstream service is registered, or if a downstream call fails for any other reason

### Downstream services
- The service Calls `metadata` and `rating` services over **gRPC**, finding their instances through Consul on each request and picking a healthy instance at random; each call has a 5 second timeout
- It keeps one gRPC connection per downstream instance address, creating it on first use and reusing it for later requests; all connections are closed on shutdown

### gRPC API (not served yet)

[proto/movieapi/movieapi.proto](../proto/movieapi/movieapi.proto) defines a `movieapi.MovieService` with a `GetMovie` RPC that returns a `Movie` message (metadata plus an optional `rating`). It's implemented in `internal/handler/grpc`, but `main.go` doesn't start a gRPC server yet, so the service exposes only the HTTP endpoints above. The handler:

- Returns `InvalidArgument` if `movie_id` is empty or not a valid ID
- Returns `NotFound` if the movie doesn't exist in the metadata service
- Returns the movie with `rating` unset if it has no ratings yet
- Returns `Internal` if a downstream call fails

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
│   ├── handler/                    # API layer for request handling
│   │   ├── grpc/                   # gRPC server for MovieService (not wired in main.go yet)
│   │   └── http/                   # HTTP handlers (used by main.go)
│   ├── controller/                 # Business logic layer combining metadata and ratings
│   │   └── mock_controller/        # Generated gateway mocks for controller tests
│   └── gateway/                    # Gateway layer for calling downstream services
│       ├── gateway.go              # Shared gateway errors and client timeout
│       ├── grpcutil/               # gRPC connection pool and instance selection
│       ├── httputil/               # HTTP helpers used by the HTTP gateways
│       ├── metadata/
│       │   ├── grpc/               # gRPC gateway to the metadata service (used by main.go)
│       │   └── http/               # HTTP gateway to the metadata service
│       └── rating/
│           ├── grpc/               # gRPC gateway to the rating service (used by main.go)
│           └── http/               # HTTP gateway to the rating service
└── pkg/
    └── model/                      # Shared data model for the service
```
