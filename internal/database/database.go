package database

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"os"

	_ "embed"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	// _ "github.com/mattn/go-sqlite3"
)

//go:embed "migrations/*"
var ddl embed.FS

var (
	dburl = os.Getenv("BLUEPRINT_DB_URL")
)

type service struct {
	db *pgx.Conn
}

func NewSQCL() *Queries {
	db, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	sqlDB, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Fatal(err)
	}

	goose.SetBaseFS(ddl)
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	queries := New(db)

	return queries
}
