package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/qxuken/short/internal/config"
)

func HeaderAuthMiddleware(conf *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAuthorized := false
			if bearerToken, ok := r.Header["Authorization"]; ok {
				token, _ := strings.CutPrefix(bearerToken[0], "Bearer ")
				isAuthorized, _ = VerifyHash(conf, []byte(token))
			}
			if conf.Verbose {
				log.Println("isAuthorized", isAuthorized)
			}
			ctx := context.WithValue(r.Context(), "isAuthorized", isAuthorized)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func CokieAuthMiddleware(conf *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAuthorized := false
			if authCokie, err := r.Cookie("authToken"); err == nil {
				isAuthorized, _ = VerifyHash(conf, []byte(authCokie.Value))
			}
			if conf.Verbose {
				log.Println("isAuthorized", isAuthorized)
			}
			ctx := context.WithValue(r.Context(), "isAuthorized", isAuthorized)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
