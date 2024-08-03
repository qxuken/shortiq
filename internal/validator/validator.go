package validator

import (
	"errors"
	netUrl "net/url"

	"github.com/qxuken/short/internal/config"
	"github.com/qxuken/short/internal/db"
)

func ValidateRedirectUrl(conf *config.Config, url string, touched bool) error {
	if !touched && url == "" {
		return nil
	}
	u, err := netUrl.Parse(url)
	switch {
	case err != nil:
		return errors.New("Must be url")
	case u.Scheme == "":
		return errors.New("Must contain valid scheme")
	case !(u.Scheme == "http" || u.Scheme == "https"):
		return errors.New("Scheme must be either http or https")
	case u.Host == "":
		return errors.New("Must contain valid host")
	case u.Host == conf.PublicUrl.Host:
		return errors.New("Cannot create short url onto itself")
	}
	return nil
}

func ValidateShortHandle(db db.DB, url string) error {
	l := len(url)
	if l < 5 {
		return errors.New("Handle must be at least 5 characters long")
	}

	if l > 64 {
		return errors.New("Handle must be at max 64 characters long")
	}

	for _, c := range url {
		if c == '-' || c == '_' || (c >= '0' && c <= '9') || (c >= 'A' && c <= 'z') {
			continue
		}

		return errors.New("Must only contain numbers, english letters or symbols: -,_")
	}

	_, err := db.GetLink(url)
	if err == nil {
		return errors.New("Handle is already taken")
	}

	return nil
}
