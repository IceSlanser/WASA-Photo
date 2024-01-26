package utils

import (
	"encoding/json"
	"net/http"
	"regexp"
)

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

func IsLegal(s string) (bool, error) {
	if len(s) < 3 || len(s) > 16 {
		return false, nil
	}

	pattern := `^.*?$`
	return regexp.MatchString(pattern, s)
}
