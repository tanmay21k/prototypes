# SQLC Database Template

This repository is a small Go web-service template built around:

- Go standard package layout
- Chi for HTTP routing
- PostgreSQL as the database
- sqlc for type-safe SQL code generation
- a simple layered architecture for easy reuse in future services

The goal is to keep the project easy to read, easy to extend, and aligned with common Go conventions.

---

## Project purpose

This project demonstrates a minimal but realistic backend structure:

- HTTP layer handles requests and responses
- feature packages handle business logic
- generated data access code handles SQL access
- the database adapter is isolated from the rest of the application

This makes it a good reference project for building a Go API that uses sqlc without creating a large, hard-to-follow application.

---

## Folder structure and why it is designed this way

```text
.
├── cmd/
│   ├── api.go
│   └── main.go
├── internal/
│   ├── adapters/
│   │   └── postgresql/
│   │       └── sqlc/
│   │           ├── db.go
│   │           ├── models.go
│   │           ├── querier.go
│   │           ├── queries.sql.go
│   │           └── queries.sql
│   ├── json/
│   │   └── json.go
│   └── products/
│       ├── handler.go
│       └── service.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── sqlc.yaml
└── README.md
```

### cmd/

This is the application entrypoint.

- `cmd/main.go` contains the startup logic and initializes the application configuration.
- `cmd/api.go` contains the HTTP server setup, route registration, and the wiring between dependencies.

Why this folder exists:

- application startup should be isolated from business logic
- the top-level app bootstrap should remain small and easy to read
- it keeps the project consistent with Go conventions where the entrypoint is separate from the domain code

This is where you wire together:

- HTTP routing
- configuration
- database access
- feature packages

### internal/

The `internal` directory holds code that is private to this module and should not be imported from outside the project.

This is a common Go pattern for application-specific packages that should not be publicly exposed.

#### internal/products/

This package contains the product feature logic.

- `handler.go` contains HTTP handlers
- `service.go` contains business logic and application-facing service methods

Why this is separated:

- handlers deal with HTTP concerns like request parsing and response writing
- service logic deals with business rules and orchestration
- this separation keeps the feature easy to test and extend

#### internal/json/

This is a tiny helper package used to write JSON responses consistently.

Why this exists:

- avoids repeated response boilerplate in handlers
- keeps handlers focused on request flow instead of encoding logic

#### internal/adapters/postgresql/sqlc/

This is the database adapter layer generated and managed by sqlc.

Generated files such as:

- `db.go`
- `models.go`
- `querier.go`
- `queries.sql.go`

contain:

- database connection wrappers
- generated SQL query functions
- typed models matching the schema
- interfaces for repository access

Why this is structured this way:

- the application should not depend directly on raw SQL everywhere
- sqlc gives type-safe access to PostgreSQL queries
- the database layer is isolated behind generated interfaces, so the rest of the code stays clean

The important idea is that application code talks to `repo.Querier`, not to a concrete database implementation scattered across the project.

---

## Why this layering matters

This project intentionally separates responsibilities into layers:

1. `cmd` = composition root
2. `internal/products` = feature logic
3. `internal/adapters/postgresql/sqlc` = data access layer
4. `internal/json` = shared response helper

This reduces coupling and makes it easier to:

- swap database implementations later
- add new features without touching unrelated code
- test business logic independently from HTTP concerns
- generate typed SQL access from schema and SQL files instead of hand-writing repository code

---

## SQLC and generated code

The project uses sqlc via `sqlc.yaml`.

```yaml
version: "2"
sql:
- schema: "./internal/adapters/postgresql/migrations"
  queries: "./internal/adapters/postgresql/sqlc/queries.sql"
  engine: "postgresql"
  gen:
    go:
      package: "repo"
      out: "./internal/adapters/postgresql/sqlc"
      sql_package: "pgx/v5"
      emit_json_tags: true
      emit_interface: true
```

This tells sqlc to:

- read the schema and SQL query files
- generate Go types and query methods
- write them into the repository package under `internal/adapters/postgresql/sqlc`

This is why the generated DB code lives under a dedicated adapter package instead of being mixed into the HTTP or service layer.

---

## Running the full stack locally

This project expects PostgreSQL to be available before the API boots.

### 1) Start the database with Docker Compose

The Postgres container is configured with these credentials in [docker-compose.yaml](docker-compose.yaml):

- username: `myuser`
- password: `mypassword`
- database: `ecom`
- port: `5432`

Start it from the project root:

```bash
docker compose up -d
```

To verify the database is up:

```bash
docker compose ps
```

You can also connect locally with:

```bash
psql "postgres://myuser:mypassword@localhost:5432/ecom?sslmode=disable"
```

### 2) Run the Go API

From the project root:

```bash
go run .
```

This is the safest way to start the service in a multi-file Go package. Running a single file like `go run cmd/main.go` is not the correct pattern when other files in the same package provide the types and methods that `main.go` depends on.

The app listens on:

```text
http://localhost:8080
```

### 3) Smoke test

```bash
curl http://localhost:8080/
```

and:

```bash
curl http://localhost:8080/products
```

---

## Example API shape

The current service is a simple product API template.

```text
GET /products
GET /product/{productId}
```

The request flow is roughly:

- router in `cmd/api.go`
- HTTP handler in `internal/products/handler.go`
- service method in `internal/products/service.go`
- generated SQLC repository in `internal/adapters/postgresql/sqlc`
- PostgreSQL database

---

## Migration note

This template is designed to work with PostgreSQL, and the `sqlc` configuration expects schema and migration files under the PostgreSQL adapter path.

A typical migration command looks like:

```bash
goose -dir ./internal/adapters/postgresql/migrations postgres "postgres://myuser:mypassword@localhost:5432/ecom?sslmode=disable" up
```

This keeps database migration logic close to the database adapter rather than spreading it across the app.

---

## Summary

This repository is structured to follow a clean, reusable backend architecture:

- `cmd` for startup and composition
- `internal/products` for feature logic
- `internal/adapters/postgresql/sqlc` for database access
- `internal/json` for shared response helpers
- `sqlc.yaml` for generated database code

That structure becomes a useful template because it keeps the app readable, testable, and easy to extend as a real project grows.
