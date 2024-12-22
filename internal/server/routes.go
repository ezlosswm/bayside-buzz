package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"bayside-buzz/cmd/web/pages"
	"bayside-buzz/internal/domain"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
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
	r.HandleFunc("/event/{id:[0-9]+}", s.EventPage)
	r.HandleFunc("/contact", s.ContactPage).Methods(http.MethodGet)

	r.HandleFunc("/login", s.LoginPage).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/logout", s.HandleLogout).Methods(http.MethodPost)
	r.HandleFunc("/register", s.RegisterPage).Methods(http.MethodGet, http.MethodPost)

	r.HandleFunc("/dashboard", s.Authenticate(s.DashboardPage))
	r.HandleFunc("/dashboard/create_event", s.Authenticate(s.CreateEventPage))
	r.HandleFunc("/dashboard/create_event/{id:[0-9]+}", s.Authenticate(s.HandleDeleteEvent)).Methods(http.MethodDelete)
	r.HandleFunc("/dashboard/create_organizer", s.Authenticate(s.CreateOrganizerPage)).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/dashboard/create_organizer/{id:[0-9]+}", s.Authenticate(s.HandleDeleteOrganizer)).Methods(http.MethodDelete)

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
	var url = r.Host

	const (
		title       = "Bayside Buzz"
		description = "Discover Events in Corozal, Belize - Bayside Buzz"
		pageType    = "website"
		image       = "/assets/images/corozal-sign.jpg"
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	if r.Method == "GET" {
		organizers, err := s.db.GetOrganizers(context.Background())
		if err != nil {
			slog.Error("Error getting organizers", "error", err)
		}

		events, err := s.db.GetEventsWithTags(context.Background())
		if err != nil {
			slog.Error("Error getting events", "error", err)

			if pgErr, ok := err.(*pgconn.PgError); ok {
				slog.Error("PostgreSQL error details",
					"code", pgErr.Code,
					"message", pgErr.Message,
					"detail", pgErr.Detail,
					"hint", pgErr.Hint)
			}
		}
		pages.Home(pageData, events, organizers).Render(context.Background(), w)
	}
}

func (s *Server) EventPage(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	events, _ := s.db.GetEventWithTags(context.Background(), int32(id))

	var (
		url   = r.Host
		image = fmt.Sprintf("%v", events.Imgpath)
	)

	const (
		description = "Discover Events in Corozal, Belize - Bayside Buzz"
		pageType    = "article"
	)

	// update OG info to match the event
	title := strings.Join([]string{SITE_NAME, events.Title}, " - ")
	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	pages.Event(pageData, events).Render(context.Background(), w)
}

func (s *Server) ContactPage(w http.ResponseWriter, r *http.Request) {
	var (
		url   = r.Host
		title = strings.Join([]string{SITE_NAME, "Contact Us"}, " - ")
	)
	const (
		description = "Have questions, need assistance or want to advertise your business? Reach out to the Bayside Breeze team through our contact form or find our contact details here."
		pageType    = "website"
		image       = "/assets/images/corozal-sign.jpg"
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	pages.Contact(pageData).Render(context.Background(), w)
}
