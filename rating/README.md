# Rating Service

This Go module provides **an HTTP service** for storing and retrieving ratings for various record types (e.g. movies). It stores each user's rating for a record and computes the record's average rating.

Its design uses a **layered architecture**, separating the HTTP handling (i.e. API layer), business logic (i.e. controller layer), and data storage (i.e. repository layer) concerns.

## Endpoints

| Method | Path           | Description                           |
| ------ | -------------- | ------------------------------------- |
| `GET`  | `/{type}/{id}` | Get the average rating for a record   |
| `PUT`  | `/{type}/{id}` | Store a user's rating for a record    |

`GET` returns the average rating as a JSON number (e.g. `8.5`). `PUT` takes the user's rating as JSON:

```json
{"userID": "user-42", "value": 8}
```

## Behavior

- `PUT` uses the record type and ID from the path, ignoring any in the body
- Ratings are stored per record type, record ID, and user ID; a user's new rating for a record replaces their previous one
- `PUT` returns 400 if `userID` is empty or the value isn't a whole number from 1 to 10
- `GET` returns 404 if the record has no ratings
- Ratings are held in memory and are lost when the service stops

## Startup & graceful shutdown

- On startup, the service registers itself with Consul (at `localhost:8500`) as `rating` with a unique instance ID, then sends a heartbeat every second to keep the instance healthy
- On `SIGINT` (Ctrl+C) or `SIGTERM`, it first deregisters the instance, so callers stop discovering it, and then shuts down gracefully, giving in-flight requests up to 5 seconds to finish before closing the remaining connections
- A second Ctrl+C during shutdown forces the service to quit
- If the server fails to start (e.g. the port is already in use), the instance is deregistered before the service exits

## Running

Start the Consul server first (see the [root README](../README.md#service-discovery)), then:

```sh
go run ./rating                  # listens on localhost:8082
go run ./rating -port 9082       # listens on localhost:9082
go test -race ./rating/...
```

Use the `-port` flag to run more than one instance; each registers separately with Consul. Each instance has its own in-memory store, so instances don't share data.

## Project structure

```text
rating/
├── main.go                         # Wires the service layers and starts the HTTP server
├── internal/
│   ├── handler/                    # API layer for HTTP request handling
│   ├── controller/                 # Business logic layer for rating operations
│   └── repository/                 # Data access layer for rating storage/retrieval
│       └── memory/                 # In-memory rating storage, with race-condition tests
└── pkg/
    └── model/                      # Shared data model for the service
```
