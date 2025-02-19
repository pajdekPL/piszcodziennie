package utils

import (
	"encoding/json"
	"errors"
	"regexp"
)

type SupabaseError struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"msg"`
}

func ExtractSupabaseError(errMsg string) (SupabaseError, error) {
	var supaErr SupabaseError

	// Find the JSON part using regex
	re := regexp.MustCompile(`\{.*\}`)
	jsonPart := re.FindString(errMsg)

	if jsonPart == "" {
		return supaErr, errors.New("no valid JSON found in error message")
	}

	// Parse the extracted JSON
	if err := json.Unmarshal([]byte(jsonPart), &supaErr); err != nil {
		return supaErr, err
	}

	return supaErr, nil
}
