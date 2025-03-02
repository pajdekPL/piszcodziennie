package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"piszcodziennie/internal/database"
	"time"
)

func (cfg *Config) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := cfg.firebaseClient.VerifyIDToken(r.Context(), req.Token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token expired or not valid", err)
		return
	}
	emailVerified, ok := token.Claims["email_verified"]

	if !ok {
		respondWithError(w, http.StatusInternalServerError, "", fmt.Errorf("email_verified claim not found in token"))
	}

	if emailVerified != true {
		respondWithError(w, http.StatusUnauthorized, "Email is not verified", nil)
	}
	fmt.Println(token.UID)
	user, err := cfg.db.GetUser(r.Context(), token.UID)
	if !user.EmailVerified {
		_, err := cfg.db.SetUserEmailVerified(r.Context(), database.SetUserEmailVerifiedParams{
			ID:            token.UID,
			EmailVerified: true,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "", err)
			return
		}
	}
	fmt.Printf("%v+\n", user)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}
	// Set the token in an HTTP-only, secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    req.Token,
		HttpOnly: true,
		Secure:   true, // Ensures it only works over HTTPS in production
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour), // 1-hour expiry
	})

	renderTemplate(w, "dashboard.html", map[string]string{"email": token.Claims["email"].(string)})

}
