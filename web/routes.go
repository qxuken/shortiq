package web

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	servertiming "github.com/mitchellh/go-server-timing"

	"github.com/qxuken/short/internal/config"
	dbModule "github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/shortener"
	"github.com/qxuken/short/internal/validator"
	"github.com/qxuken/short/web/template/component"
	"github.com/qxuken/short/web/template/page"
)

func WebRouter(conf *config.Config, db dbModule.DB) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return servertiming.Middleware(h, nil)
		})
		if conf.Verbose {
			r.Use(middleware.RequestID)
			r.Use(middleware.Logger)
		} else {
			r.Use(middleware.Compress(3))
		}

		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "./assets"))
		fileServer(r, "/assets", filesDir)

		r.Get("/", templ.Handler(page.Index(conf.Verbose, conf.PublicUrlStr, "", "")).ServeHTTP)

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			vt := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			url := form.Get("redirect_url")
			short := form.Get("short_url")
			shortType := form.Get("short_type")
			vt.Stop()

			ut := timing.NewMetric("validating_url").Start()
			urlErr := validator.ValidateRedirectUrl(conf, url, true)
			ut.Stop()

			st := timing.NewMetric("creating_or_validating_short_url").Start()
			var shortErr error
			if shortType == "custom" {
				shortErr = validator.ValidateShortHandle(db, short)
			} else if urlErr == nil {
				short, shortErr = shortener.ShortUrlChecked(db, conf.HandleLen)
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

				c := component.CreateLink(conf.PublicUrlStr, shortType, urlErrStr, shortErrStr)
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

			fullShort := fmt.Sprintf("%v/u/%v", conf.PublicUrlStr, short)
			statsShort := "/s/" + short
			w.Header().Add("HX-Push-Url", statsShort)
			templ.Handler(component.LinkStats(templ.SafeURL(fullShort), "Your link is ready")).ServeHTTP(w, r)
		})

		r.Post("/f/generated", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			vt := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			url := form.Get("redirect_url")
			vt.Stop()

			ut := timing.NewMetric("validating_url").Start()
			urlErr := validator.ValidateRedirectUrl(conf, url, false)
			var urlErrStr string
			if urlErr != nil {
				urlErrStr = urlErr.Error()
			} else {
				urlErrStr = ""
			}
			ut.Stop()

			c := component.CreateLink(conf.PublicUrlStr, "generated", urlErrStr, "")
			templ.Handler(c).ServeHTTP(w, r)
			return

		})
		r.Post("/f/custom", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			vt := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			url := form.Get("redirect_url")
			vt.Stop()

			ut := timing.NewMetric("validating_url").Start()
			urlErr := validator.ValidateRedirectUrl(conf, url, false)
			var urlErrStr string
			if urlErr != nil {
				urlErrStr = urlErr.Error()
			} else {
				urlErrStr = ""
			}
			ut.Stop()

			c := component.CreateLink(conf.PublicUrlStr, "custom", urlErrStr, "")
			templ.Handler(c).ServeHTTP(w, r)
			return

		})

		r.Post("/f/redirect_url", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			et := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			url := form.Get("redirect_url")
			et.Stop()

			vt := timing.NewMetric("validate_values").Start()
			err := validator.ValidateRedirectUrl(conf, url, true)
			var errStr string
			if err != nil {
				errStr = err.Error()
			} else {
				errStr = ""
			}
			vt.Stop()

			templ.Handler(component.RedirectUrlInput(errStr)).ServeHTTP(w, r)
		})

		r.Post("/f/short_url", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			et := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			short := form.Get("short_url")
			et.Stop()

			vt := timing.NewMetric("validate_values").Start()
			err := validator.ValidateShortHandle(db, short)
			var errStr string
			if err != nil {
				errStr = err.Error()
			} else {
				errStr = ""
			}
			vt.Stop()

			templ.Handler(component.ShortUrlInput(errStr)).ServeHTTP(w, r)
		})

		r.Get("/s/{short_url}", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			st := timing.NewMetric("extract_url").Start()
			short := r.PathValue("short_url")
			st.Stop()

			fullShort := fmt.Sprintf("%v/u/%v", conf.PublicUrlStr, short)
			c := page.Stats(conf.Verbose, templ.SafeURL(fullShort), "")
			templ.Handler(c).ServeHTTP(w, r)
		})
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
