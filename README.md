# Movie Backend

This repository is a Go-based movie backend composed of separate service modules that work together in a shared workspace. The current implementation includes the metadata and rating services, each as its own Go module.

## Structure

- `go.work` — includes the repository root and the service modules in the workspace.
- `metadata/` — A movie metadata service for storing and retrieving movie information; see [metadata/README.md](metadata/README.md).
- `rating/` — A rating service for storing per-user ratings and computing the average rating for a record; see [rating/README.md](rating/README.md).

## Modules

The repository is organized as a multi-module Go workspace so each service can evolve independently while sharing a common project setup.

## Development

Run go commands for the relevant module directory from the workspace root, for example:

- `go test ./metadata/...` & `go run .metadata/...`
- `go test ./rating/...` & `go run .rating/...`

Each service keeps its own README and internal package layout for controller, handler, repository, and model concerns.