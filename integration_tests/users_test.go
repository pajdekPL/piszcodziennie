//go:build integration
// +build integration

package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/health")
	assert.NoError(t, err, "Health check should not return an error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Health check should return 200 OK")
}

func TestUserRegistration(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:8080/register", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	assert.NoError(t, err, "API should be reachable")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400 for invalid request")
}
