# Metadata Service

This Go module provides **an HTTP service** for storing and retrieving movie metadata, using an in-memory repository for storage.

Its design uses a **layered architecture**, separating the HTTP handling (i.e. API layer), business logic (i.e. controller layer), and data storage (i.e. repository layer) concerns.

## Endpoints

| Method | Path    | Description                 |
| ------ | ------- | --------------------------- |
| `GET`  | `/{id}` | Get a movie's metadata      |
| `PUT`  | `/{id}` | Store a movie's metadata    |

Metadata is sent and returned as JSON:

```json
{"id": "the-matrix", "title": "The Matrix", "releaseYear": 1999, "description": "..."}
```

## Behavior

- `PUT` uses the ID from the path, ignoring any `id` in the body
- `PUT` returns 409 if metadata for the ID already exists; stored metadata can't be updated
- `GET` returns 404 if the movie doesn't exist
- Metadata is held in memory and is lost when the service stops

## Running

```sh
go run ./metadata        # listens on localhost:8081
```

## Project structure

```text
metadata/
├── main.go                         # Wires the service layers and starts the HTTP server
├── internal/
│   ├── handler/                    # API layer for HTTP request handling
│   ├── controller/                 # Business logic layer for metadata operations
│   └── repository/                 # Data access layer for metadata storage/retrieval
│       └── memory/                 # In-memory metadata storage
└── pkg/
    └── model/                      # Shared data model for the service
```
