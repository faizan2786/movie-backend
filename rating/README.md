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

## Running

```sh
go run ./rating          # listens on localhost:8082
go test -race ./rating/...
```

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
