package internal

import (
	"errors"
	"regexp"

	"github.com/qxuken/short/internal/db"
)

const (
	UrlRegexStr    = `https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	UrlShortHandle = `^[-_a-zA-Z0-9]{5,256}$`
)

func ValidateRedirectUrl(url string) error {
	if len(url) < 10 {
		return errors.New("Must be at least 10 characters long")
	}
	urlRegex := regexp.MustCompile(UrlRegexStr)
	if !urlRegex.MatchString(url) {
		return errors.New("Must be url")
	}
	return nil
}

func ValidateShortHandle(db db.DB, url string) error {
	l := len(url)
	if l < 5 {
		return errors.New("Handle must be at least 5 characters long")
	}

	if l > 256 {
		return errors.New("Handle must be at max 256 characters long")
	}

	urlRegex := regexp.MustCompile(UrlShortHandle)
	if !urlRegex.MatchString(url) {
		return errors.New("Must only contain numbers, letters or symbols: -,_")
	}

	_, err := db.GetLink(url)
	if err == nil {
		return errors.New("Handle is already taken")
	}

	return nil
}
