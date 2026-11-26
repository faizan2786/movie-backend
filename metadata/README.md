# Metadata Service

This Go module provides **a gRPC service** for storing and retrieving movie metadata, using an in-memory repository for storage.

Its design uses a **layered architecture**, separating the gRPC request handling (i.e. API/Service layer), business logic (i.e. controller layer), and data storage (i.e. repository layer) concerns.

## RPCs

The `metadata.MetadataService` service is defined in [proto/metadata/metadata.proto](../proto/metadata/metadata.proto):

| RPC           | Request fields | Description              |
| ------------- | -------------- | ------------------------ |
| `GetMetadata` | `movie_id`     | Get a movie's metadata   |
| `PutMetadata` | `metadata`     | Store a movie's metadata |

Both RPCs use the `Metadata` message, which has `movie_id`, `title`, `release_year`, and `description` fields. `GetMetadata` returns it in the `metadata` field, and `PutMetadata` returns an empty response.

## Service Behaviour

### GetMetadata

- Returns `InvalidArgument` if `movie_id` is empty
- Returns `NotFound` if the movie doesn't exist

### PutMetadata

- Metadata are held in memory and are lost when the service stops
- Creates the metadata if the movie doesn't exist, otherwise replaces it
- Returns `InvalidArgument` if `movie_id`, `title` or `release_year` is missing
- Returns `InvalidArgument` if `title` or `release_year` differs from the existing metadata; only `description` can be changed

## Running

Start the Consul server first (see the [root README](../README.md#service-discovery)), then:

```sh
go run ./metadata                # listens on localhost:8081
go run ./metadata -port 9081     # listens on localhost:9081
```

Use the `-port` flag to run more than one instance; each registers separately with Consul. Each instance has its own in-memory store, so instances don't share data.

### Example gRPC requests

The server registers gRPC reflection, so you can call it with [grpcurl](https://github.com/fullstorydev/grpcurl) without passing any `.proto` files.

Store a movie's metadata (returns an empty response on success):

```sh
grpcurl -plaintext -d @ localhost:8081 metadata.MetadataService.PutMetadata << EOM
{
  "metadata": {
    "movie_id": "the-matrix",
    "title": "The Matrix",
    "release_year": 1999,
    "description": "A hacker learns the truth about his reality."
  }
}
EOM
```
Output:
```text
{}
```

Get a movie's metadata:

```sh
grpcurl -plaintext -format text -d 'movie_id:"the-matrix"' localhost:8081 metadata.MetadataService.GetMetadata
```
Output:
```text
metadata: <
  movie_id: "the-matrix"
  title: "The Matrix"
  release_year: 1999
  description: "A hacker learns the truth about his reality."
>
```

## Project structure

```text
metadata/
├── main.go                         # Wires the service layers and starts the gRPC server
├── internal/
│   ├── handler/                    # API layer
│   │   ├── grpc/                   # gRPC server (used by main.go)
│   │   └── http/                   # Legacy HTTP handlers (not wired up)
│   ├── controller/                 # Business logic layer for metadata operations
│   └── repository/                 # Data access layer for metadata storage/retrieval
│       └── memory/                 # In-memory metadata storage
└── pkg/
    └── model/                      # Shared data model for the service
```
