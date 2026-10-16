# Learning Stuff

A bookmarks app built while learning: a Go API (`api/`) and a Next.js frontend (`web/`).

## Local development

### Prerequisites

- Go
- Docker (for local Postgres)

### Database

Start Postgres in Docker:

    docker compose up -d db

Or with Apple's `container` CLI:

    container run --name bookmarks-db -d \
      -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=bookmarks \
      -p 5432:5432 postgres:17

Export the connection string:

    export DATABASE_URL="postgres://postgres:postgres@localhost:5432/bookmarks?sslmode=disable"

To stop it and delete all data:

    docker compose down -v

### API

    cd api
    go run .

The API reads `DATABASE_URL` from the environment and exits with a clear error
if it is not set. SQL migrations in `api/migrations/` run automatically on
startup.

Run the tests (database-backed tests are skipped unless `DATABASE_URL` is set):

    go test ./...
