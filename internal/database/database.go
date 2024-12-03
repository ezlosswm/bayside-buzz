package database

import (
	"database/sql"
	"log"
	"os"

	_ "embed"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed "migrations/20241201192801_add_user_table.sql"
var ddl string

var (
	dburl = os.Getenv("BLUEPRINT_DB_URL")
)

type service struct {
	db *sql.DB
}

func NewSQCL() *Queries {
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	queries := New(db)

	return queries
}
