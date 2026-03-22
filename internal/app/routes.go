package app

import (
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	"github.com/qxuken/short/internal/api"
	"github.com/qxuken/short/internal/auth"
	"github.com/qxuken/short/internal/config"
	mdb "github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/redirect"
	"github.com/qxuken/short/web"
)

func exportV1(mainDb mdb.MainDb, auxDb mdb.AuxiliaryDB) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/analytics", api.ExportRedirectAnalyticsCsv(auxDb))
		r.Get("/links", api.ExportRedirectLinksCsv(mainDb))
	}
}

func apiV1Router(conf *config.Config, mainDb mdb.MainDb, auxDb mdb.AuxiliaryDB) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Use(cors.AllowAll().Handler)
		r.Use(auth.HeaderAuthMiddleware(conf))
		r.Use(api.AuthorizedOnly)

		r.Route("/export", exportV1(mainDb, auxDb))
		r.Post("/short", api.CreateShortUrlHandler(conf, mainDb))
		r.Get("/stats", api.GetAllStats(mainDb, auxDb))
		r.Get("/stats/{short}", api.GetLinkStats("", auxDb, mainDb))
	}
}

func appRouter(conf *config.Config, mainDb mdb.MainDb, auxDb mdb.AuxiliaryDB, r chi.Router) {
	r.Get("/u/{short}", redirect.RedirectRoute(mainDb, auxDb))
	r.Route("/api/v1", apiV1Router(conf, mainDb, auxDb))
	r.Group(web.WebRouter(conf, mainDb, auxDb))
}

func HWAddr(conf *config.Config) string {
	return path.Join(conf.DataPath, "main.db")
}
