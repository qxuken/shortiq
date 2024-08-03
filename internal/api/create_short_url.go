package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/qxuken/short/internal/config"
	mdb "github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/shortener"
	"github.com/qxuken/short/internal/validator"
)

type CreateShortUrl struct {
	RedirectUrl  string `json:"redirectUrl"`
	ShortUrl     string `json:"shortUrl"`
	FullShortUrl string `json:"fullShortUrl"`
}

func (a *CreateShortUrl) Bind(r *http.Request) error {
	if a.RedirectUrl == "" {
		return errors.New("missing required 'redirectUrl' field")
	}

	return nil
}

func (rd *CreateShortUrl) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func CreateShortUrlHandler(conf *config.Config, db mdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &CreateShortUrl{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		if err := validator.ValidateRedirectUrl(conf, data.RedirectUrl, true); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		if data.ShortUrl == "" {
			short, shortErr := shortener.ShortUrlChecked(db, conf.HandleLen)
			if shortErr != nil {
				render.Render(w, r, ErrInternalError(shortErr))
				return
			}
			data.ShortUrl = short
		} else {
			if err := validator.ValidateShortHandle(db, data.ShortUrl); err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
		}
		if err := db.SetLink(data.RedirectUrl, data.ShortUrl); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		data.FullShortUrl = fmt.Sprintf("%v/u/%v", conf.PublicUrlStr, data.ShortUrl)
		render.Status(r, http.StatusCreated)
		render.Render(w, r, data)
	}
}
