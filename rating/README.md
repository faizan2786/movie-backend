# Rating Service



## Project structure

```text
rating/
├── main.go                         # Wires various service layers and starts the HTTP server
├── internal/
│   ├── handler/                    # API layer for HTTP request handling
│   ├── controller/                 # Business logic layer for rating operations
│   └── repository/                 # Data access layer for ratings storage/retrieval
│       └── memory/                 # An in-memory ratings storage
└── pkg/
    └── model/                      # Shared data model for the service
```