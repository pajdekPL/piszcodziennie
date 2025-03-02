package main

import (
	"fmt"
	"net/http"
)

type ReturnProtected struct {
	Message string `json:"Message"`
}

func (cfg *Config) protectedHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(userIDKey)
	email := ctx.Value(emailKey)
	renderTemplate(w, "dashboard.html", map[string]string{"email": email.(string), "protected": fmt.Sprintf("twoj stary pijany id: %s", userID.(string))})
}
