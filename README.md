# Movie Backend

This repository is a Go-based movie backend composed of separate service modules that work together in a shared workspace. It includes the metadata and rating services, each as its own Go module, and the movieapi service, which combines them into a single public API.

## Project structure

The repository is organized as a **multi-module Go workspace** so each service can evolve independently while sharing a common project setup.

### Repo overview

- `go.work` — includes the repository root and the service modules in the workspace.
- `movieapi/` — The public HTTP API for the backend; see [movieapi/README.md](movieapi/README.md).
- `metadata/` — A movie metadata service for storing and retrieving movie information; see [metadata/README.md](metadata/README.md).
- `rating/` — A rating service for storing per-user ratings and computing the average rating for a record; see [rating/README.md](rating/README.md).

## Services

| Service    | Default address  | Role                                                              |
| ---------- | ---------------- | ----------------------------------------------------------------- |
| `movieapi` | `localhost:8083` | A public API server; calls `metadata` and `rating` over HTTP      |
| `metadata` | `localhost:8081` | Stores movie metadata                                             |
| `rating`   | `localhost:8082` | Stores ratings and computes averages                              |

Each service accepts a `-port` flag to override its default port (e.g. `go run ./rating -port 9082`), so you can run several instances of the same service.

```text
client ──▶ movieapi ──┬──▶ metadata
                      └──▶ rating
```

## Service discovery

Services find each other through a [Consul](https://developer.hashicorp.com/consul) registry at `localhost:8500`, using the shared `pkg/discovery` package:

- On startup, each service registers an instance with its address and sends a heartbeat every second; Consul marks an instance unhealthy if it misses heartbeats for 5 seconds.
- On shutdown, each service deregisters its instance before draining in-flight requests.
- `movieapi` looks up the healthy `metadata` and `rating` instances for each request and picks one at random.

Start a local Consul dev server with Docker before running the services:

```sh
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client='0.0.0.0'
```

The Consul UI is then available at <http://localhost:8500>, where you can see the registered services and their health.

## Development

Run go commands for the relevant module directory from the workspace root, for example:

- `go test ./metadata/...` & `go run ./metadata`
- `go test ./rating/...` & `go run ./rating`
- `go test ./movieapi/...` & `go run ./movieapi`

Start the Consul server first, since every service registers with it on startup and exits if it can't. The movieapi service depends on the other two services, so start `metadata` and `rating` before `movieapi`.

Each service keeps its own README and internal package layout for its controller, handler, model, and repository or gateway concerns.
