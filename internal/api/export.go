package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gocarina/gocsv"
	mdb "github.com/qxuken/short/internal/db"
)

func ExportRedirectAnalyticsCsv(db mdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := db.GetLinkAnalytics()
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}
		w.Header().Add("Content-type", "text/csv")
		w.Header().Add("Content-Disposition", `attachment; filename="redirect_analytics.csv"`)
		gocsv.Marshal(data, w)
	}
}

func ExportRedirectLinksCsv(db mdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := db.GetLinks()
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}
		w.Header().Add("Content-type", "text/csv")
		w.Header().Add("Content-Disposition", `attachment; filename="redirect_links.csv"`)
		gocsv.Marshal(data, w)
	}
}
