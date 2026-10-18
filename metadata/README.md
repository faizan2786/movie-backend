# Metadata Service

This Go module provides **an HTTP service** for storing and retrieving movie metadata. It exposes `GET /{id}` and `PUT /{id}` endpoints and also provides an in-memory repository for storing metadata.

Its design uses a **layered architecture**, separating the HTTP handling (i.e. API layer), business logic (i.e. controller layer), and data storage (i.e. repository layer) concerns.

## Project structure

```text
metadata/
├── main.go                         # Wires various service layers and starts the HTTP server
├── internal/
│   ├── handler/                    # API layer for HTTP request handling
│   ├── controller/                 # Business logic layer for Metadata operations
│   └── repository/                 # Data access layer for Metadata storage/retrieval
│       └── memory/                 # An in-memory metadata storage
└── pkg/
    └── model/                      # Shared data model for the service
```