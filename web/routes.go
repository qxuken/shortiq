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

var exampleShort = internal.ShortUrl(exampleUrl)

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
			short, _ := internal.ShortUrlChecked(db, exampleUrl)
			st.Stop()

			if r.Header.Get("Hx-Request") == "true" {
				c := component.CreateLink(conf.PublicUrl, exampleUrl, short, "generated", "", "")
				templ.Handler(c).ServeHTTP(w, r)
				return
			}

			c := page.Index(conf.Debug, conf.PublicUrl, exampleUrl, short, "", "")
			templ.Handler(c).ServeHTTP(w, r)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			vt := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			url := form.Get("url")
			short := form.Get("short")
			shortType := form.Get("short_type")
			vt.Stop()

			ut := timing.NewMetric("validating_url").Start()
			urlErr := internal.ValidateRedirectUrl(url)
			ut.Stop()

			st := timing.NewMetric("validating_short_url").Start()
			shortErr := internal.ValidateShortHandle(db, short)
			if shortType != "custom" && urlErr == nil && shortErr != nil {
				short, shortErr = internal.ShortUrlChecked(db, url)
			}
			st.Stop()

			if urlErr != nil || shortErr != nil {
				var urlErrStr string
				if urlErr != nil {
					urlErrStr = urlErr.Error()
				} else {
					urlErrStr = ""
				}

				var shortErrStr string
				if shortErr != nil {
					shortErrStr = shortErr.Error()
				} else {
					shortErrStr = ""
				}

				c := component.CreateLink(conf.PublicUrl, url, short, shortType, urlErrStr, shortErrStr)
				templ.Handler(c).ServeHTTP(w, r)
				return
			}

			dt := timing.NewMetric("save").Start()
			err := db.SetLink(url, short)
			dt.Stop()

			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}

			fullShort := conf.PublicUrl + "/u/" + short
			templ.Handler(component.SuccessfullyCreated(templ.SafeURL(fullShort))).ServeHTTP(w, r)
		})

		r.Post("/f/url", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			et := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			url := form.Get("url")
			et.Stop()

			vt := timing.NewMetric("validate_values").Start()
			err := internal.ValidateRedirectUrl(url)
			var errStr string
			if err != nil {
				errStr = err.Error()
			} else {
				errStr = ""
			}
			vt.Stop()

			templ.Handler(component.RedirectUrlInput(url, errStr)).ServeHTTP(w, r)
		})

		r.Post("/f/short", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			et := timing.NewMetric("extracting_url").Start()
			r.ParseForm()
			form := r.PostForm
			url := form.Get("url")
			short := form.Get("short")
			shortType := form.Get("short_type")
			et.Stop()

			st := timing.NewMetric("short_url").Start()

			var err error
			if shortType == "custom" {
				err = internal.ValidateShortHandle(db, short)
				if err != nil {
					short, err = internal.ShortUrlChecked(db, url)
				}
			} else {
				short, err = internal.ShortUrlChecked(db, url)
			}
			st.Stop()

			var errStr string
			if err != nil {
				errStr = err.Error()
			} else {
				errStr = ""
			}

			templ.Handler(component.ExampleUrl(conf.PublicUrl, short, errStr)).ServeHTTP(w, r)
		})

		r.Post("/f/custom", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			et := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			short := form.Get("short")
			et.Stop()

			vt := timing.NewMetric("validate_values").Start()
			err := internal.ValidateShortHandle(db, short)
			var errStr string
			if err != nil {
				errStr = err.Error()
			} else {
				errStr = ""
			}
			vt.Stop()

			templ.Handler(component.ShortUrlInput(short, errStr)).ServeHTTP(w, r)
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
