# Project bayside-buzz

Bayside Buzz is an event listing website with a dashboard for new listings.

## Installation

_Requirements: make, docker, air, goose, sqlc, gorilla/mux/ godotenv_

1. Ensure all depencies are installed

```Shell
go mod tidy
```

2. Copy `env.example` to `.env` file.

```Shell
cp env.example .env
```

3. Get [Appwrite](https://appwrite.io/docs) API Key, Project ID & Bucket ID.

4. Build and run PostgresSQL server

```Shell
docker build -t bayside-buzz-postgres .

docker run -p 5432:5432 -e POSTGRES_PASSWORD=mypassword bayside-buzz-postgres
```

5. Build and run the API server

```Shell
# Dev
make watch

# Build
make build

# Run
./main
```

## Pages

Home

- Lists all posted events and available organizers for the events.

Single Event

- Dynamic route that displays all information for an event.

Contact

Login & Register (Limits to only one admin user)

Dashboard (Private)

- Shows the total number of events & organizers
- Lists all events

Create Organizers (Private)

- Allows you to create a new organizer
- Lists all organizers

Create Event

- Creates a new event

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```

Run migration up

```bash
make migrate-up
```

Rollback migrations

```bash
make migrate-rollback
```

Generates SQLc queries

```bash
make sqlc-generate
```

## Tech Stack

- Go & Gorilla Mux
- Templ
- Tailwind
- HTMX
- Postgres
- Appwrite

## To Do

- [ ] Fix UI/UX for dashboard
- [ ] Add a way to manage paid ads
- [x] Ability to share an event
- [x] Switch to Postgres
- [x] Deploy
