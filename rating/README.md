# Rating Service

This Go module provides **an HTTP service** backend for storing and retrieving ratings for various record types. It stores user submissions for a given record and computes the average rating for that record.

Its design uses a **layered architecture**, separating the HTTP handling (i.e. API layer), business logic (i.e. controller layer), and data storage (i.e. repository layer) concerns.

## What it does

- Stores a rating for a specific record using its `recordType`, `recordID`, and `userID`
- Validates empty IDs and out-of-range rating values
- Retrieves the current average rating for a record

## Architecture

- `main.go` — starts the service and wires together the repository, controller, and handler
- `internal/handler/` — HTTP-facing request handling layer
- `internal/controller/` — business logic and repository interaction
- `internal/repository/` — repository contract and in-memory implementation
- `pkg/model/` — shared types for record IDs, user IDs, and ratings

## Project structure

```text
rating/
├── main.go                         # Starts the service and configures the HTTP server
├── internal/
│   ├── handler/                    # HTTP handlers for the rating service
│   ├── controller/                 # Business logic for rating operations
│   └── repository/
│       ├── errors.go               # Shared repository errors
│       └── memory/                 # In-memory storage for ratings
│           ├── db.go                # Repository implementation
│           ├── db_test.go           # Repository tests
│           └── db_race_test.go      # Race-condition coverage
├── pkg/
│   └── model/
│       └── model.go                # Shared rating data model
└── README.md                       # Service documentation
```

## Repository behavior

The in-memory database keeps ratings grouped by:

- record type
- record ID
- user ID

`GetRating` returns a `float64` average for the record and returns `repository.ErrNotFound` when the record does not exist.

`UpsertRating` enforces these validation rules:
- empty `recordType`, `recordID`, or `userID` is rejected
- values below 1 or above 10 are rejected