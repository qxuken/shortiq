package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

func AuthorizedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isAuthorized, ok := ctx.Value("isAuthorized").(bool)
		if !ok || !isAuthorized {
			render.Render(w, r, ErrUnauthorizedRequest(errors.New("Must include valid token")))
			return
		}
		next.ServeHTTP(w, r)
	})
}
