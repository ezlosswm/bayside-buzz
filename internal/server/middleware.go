package server

import (
	"bayside-buzz/internal/database"
	"context"
	"log/slog"
	"net/http"
)

type contextKey string

const contextKeyUser = contextKey("user")

func (s *Server) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			session, err := s.store.Get(r, "login")
			if err != nil {
				slog.Error("error getting user", err)
				redirect(w, r)
				return
			}

			userId, ok := session.Values["userId"]
			if !ok {
				slog.Warn("No userId in session", userId)
				redirect(w, r)
				return
			}

			// Convert userId to int32
			var userIdInt32 int32
			switch v := userId.(type) {
			case int:
				userIdInt32 = int32(v)
			case int32:
				userIdInt32 = v
			default:
				slog.Error("Invalid userId type in session", "type", v)
				redirect(w, r)
				return
			}

			u := func() *database.User {
				user, err := s.db.FindUser(context.Background(), userIdInt32)
				if err != nil {
					slog.Error("Error finding user in DB", err)
					redirect(w, r)
					return nil
				}
				return &user
			}()

			c := context.WithValue(r.Context(), contextKeyUser, u)
			next.ServeHTTP(w, r.WithContext(c))
		},
	)
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

func redirect(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/login")
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
