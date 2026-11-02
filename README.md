# Movie Backend

This repository is a Go-based movie backend composed of separate service modules that work together in a shared workspace. It includes the metadata and rating services, each as its own Go module, and the movieapi service, which combines them into a single public API.

## Structure

- `go.work` — includes the repository root and the service modules in the workspace.
- `movieapi/` — The public HTTP API; it combines a movie's metadata and average rating, and accepts user ratings for movies; see [movieapi/README.md](movieapi/README.md).
- `metadata/` — A movie metadata service for storing and retrieving movie information; see [metadata/README.md](metadata/README.md).
- `rating/` — A rating service for storing per-user ratings and computing the average rating for a record; see [rating/README.md](rating/README.md).

## Services

| Service    | Address          | Role                                                             |
| ---------- | ---------------- | ---------------------------------------------------------------- |
| `movieapi` | `localhost:8083` | A public API server; calls `metadata` and `rating` over HTTP |
| `metadata` | `localhost:8081` | Stores movie metadata                                            |
| `rating`   | `localhost:8082` | Stores ratings and computes averages                             |

```text
client ──▶ movieapi ──┬──▶ metadata
                      └──▶ rating
```

## Modules

The repository is organized as a multi-module Go workspace so each service can evolve independently while sharing a common project setup.

## Development

Run go commands for the relevant module directory from the workspace root, for example:

- `go test ./metadata/...` & `go run ./metadata`
- `go test ./rating/...` & `go run ./rating`
- `go test ./movieapi/...` & `go run ./movieapi`

The movieapi service depends on the other two services, so start `metadata` and `rating` before `movieapi`.

Each service keeps its own README and internal package layout for its controller, handler, model, and repository or gateway concerns.
