package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func GetTokenFromHeader(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header not present")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("token not present")
	}
	return token, nil
}
