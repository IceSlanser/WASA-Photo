package utils

import "regexp"

func IsLegal(s string) (bool, error) {
	if len(s) < 3 || len(s) > 16 {
		return false, nil
	}

	pattern := `^.*?$`
	return regexp.MatchString(pattern, s)
}
