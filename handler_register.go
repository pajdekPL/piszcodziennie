package main

import (
	"net/http"
	"piszcodziennie/internal/database"
	"strings"

	"firebase.google.com/go/auth"
)


func (cfg *Config) registerHandler(w http.ResponseWriter, req *http.Request) {
    if err := req.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    email := req.FormValue("email")
    password := req.FormValue("password")

    if email == "" || password == "" {
        http.Error(w, "Missing email or password", http.StatusBadRequest)
        return
    }

	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)

	user, err := cfg.firebaseClient.CreateUser(req.Context(),params)

	if err != nil {
		if strings.Contains(err.Error(), "EMAIL_EXISTS") {
			w.WriteHeader(http.StatusOK)
			renderTemplate(w, "register_failed.html", map[string]string{"error": "Email already exists"})
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		renderTemplate(w, "register_failed.html", map[string]string{"error": err.Error()})
		return
	}

	link, err:= cfg.firebaseClient.EmailVerificationLink(req.Context(), user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating email verification link", err)
		return
	}

	err = cfg.sendEmail(user.Email, link)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error sending email", err)
		return
	}

	_, err = cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		ID:    user.UID,
		Email: email,
	})

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			respondWithError(w, http.StatusBadRequest, "User already exists", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}
	renderTemplate(w, "register_success.html", map[string]string{"email": user.Email})
}