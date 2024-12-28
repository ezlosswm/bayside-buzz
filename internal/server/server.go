package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/storage"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"

	"bayside-buzz/internal/database"
)

var (
	cookie         = os.Getenv("COOKIE")
	appwriteClient client.Client
)

type Server struct {
	port int

	db      *database.Queries
	store   *sessions.CookieStore
	storage *storage.Storage
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	appwriteClient = appwrite.NewClient(
		appwrite.WithProject(os.Getenv("PROJECT_KEY")),
		appwrite.WithKey(os.Getenv("API_KEY")),
	)

	session := sessions.NewCookieStore([]byte(cookie))
	session.Options.HttpOnly = true
	session.Options.SameSite = http.SameSiteLaxMode

	NewServer := &Server{
		port: port,

		db:      database.NewSQCL(),
		store:   session,
		storage: storage.New(appwriteClient),
	}

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
