package utils

import (
	"fmt"
	"net/url"
	"regexp"
)

func FormatStars(count int) string {
	if count >= 1000 {
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`(?i)^[a-z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)
	return emailRegex.MatchString(e)

}

func IsValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
