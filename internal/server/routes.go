package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"bayside-buzz/cmd/web/components"
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
	r.HandleFunc("/organizers/{organizer}", s.FilterByOrganizer).Methods(http.MethodGet)
	r.HandleFunc("/event/{id:[0-9]+}", s.EventPage)
	r.HandleFunc("/contact", s.ContactPage).Methods(http.MethodGet)

	r.HandleFunc("/login", s.LoginPage).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/logout", s.HandleLogout).Methods(http.MethodPost)

	dashboard := r.PathPrefix("/dashboard").Subrouter()
	dashboard.HandleFunc("", s.Authenticate(s.HandleDashboard))
	dashboard.HandleFunc("/events", s.Authenticate(s.HandleEvents))
	dashboard.HandleFunc("/events/{id:[0-9]+}/delete", s.Authenticate(s.HandleEventsDelete))
	dashboard.HandleFunc("/events/{id:[0-9]+}/edit", s.Authenticate(s.HandleEventsEdit))

	dashboard.HandleFunc("/organizers", s.Authenticate(s.HandleOrganizers))
	dashboard.HandleFunc("/organizers/{id:[0-9]+}/delete", s.Authenticate(s.HandleOrganizersDelete))
	dashboard.HandleFunc("/organizers/{id:[0-9]+}/edit", s.Authenticate(s.HandleOrganizersEdit))

	return r
}

func (s *Server) checkSession(r *http.Request) bool {
	session, err := s.store.Get(r, "login")
	if err != nil {
		slog.Error("error getting user", "error", err)
		return false
	}

	_, ok := session.Values["userId"]

	return ok
}

func (s *Server) HomePage(w http.ResponseWriter, r *http.Request) {
	settings := domain.NewSettings()
	settings.PageData.URL = r.Host

	if r.Method == "GET" {
		organizers, err := s.db.GetOrganizers(context.Background())
		if err != nil {
			slog.Error("Error getting organizers", "error", err)
		}

		events, err := s.db.GetEvents(context.Background())
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

		ok := s.checkSession(r)
		if !ok {
			pages.Home(settings, events, organizers).Render(context.Background(), w)
			return
		}

		settings.IsLoggedIn = true
		pages.Home(settings, events, organizers).Render(context.Background(), w)
	}
}

func (s *Server) FilterByOrganizer(w http.ResponseWriter, r *http.Request) {
	orgParam := mux.Vars(r)["organizer"]
	if orgParam == "all" {
		allEvents, err := s.db.GetEvents(context.Background())
		if err != nil {
			log.Println("Error getting all events", "error", err)
		}
		components.AllEvents(allEvents).Render(context.Background(), w)
	}

	eventsByOrg, err := s.db.GetEventsByOrganizer(context.Background(), orgParam)
	if err != nil {
		log.Println("Error getting events by organizer", "error", err)
	}

	components.AllEvents(eventsByOrg).Render(context.Background(), w)
}

func (s *Server) EventPage(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	events, _ := s.db.GetEventWithTags(context.Background(), int32(id))

	settings := domain.NewSettings()
	settings.PageData.URL = r.Host
	settings.PageData.Title = strings.Join([]string{settings.PageData.SiteName, events.Title}, " - ")

	ok := s.checkSession(r)
	if !ok {
		pages.Event(settings, events).Render(context.Background(), w)
		return
	}

	settings.IsLoggedIn = true
	pages.Event(settings, events).Render(context.Background(), w)
}

func (s *Server) ContactPage(w http.ResponseWriter, r *http.Request) {
	settings := domain.NewSettings()
	settings.PageData.URL = r.Host

	ok := s.checkSession(r)
	if !ok {
		pages.Contact(settings).Render(context.Background(), w)
		return
	}

	settings.IsLoggedIn = true
	pages.Contact(settings).Render(context.Background(), w)
}
