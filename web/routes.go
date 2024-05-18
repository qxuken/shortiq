package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/web/template/page"
)

func WebRouter(db db.DB) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Compress(3))

		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "./assets"))
		FileServer(r, "/assets", filesDir)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			page.Index().Render(r.Context(), w)
		})

	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
