# Project Phoenix Backend

The backend is a standard-library-only Go HTTP API with a layered structure:

- `cmd/api`: application bootstrap
- `internal/config`: environment-backed configuration
- `internal/handlers`: HTTP handlers and JSON responses
- `internal/server`: HTTP server and routing
- `internal/version`: centralized version metadata

## Run

```sh
go run ./cmd/api
```

The server listens on `:8080` by default. Set `PHOENIX_HTTP_ADDR` to override the address.

## Endpoints

- `GET /health` returns `{"status":"ok"}`
- `GET /version` returns the current application version
- `POST /projects` creates a project from a JSON body containing `name` and an optional `description`
- `GET /projects` returns all projects
- `GET /projects/{id}` returns a project by ID
