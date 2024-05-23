package internal

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/golang-lru/v2"
	mdb "github.com/qxuken/short/internal/db"
)

func logRedirect(db mdb.DB, r *http.Request) {
	short := r.PathValue("short")
	c := r.Header.Get("CF-IPCountry")
	o := r.Header.Get("Referer")
	if o == "" {
		o = r.Header.Get("Origin")
	}
	ts := time.Now().Unix()

	db.LogVisit(short, mdb.LinkVisit{Country: c, Origin: o, Ts: ts})
}

func RedirectRoute(db mdb.DB) func(w http.ResponseWriter, r *http.Request) {
	cache, err := lru.New[string, string](512)
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		short := r.PathValue("short")

		defer logRedirect(db, r)

		url, ok := cache.Get(short)
		if ok {
			http.Redirect(w, r, url, 302)
			return
		}

		url, err := db.GetLink(short)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				http.Error(w, http.StatusText(404), 404)
			default:
				http.Error(w, http.StatusText(500), 500)
			}
			return
		}

		cache.Add(short, url)

		http.Redirect(w, r, url, 302)

	}

}
