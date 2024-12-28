package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"

	"bayside-buzz/internal/database"
)

var (
	cookie = os.Getenv("COOKIE")
)

type Server struct {
	port int

	db    *database.Queries
	store *sessions.CookieStore
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	session := sessions.NewCookieStore([]byte(cookie))
	session.Options.HttpOnly = true
	session.Options.SameSite = http.SameSiteLaxMode

	NewServer := &Server{
		port: port,

		db:    database.NewSQCL(),
		store: session,
	}

    fmt.Println("heyyy")

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
