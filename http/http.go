package http

import (
	"net/http"

	"github.com/nasfiles/api"
)

type ApiHandler func(w http.ResponseWriter, r *http.Request, c *api.Config)

func Wrap(h ApiHandler, c *api.Config) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			code int
			err  error
		)

		h(w, r, c)
	}
}
