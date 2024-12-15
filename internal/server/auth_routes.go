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

func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Get(r, "login")
	if err != nil {
		slog.Error("Error with getting session info.\n", err)
		return
	}

	delete(session.Values, "userId")

	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/login")
}

func (s *Server) LoginPage(w http.ResponseWriter, r *http.Request) {
	var (
		title = strings.Join([]string{SITE_NAME, "Login"}, " - ")
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
			w.Header().Add("HX-Refresh", "true")
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := s.authenticateUser(email, password)
		if err != nil {
			pages.Login(pageData, true).Render(context.Background(), w)
			return
		}

		// CREATES THE SESSION
		session, err := s.store.Get(r, "login")
		if err != nil {
			slog.Error("Failed to get session", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		session.Values["userId"] = user.ID

		if err := session.Save(r, w); err != nil {
			slog.Error("Failed to save session", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("HX-Redirect", "/dashboard")
	}
}

func (s *Server) authenticateUser(email, password string) (*database.User, error) {
	// Authenticate user
	user, err := s.db.GetUser(context.Background(), email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Server) RegisterPage(w http.ResponseWriter, r *http.Request) {
	var (
		title = strings.Join([]string{SITE_NAME, "Register"}, " - ")
		url   = r.Host
	)
	const (
		description = "Register to Bayside Breeze"
		pageType    = "website"
		image       = "" // get an image
	)

	pageData := domain.NewPageData(SITE_NAME, title, description, pageType, image, url)

	if r.Method == "GET" {
		count, err := s.db.CountUsers(context.Background())
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
			w.Header().Add("HX-Refresh", "true")
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

		if err = s.db.CreateUser(context.Background(), userParams); err != nil {
			slog.Error("error creating user: \n", err)
			return
		}

		// just testing redirects
		w.Header().Add("HX-Redirect", "/login")
	}
}
