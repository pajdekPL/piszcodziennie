package utils

import (
	"net/http"
	"testing"
)

func TestGetTokenFromHeader(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "valid token",
			header:    "Bearer token",
			wantToken: "token",
			wantErr:   false,
		},
		{
			name:      "missing token",
			header:    "",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "invalid header format",
			header:    "Invalid token",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "empty token",
			header:    "Bearer ",
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := http.Request{
				Header: make(http.Header),
			}
			req.Header.Set("Authorization", tt.header)
			token, err := GetTokenFromHeader(&req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokenFromHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if token != tt.wantToken {
				t.Errorf("GetTokenFromHeader() token = %v, want %v", token, tt.wantToken)
			}
		})
	}
}
