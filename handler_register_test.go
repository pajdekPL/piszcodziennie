//go:build integration
// +build integration

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRegistration(t *testing.T) {
	config := CreateTestConfig()
	router := SetupTestRouter(config)

	invalidEmailBody := bytes.NewBuffer([]byte(`{"email": "invalid-email", "password": "securepass"}`))
	properRequest := bytes.NewBuffer([]byte(`{"email": "test1@example.com", "password": "securepass"}`))
	onlyEmailBody := bytes.NewBuffer([]byte(`{"email": "test2@example.com"}`))
	tests := []struct {
		name         string
		requestBody  *bytes.Buffer
		bodyContains string
		expectedCode int
	}{
		{
			name:         "invalid email",
			requestBody:  invalidEmailBody,
			bodyContains: "Invalid request",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "valid request",
			requestBody:  properRequest,
			bodyContains: "User registered",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "only email",
			requestBody:  onlyEmailBody,
			bodyContains: "Invalid request",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/register", tt.requestBody)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.bodyContains)
		})
	}
}

func TestUserRegistrationSameEmail(t *testing.T) {
	config := CreateTestConfig()
	router := SetupTestRouter(config)

	properRequest := bytes.NewBuffer([]byte(`{"email": "test3@example.com", "password": "securepass"}`))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", properRequest)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/register", properRequest)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
