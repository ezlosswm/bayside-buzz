package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"bayside-buzz/cmd/web/pages"
	"bayside-buzz/internal/domain"

	"github.com/gorilla/mux"
)

const SITE_NAME = "Bayside Buzz"

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(s.corsMiddleware)

	// fileServer := http.FileServer(http.FS(web.Files))
	fileServer := http.FileServer(http.Dir("cmd/web"))

	r.PathPrefix("/assets/").Handler(fileServer)

	r.HandleFunc("/", s.HomePage).Methods(http.MethodGet)
	r.HandleFunc("/event", s.EventPage).Methods(http.MethodGet) // dynamic route
	r.HandleFunc("/contact", s.ContactPage).Methods(http.MethodGet)

	r.HandleFunc("/login", s.LoginPage).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/register", s.RegisterPage).Methods(http.MethodGet, http.MethodPost)

	r.HandleFunc("/dashboard", s.DashboardPage)
	r.HandleFunc("/dashboard/create_event", s.CreateEventPage)
	r.HandleFunc("/dashboard/create_organizer", s.CreateOrganizerPage).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/dashboard/create_organizer/{id:[0-9]+}", s.HandleDeleteOrganizer).Methods(http.MethodDelete)

	return r
}

// CORS middleware
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Wildcard allows all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Credentials not allowed with wildcard origins

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) HomePage(w http.ResponseWriter, r *http.Request) {
	var (
		url = r.Host
	)
	const (
		title       = "Bayside Buzz"
		description = "Discover Events in Corozal, Belize - Bayside Buzz"
		pageType    = "website"
		image       = "" // get an image
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	ctx := context.Background()

	if r.Method == "GET" {
		organizers, _ := s.db.GetOrganizers(ctx)

		pages.Home(pageData, &organizers).Render(context.Background(), w)
	}
}

func (s *Server) EventPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			url   = r.Host
			title = strings.Join([]string{"Get event name", SITE_NAME}, " - ")
		)
		const (
			description = "Discover Events in Corozal, Belize - Bayside Buzz"
			pageType    = "article"
			image       = "" // event's image
		)

		pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)
		slog.Info("Page Data: \n", pageData)

		pages.Event(pageData).Render(context.Background(), w)
	}
}

func (s *Server) ContactPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			url   = r.Host
			title = strings.Join([]string{"Contact Us", SITE_NAME}, " - ")
		)
		const (
			description = "Have questions, need assistance or want to advertise your business? Reach out to the Bayside Breeze team through our contact form or find our contact details here."
			pageType    = "website"
			image       = "" // get image
		)

		pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)
		slog.Info("Page Data: \n", pageData)

		pages.Contact(pageData).Render(context.Background(), w)
	}
}
