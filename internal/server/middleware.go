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
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userId, ok := session.Values["userId"]
			if !ok {
				slog.Warn("No userId in session", userId)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}


			// Convert userId to int64
			var userIdInt64 int64
			switch v := userId.(type) {
			case int:
				userIdInt64 = int64(v)
			case int64:
				userIdInt64 = v
			default:
				slog.Error("Invalid userId type in session", "type", v)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			u := func() *database.User {
				ctxB := context.Background()
				user, err := s.db.FindUser(ctxB, userIdInt64)
				if err != nil {
					slog.Error("Error finding user in DB", err)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return nil
				}
				return &user
			}()

			slog.Info("Authenticated user", "userID", u.ID)
			c := context.WithValue(r.Context(), contextKeyUser, u)
			next.ServeHTTP(w, r.WithContext(c))
		},
	)
}
