package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	"github.com/qxuken/short/internal/api"
	"github.com/qxuken/short/internal/config"
	mdb "github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/redirect"
	"github.com/qxuken/short/web"
)

func exportV1(db mdb.DB) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/analytics", api.ExportRedirectAnalyticsCsv(db))
		r.Get("/links", api.ExportRedirectLinksCsv(db))
	}
}

func apiV1Router(conf *config.Config, db mdb.DB) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Use(cors.AllowAll().Handler)

		r.Route("/export", exportV1(db))
		r.Post("/short", api.CreateShortUrlHandler(conf, db))
	}
}

func appRouter(conf *config.Config, db mdb.DB, r chi.Router) {
	r.Get("/u/{short}", redirect.RedirectRoute(db))
	r.Route("/api/v1", apiV1Router(conf, db))
	r.Group(web.WebRouter(conf, db))
}
