package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gocarina/gocsv"
	mdb "github.com/qxuken/short/internal/db"
)

func ExportRedirectAnalyticsCsv(auxDb mdb.AuxiliaryDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := auxDb.GetLinkAnalytics()
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}
		w.Header().Add("Content-type", "text/csv")
		w.Header().Add("Content-Disposition", `attachment; filename="redirect_analytics.csv"`)
		gocsv.Marshal(data, w)
	}
}

func ExportRedirectLinksCsv(mainDb mdb.MainDb) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := mainDb.GetLinks()
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}
		w.Header().Add("Content-type", "text/csv")
		w.Header().Add("Content-Disposition", `attachment; filename="redirect_links.csv"`)
		gocsv.Marshal(data, w)
	}
}
