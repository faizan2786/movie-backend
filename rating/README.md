# Rating Service

This Go module provides **a gRPC service** for storing and retrieving ratings for various record types (e.g. movies). It stores each user's rating for a record and computes the record's average rating.

Its design uses a **layered architecture**, separating the gRPC request handling (i.e. API / Service layer), business logic (i.e. controller layer), and data storage (i.e. repository layer) concerns.

## RPCs

The `rating.RatingService` service is defined in [proto/rating/rating.proto](../proto/rating/rating.proto):

| RPC                   | Request fields                                          | Description                         |
| --------------------- | ------------------------------------------------------- | ----------------------------------- |
| `GetAggregatedRating` | `record_type`, `record_id`                              | Get the average rating for a record |
| `PutRating`           | `record_type`, `record_id`, `user_id`, `rating_value`   | Store a user's rating for a record  |

`record_type` is a `RecordType` enum; `RECORD_TYPE_MOVIE` (`1`) is currently the only supported type. `GetAggregatedRating` returns the average rating in the optional `rating` field. `PutRating` returns an empty response.

## Service Behaviour

- Ratings are held in memory and are lost when the service stops
- Ratings are stored per (record type, record ID, and user ID); a user's new rating for a record replaces their previous one
- Both RPCs return `InvalidArgument` if `record_type` is unspecified or unknown, or if `record_id` is empty
- `PutRating` returns `InvalidArgument` if `user_id` is empty or `rating_value` isn't from 1 to 10
- `GetAggregatedRating` returns `OK` with no `rating` set if the record has no ratings yet

## Running

Start the Consul server first (see the [root README](../README.md#service-discovery)), then:

```sh
go run ./rating                  # listens on localhost:8082
go run ./rating -port 9082       # listens on localhost:9082
go test -race ./rating/...
```

Use the `-port` flag to run more than one instance; each registers separately with Consul. Each instance has its own in-memory store, so instances don't share data.

### Example gRPC requests

The server registers gRPC reflection, so you can call it with [grpcurl](https://github.com/fullstorydev/grpcurl) without passing any `.proto` files. `record_type` `1` is `RECORD_TYPE_MOVIE`.

Store a user's rating for a record (returns an empty response on success):

```sh
grpcurl -plaintext -d @ localhost:8082 rating.RatingService.PutRating << EOM
{
  "record_type": 1,
  "record_id": "123",
  "user_id": "1",
  "rating_value": 8
}
EOM
```
Output:
```text
{}
```

Get the average rating for a record:

```sh
grpcurl -plaintext -format text -d 'record_id:"123", record_type:1' localhost:8082 rating.RatingService.GetAggregatedRating
```
Output:
```text
rating: 8
```

## Project structure

```text
rating/
├── main.go                         # Wires the service layers and starts the gRPC server
├── internal/
│   ├── handler/                    # API layer
│   │   ├── grpc/                   # gRPC server (used by main.go)
│   │   └── http/                   # Legacy HTTP handlers (not wired up)
│   ├── controller/                 # Business logic layer for rating operations
│   └── repository/                 # Data access layer for rating storage/retrieval
│       └── memory/                 # In-memory rating storage, with race-condition tests
└── pkg/
    └── model/                      # Shared data model for the service
```
