package utils

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"regexp"
)

// GetMyUsername gets UserName from the request body
func GetMyUsername(r *http.Request) (string, error) {
	type Temp struct {
		Username string `json:"username"`
	}
	var username Temp
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		return "", err
	}
	return username.Username, nil
}

// IsLegal verifies if the input is a standard string
func IsLegal(s string) (bool, error) {
	if len(s) < 3 || len(s) > 16 {
		return false, nil
	}

	pattern := `^.*?$`
	return regexp.MatchString(pattern, s)
}

// IsMediaFile verifies the input filename extensions
func IsMediaFile(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".mp3", ".mp4":
		return true
	}
	return false
}
