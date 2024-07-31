package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gocarina/gocsv"
	mdb "github.com/qxuken/short/internal/db"
)

func ExportRouter(db mdb.DB) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/a", ExportRedirectAnalyticsCsv(db))
		r.Get("/l", ExportRedirectLinksCsv(db))
	}
}

func ExportRedirectAnalyticsCsv(db mdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := db.GetLinkAnalytics()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Add("Content-Disposition", `attachment; filename="redirect_analytics.csv"`)
		gocsv.Marshal(data, w)
	}
}

func ExportRedirectLinksCsv(db mdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := db.GetLinks()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Add("Content-Disposition", `attachment; filename="redirect-links.csv"`)
		gocsv.Marshal(data, w)
	}
}
