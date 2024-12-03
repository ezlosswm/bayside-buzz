package server

import (
	"bayside-buzz/cmd/web/pages"
	"bayside-buzz/internal/database"
	"bayside-buzz/internal/domain"
	"context"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) LoginPage(w http.ResponseWriter, r *http.Request) {
	var (
		title = strings.Join([]string{"Login", SITE_NAME}, " - ")
		url   = r.Host
	)
	const (
		description = "Access your Bayside Breeze account to manage the events happening in Corozal Town, Belize."
		pageType    = "website"
		image       = "" // get an image
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	if r.Method == "GET" {
		pages.Login(pageData, false).Render(context.Background(), w)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			slog.Error("error parsing registration form\n", err)
			w.Header().Set("hx-refresh", "true")
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if err := s.authenticateUser(email, password); err != nil {
			slog.Error("error with authenticating the user", err)
			pages.Login(pageData, true).Render(context.Background(), w)
			return
		}

		// just testing redirects
		w.Header().Add("hx-redirect", "/dashboard")
	}
}

func (s *Server) authenticateUser(email, password string) error {
	// Authenticate user
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return err
	}

	return nil
}

func (s *Server) RegisterPage(w http.ResponseWriter, r *http.Request) {
	var (
		title = strings.Join([]string{"Register", SITE_NAME}, " - ")
		url   = r.Host
	)
	const (
		description = "Register to Bayside Breeze"
		pageType    = "website"
		image       = "" // get an image
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	if r.Method == "GET" {
		ctx := context.Background()
		count, err := s.db.CountUsers(ctx)
		if err != nil {
			slog.Error("error counting user: \n", err)
			return
		}

		if count == 1 {
			pages.Register(pageData, true).Render(context.Background(), w)
		}

		pages.Register(pageData, false).Render(context.Background(), w)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			slog.Error("error parsing registration form\n", err)
			w.Header().Set("hx-refresh", "true")
			return
		}

		// Get form data
		fullName := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Authenticate User
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("error generating hash: \n", err)
			return
		}

		userParams := database.CreateUserParams{
			Name:         fullName,
			Email:        email,
			PasswordHash: string(passwordHash),
		}

		ctx := context.Background()
		if err = s.db.CreateUser(ctx, userParams); err != nil {
			slog.Error("error creating user: \n", err)
			return
		}

		// just testing redirects
		w.Header().Add("hx-redirect", "/login")
	}
}
