package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gorilla/csrf"

	servertiming "github.com/mitchellh/go-server-timing"

	"github.com/qxuken/short/internal/auth"
	"github.com/qxuken/short/internal/config"
	dbModule "github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/shortener"
	"github.com/qxuken/short/internal/validator"
	"github.com/qxuken/short/web/template/component"
	"github.com/qxuken/short/web/template/page"
)

func WebRouter(conf *config.Config, db dbModule.DB) func(chi.Router) {
	return func(r chi.Router) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "./assets"))
		fileServer(r, "/assets", filesDir)

		r.Group(pages(conf, db))
	}
}

func pages(conf *config.Config, db dbModule.DB) func(chi.Router) {
	return func(r chi.Router) {

		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "app.conf.verbose", conf.Verbose)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
		r.Use(func(next http.Handler) http.Handler {
			return servertiming.Middleware(next, nil)
		})
		if conf.Verbose {
			r.Use(middleware.RequestID)
			r.Use(middleware.Logger)
		} else {
			r.Use(middleware.Compress(3))
		}

		r.Use(auth.CokieAuthMiddleware(conf))
		r.Use(csrf.Protect(conf.AppSecret, csrf.ErrorHandler(templ.Handler(page.CSRFError()))))

		r.Group(authorizedRouter(conf, db))
		r.Group(unauthorizedRouter(conf))

	}
}

func authorizedRouter(conf *config.Config, db dbModule.DB) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(authorizedOnly)

		r.Get("/", templ.Handler(page.Index(conf.PublicUrlStr, "", "")).ServeHTTP)

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
				templ.Handler(page.ServerError()).ServeHTTP(w, r)
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
			c := page.Stats(templ.SafeURL(fullShort), "")
			templ.Handler(c).ServeHTTP(w, r)
		})
	}
}

func unauthorizedRouter(conf *config.Config) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/login", templ.Handler(page.Auth("")).ServeHTTP)
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			timing := servertiming.FromContext(r.Context())

			vt := timing.NewMetric("extracting_values").Start()
			r.ParseForm()
			form := r.PostForm
			token := form.Get("token")
			vt.Stop()

			ut := timing.NewMetric("validating_token").Start()
			ok, tokenValidationErr := auth.VerifyHash(conf, []byte(token))
			ut.Stop()

			if ok && tokenValidationErr == nil {
				cookie := http.Cookie{
					Name:     "authToken",
					Value:    token,
					Path:     "/",
					MaxAge:   3600,
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteLaxMode,
				}
				http.SetCookie(w, &cookie)
				if _, ok := r.Header["Hx-Request"]; ok {
					w.Header().Add("HX-Redirect", "/")
					w.Header().Add("HX-Replace-Url", "/")
				} else {
					http.Redirect(w, r, "/", 307)
				}
				return
			}

			c := component.AuthForm("Invalid token")
			templ.Handler(c).ServeHTTP(w, r)
			return
		})
	}
}

func authorizedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isAuthorized, ok := ctx.Value("isAuthorized").(bool)
		if !ok || !isAuthorized {
			http.Redirect(w, r, "/login", 302)
			return
		}
		next.ServeHTTP(w, r)
	})
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
