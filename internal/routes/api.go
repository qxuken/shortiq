package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	"github.com/qxuken/short/internal"
	"github.com/qxuken/short/internal/api"
	mdb "github.com/qxuken/short/internal/db"
)

func ApiRouter(conf *internal.Config, db mdb.DB) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Use(middleware.Recoverer)
		r.Use(cors.AllowAll().Handler)

		r.Route("/v1", func(r chi.Router) {
			r.Post("/short", api.CreateShortUrlHandler(conf, db))
		})
	}
}
