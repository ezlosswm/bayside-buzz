package server

import (
	"bayside-buzz/cmd/web/pages/dashboard"
	"bayside-buzz/internal/database"
	"bayside-buzz/internal/domain"
	"bayside-buzz/internal/lib"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// Dashboard
func (s *Server) DashboardPage(w http.ResponseWriter, r *http.Request) {
	var (
		url   = r.Host
		title = strings.Join([]string{SITE_NAME, "Dashboard"}, " - ")
	)

	const (
		description = "Manage events, view organizer details, and monitor your community activity with the Bayside Breeze dashboard."
		pageType    = "website"
		image       = "/assets/images/corozal-sign.jpg"
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	e, err := s.db.GetEvents(context.Background())
	if err != nil {
		slog.Error("error getting event data\n", e)
		return
	}

	totalOrganizers, _ := s.db.CountOrganizers(context.Background())
	totalEvents, _ := s.db.CountEvents(context.Background())

	results := domain.NewResults(totalOrganizers, totalEvents)

	dashboard.Dashboard(pageData, e, results).Render(context.Background(), w)
}

func (s *Server) CreateEventPage(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(contextKeyUser).(*database.User)

	var (
		url   = r.Host
		title = strings.Join([]string{SITE_NAME, "Create New Event"}, " - ")
	)
	const (
		description = "Host your next big event with Bayside Breeze. Fill in the details and share your event with the community today!"
		pageType    = "website"
		image       = "/assets/images/corozal-sign.jpg"
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	organizers, _ := s.db.GetOrganizers(context.Background())
	if r.Method == "GET" {
		dashboard.CreateEvent(pageData, false, organizers).Render(context.Background(), w)
	}

	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
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

		// Get all selected tags
		eventTags := r.PostForm["event__tags[]"]

		file, fileHeader, err := r.FormFile("cover__img")
		if err != nil {
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

		layout := "2006-01-02"
		parsedTime, _ := time.Parse(layout, eventDate)

		var pgDate pgtype.Date
		pgDate.Time = parsedTime
		pgDate.Valid = true

		e := database.CreateEventParams{
			Title:       eventTitle,
			Description: eventDesc,
			Date:        pgDate,
			Freq:        eventFrequency,
			Organizer:   eventOrganizer,
			Imgpath:     imgPath,
			Userid:      user.ID,
		}

		if err := s.saveEvent(context.Background(), e, eventTags); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dashboard.CreateEvent(pageData, true, organizers).Render(context.Background(), w)
	}
}

func (s *Server) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {	
	if r.Method == "DELETE" {
		vars := mux.Vars(r)
		idParam := vars["id"]

		id, err := strconv.Atoi(idParam)
		if err != nil {
			slog.Error("error converting id")
			return
		}

		if err := s.db.DeleteEvent(context.Background(), int32(id)); err != nil {
			http.Error(w, "error delete organizer.", http.StatusInternalServerError)
			return
		}
	}
}

// Organizer name is a unique value, make sure to handle error
func (s *Server) CreateOrganizerPage(w http.ResponseWriter, r *http.Request) {
	var (
		url   = r.Host
		title = strings.Join([]string{SITE_NAME, "Create New Organizer"}, " - ")
	)
	const (
		description = "Expand your network by adding a new organizer. Manage events efficiently and connect with the community through Bayside Breeze."
		pageType    = "website"
		image       = "/assets/images/corozal-sign.jpg"
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	organizers, _ := s.db.GetOrganizers(context.Background())
	if r.Method == "GET" {
		dashboard.CreateOrganizer(pageData, false, organizers).Render(context.Background(), w)
	}

	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			dashboard.CreateOrganizer(pageData, false, organizers).Render(context.Background(), w)

			// slog.Error("error parsing registration form\n", err.Error())
			w.Header().Set("HX-Refresh", "true")
			return
		}
		defer r.MultipartForm.RemoveAll()

		orgName := r.FormValue("org__name")
		orgDesc := r.FormValue("org__description")
		file, fileHeader, err := r.FormFile("org__img")
		if err != nil {
			http.Error(w, "unable to retrieve file.", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		imgPath, err := lib.FileUpload("organizers", file, *fileHeader)
		if err != nil {
			http.Error(w, "unable to create file.", http.StatusInternalServerError)
			return
		}

		newOrganizer := database.CreateOrganizerParams{
			OrganizerName: orgName,
			Description:   orgDesc,
			ImgUrl:        imgPath,
		}

		if err := s.db.CreateOrganizer(context.Background(), newOrganizer); err != nil {
			http.Error(w, "Error saving data to database.", http.StatusInternalServerError)
			return
		}

		// dashboard.CreateOrganizer(pageData, true, organizers).Render(context.Background(), w)
		w.Header().Set("HX-Refresh", "true")
	}
}

func (s *Server) HandleDeleteOrganizer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		vars := mux.Vars(r)
		idParam := vars["id"]

		id, err := strconv.Atoi(idParam)
		if err != nil {
			slog.Error("error converting id")
			return
		}

		if err := s.db.DeleteOrganizer(context.Background(), int32(id)); err != nil {
			http.Error(w, "error delete organizer.", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) saveEvent(ctx context.Context, event database.CreateEventParams, tags []string) error {
	var postgresDate time.Time
	fmt.Printf("%v\ntype%T", event.Date, event.Date)
	if event.Date.Valid {
		postgresDate = event.Date.Time
	} else {
		return fmt.Errorf("invalid date")
	}

	newEvent, err := s.db.CreateEvent(ctx, database.CreateEventParams{
		Title:       event.Title,
		Description: event.Description,
		Date:        pgtype.Date{Time: postgresDate, Valid: true},
		Freq:        event.Freq,
		Organizer:   event.Organizer,
		Imgpath:     event.Imgpath,
		Userid:      event.Userid,
	})
	if err != nil {
		slog.Error("Failed to create event",
			"error", err,
			"error_type", fmt.Sprintf("%T", err))

		if pgErr, ok := err.(*pgconn.PgError); ok {
			slog.Error("PostgreSQL error details",
				"code", pgErr.Code,
				"message", pgErr.Message,
				"detail", pgErr.Detail,
				"hint", pgErr.Hint)
		}

		return err
	}

	// Process Tags
	for _, tag := range tags {
		var tagID int32
		tagResult, err := s.db.FindTagByName(ctx, tag)
		if err == pgx.ErrNoRows {
			tagResult, err = s.db.CreateTag(ctx, tag)
		}
		if err != nil {
			return err
		}

		tagID = tagResult.ID

		// Link Event to Tag
		err = s.db.LinkEventToTag(ctx, database.LinkEventToTagParams{
			Eventid: newEvent.ID,
			Tagid:   tagID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
