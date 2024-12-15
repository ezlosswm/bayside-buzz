# Project bayside-buzz

Bayside Buzz is an event listing website with a dashboard for new listings.

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
- SQLite3

## To Do

- Fix UI/UX for dashboard
- Add a way to manage paid ads
- Ability to share an event
- Switch to Postgres
- Dockerize
- Deploy
