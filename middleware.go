package main

import (
	"context"
	"net/http"
)

type contextKey string
const userIDKey contextKey = "user_id"
const emailKey contextKey = "email"


func (cfg *Config) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err :=  r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Missing Authorization Cookie", http.StatusUnauthorized)
			return
		}

		token, err := cfg.firebaseClient.VerifyIDToken(r.Context(), authCookie.Value)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token", err)
			return
		}

		result, err := cfg.db.GetUser(r.Context(), token.UID)

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found", err)
			return
		}

		if !result.EmailVerified {
			respondWithError(w, http.StatusUnauthorized, "Email is not verified", nil)
			return
		}

		if result.IsBanned {
			http.Error(w, "User is deactivated. Contact support.", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, result.ID)
		ctx = context.WithValue(ctx, emailKey, result.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
