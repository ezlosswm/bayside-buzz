package server

import (
	"bayside-buzz/cmd/web/pages/dashboard"
	"bayside-buzz/internal/database"
	"bayside-buzz/internal/domain"
	"bayside-buzz/internal/lib"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Dashboard
func (s *Server) DashboardPage(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(contextKeyUser).(*database.User)

	var (
		url   = r.Host
		title = strings.Join([]string{"Dashboard", SITE_NAME}, " - ")
	)
	const (
		description = "Manage events, view organizer details, and monitor your community activity with the Bayside Breeze dashboard."
		pageType    = "website"
		image       = "" // get image
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	ctx := context.Background()
	if r.Method == "GET" {
		// returns user struct
		if ok {
			totalOrganizers, _ := s.db.CountOrganizers(ctx)

			results := domain.NewResults(totalOrganizers)

			dashboard.Dashboard(pageData, results).Render(context.Background(), w)
		} else {
			w.Header().Set("HX-Redirect", "/login")
			return
		}
	}
}

func (s *Server) CreateEventPage(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(contextKeyUser).(*database.User)

	var (
		url   = r.Host
		title = strings.Join([]string{"Create New Event", SITE_NAME}, " - ")
	)
	const (
		description = "Host your next big event with Bayside Breeze. Fill in the details and share your event with the community today!"
		pageType    = "website"
		image       = "" // get image
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	if r.Method == "GET" {
		if ok {
			ctx := context.Background()
			organizers, _ := s.db.GetOrganizers(ctx)
			dashboard.CreateEvent(pageData, &organizers).Render(context.Background(), w)
		} else {
			w.Header().Set("HX-Redirect", "/login")
			return
		}
	}

	if r.Method == "POST" {
		if !ok {
			w.Header().Set("HX-Redirect", "/login")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			slog.Error("error parsing registration form\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			w.Header().Set("HX-Refresh", "true")
			return
		}
		defer r.MultipartForm.RemoveAll()

		eventTitle := r.FormValue("event__title")
		eventDesc := r.FormValue("event__description")
		eventDate := r.FormValue("event__date")
		eventFrequency := r.FormValue("event__frequency")
		eventOrganizer := r.FormValue("event__organizer")

		file, fileHeader, err := r.FormFile("cover__img")
		if err != nil {
			slog.Error("err here", err)
			http.Error(w, "Unable to retrieve file.", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		if err != nil {
			http.Error(w, "Unable to read file.", http.StatusInternalServerError)
			return
		}

		imgPath, err := lib.FileUpload("events", file, *fileHeader)
		if err != nil {
			http.Error(w, "Unable to create file.", http.StatusInternalServerError)
			return
		}

		type f struct {
			t, d, e, f, o, p string
		}

		ff := f{
			t: eventTitle,
			d: eventDesc,
			e: eventDate,
			f: eventFrequency,
			o: eventOrganizer,
			p: imgPath,
		}

		slog.Info("event information\n\n", ff)
	}

}

// Organizer name is a unique value, make sure to handle error
func (s *Server) CreateOrganizerPage(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(contextKeyUser).(*database.User)

	var (
		url   = r.Host
		title = strings.Join([]string{"Create New Organizer", SITE_NAME}, " - ")
	)
	const (
		description = "Expand your network by adding a new organizer. Manage events efficiently and connect with the community through Bayside Breeze."
		pageType    = "website"
		image       = "" // get image
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	ctx := context.Background()
	if r.Method == "GET" {
		if ok {
			organizers, _ := s.db.GetOrganizers(ctx)

			dashboard.CreateOrganizer(pageData, &organizers).Render(context.Background(), w)
		} else {
			w.Header().Set("HX-Redirect", "/login")
			return
		}
	}

	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			slog.Error("error parsing registration form\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			w.Header().Set("HX-Refresh", "true")
			return
		}
		defer r.MultipartForm.RemoveAll()

		orgName := r.FormValue("org__name")
		orgDesc := r.FormValue("org__description")
		file, fileHeader, err := r.FormFile("org__img")
		if err != nil {
			slog.Error("err here", err)
			http.Error(w, "Unable to retrieve file.", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		if err != nil {
			http.Error(w, "Unable to read file.", http.StatusInternalServerError)
			return
		}

		imgPath, err := lib.FileUpload("organizers", file, *fileHeader)
		if err != nil {
			http.Error(w, "Unable to create file.", http.StatusInternalServerError)
			return
		}

		value := lib.OrganizerToValue(orgName)

		newOrganizer := database.CreateOrganizerParams{
			OrganizerName: orgName,
			Description:   orgDesc,
			Value:         value,
			ImgUrl:        imgPath,
		}

		if err := s.db.CreateOrganizer(ctx, newOrganizer); err != nil {
			http.Error(w, "Error saving data to database.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Refresh", "true")
	}
}

func (s *Server) HandleDeleteOrganizer(w http.ResponseWriter, r *http.Request) {
	// returns user struct
	_, ok := r.Context().Value(contextKeyUser).(*database.User)

	ctx := context.Background()
	if r.Method == "DELETE" {
		if !ok {
			w.Header().Set("HX-Redirect", "/login")
			return
		}

		vars := mux.Vars(r)
		idParam := vars["id"]

		id, err := strconv.Atoi(idParam)
		if err != nil {
			slog.Error("error converting id")
			return
		}

		if err := s.db.DeleteOrganizer(ctx, int64(id)); err != nil {
			http.Error(w, "Error delete organizer.", http.StatusInternalServerError)
			return
		}
	}
}
