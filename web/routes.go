package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	servertiming "github.com/mitchellh/go-server-timing"

	"github.com/qxuken/short/internal"
	"github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/web/template/component"
	"github.com/qxuken/short/web/template/page"
)

const exampleUrl = "https://github.com/qxuken"

func WebRouter(db db.DB, conf *internal.Config) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return servertiming.Middleware(h, nil)
		})
		r.Use(middleware.RequestID)
		if conf.Debug {
			r.Use(middleware.Logger)
		}
		r.Use(middleware.Compress(3))

		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "./assets"))
		FileServer(r, "/assets", filesDir)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			st := timing.NewMetric("short_url").Start()
			short := internal.ShortUrl(exampleUrl)
			st.Stop()

			if r.Header.Get("Hx-Request") == "true" {
				templ.Handler(component.CreateLink(conf.PublicUrl, exampleUrl, short)).ServeHTTP(w, r)
				return
			}

			templ.Handler(page.Index(conf.Debug, conf.PublicUrl, exampleUrl, short)).ServeHTTP(w, r)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			vt := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			url := form.Get("url")
			if url == "" {
				http.Error(w, http.StatusText(400), 400)
				return
			}
			var short string
			if form.Has("short") {
				short = form.Get("short")
				vt.Stop()
			} else {
				vt.Stop()
				st := timing.NewMetric("short_url").Start()
				var err error
				short, err = internal.ShortUrlChecked(db, exampleUrl)
				if err != nil {
					http.Error(w, http.StatusText(400), 400)
					return
				}
				st.Stop()
			}

			dt := timing.NewMetric("save").Start()
			_, err := db.SetLink(url, short)
			dt.Stop()

			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}

			fullShort := conf.PublicUrl + "/u/" + short
			templ.Handler(component.SuccessfullyCreated(templ.SafeURL(fullShort))).ServeHTTP(w, r)
		})

		r.Get("/f/short", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			et := timing.NewMetric("extracting_url").Start()
			url := r.URL.Query().Get("url")
			et.Stop()

			st := timing.NewMetric("short_url").Start()
			short := internal.ShortUrl(url)
			st.Stop()

			templ.Handler(component.ExampleUrl(conf.PublicUrl, short)).ServeHTTP(w, r)
		})

		r.Get("/f/custom", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			et := timing.NewMetric("extracting_url").Start()
			url := r.URL.Query().Get("url")
			et.Stop()

			st := timing.NewMetric("short_url").Start()
			short := internal.ShortUrl(url)
			st.Stop()

			templ.Handler(component.ShortUrlInput(short)).ServeHTTP(w, r)
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
