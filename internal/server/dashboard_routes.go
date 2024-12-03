package server

import (
	"bayside-buzz/cmd/web/pages/dashboard"
	"bayside-buzz/internal/database"
	"bayside-buzz/internal/domain"
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Dashboard
func (s *Server) DashboardPage(w http.ResponseWriter, r *http.Request) {
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
		totalOrganizers, _ := s.db.CountOrganizers(ctx)

		results := domain.NewResults(totalOrganizers)

		dashboard.Dashboard(pageData, results).Render(context.Background(), w)
	}
}

func (s *Server) CreateEventPage(w http.ResponseWriter, r *http.Request) {
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
		dashboard.CreateEvent(pageData).Render(context.Background(), w)
	}

}

// Organizer name is a unique value, make sure to handle error
func (s *Server) CreateOrganizerPage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

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

	if r.Method == "GET" {
		organizers, _ := s.db.GetOrganizers(ctx)

		dashboard.CreateOrganizer(pageData, &organizers).Render(context.Background(), w)
	}

	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			slog.Error("error parsing registration form\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			w.Header().Set("hx-refresh", "true")
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

		imgPath, err := fileUpload(file, *fileHeader)
		if err != nil {
			http.Error(w, "Unable to create file.", http.StatusInternalServerError)
			return
		}

		newOrganizer := database.CreateOrganizerParams{
			OrganizerName: orgName,
			Description:   orgDesc,
			ImgUrl:        imgPath,
		}

		if err := s.db.CreateOrganizer(ctx, newOrganizer); err != nil {
			http.Error(w, "Error saving data to database.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("hx-refresh", "true")
	}
}

func (s *Server) HandleDeleteOrganizer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method == "DELETE" {
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

func fileUpload(file multipart.File, fileHeader multipart.FileHeader) (string, error) {
	// create file
	orgDir := "/assets/images/organizers/"
	basePath := "cmd/web"
	// baseDir :="../../cmd/web"

	sanitizedFileName := filepath.Base(fileHeader.Filename)
	if err := os.MkdirAll(basePath+orgDir, os.ModePerm); err != nil {
		slog.Error("Failed to create directories", err)
		return "", err
	}

	imgPath := strings.Join([]string{orgDir, sanitizedFileName}, "")

	imgOut, err := os.Create(basePath + imgPath)
	if err != nil {
		slog.Error("err with path too", err)
		return "", err
	}
	defer imgOut.Close()

	_, err = io.Copy(imgOut, file)
	if err != nil {
		slog.Error("err here too", err)
		return "", nil
	}

	return imgPath, nil
}
