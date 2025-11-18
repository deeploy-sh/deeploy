package utils

import "regexp"

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`(?i)^[a-z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)
	return emailRegex.MatchString(e)
}
