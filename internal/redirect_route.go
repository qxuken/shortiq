package internal

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/qxuken/short/internal/db"
)

func RedirectRoute(db db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := db.GetLink(r.PathValue("short"))
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				http.Error(w, http.StatusText(404), 404)
			default:
				http.Error(w, http.StatusText(500), 500)
			}
			return
		}
		http.Redirect(w, r, url, 302)
	}
}
