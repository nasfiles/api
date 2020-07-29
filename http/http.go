package http

import (
	"encoding/json"
	"net/http"

	"github.com/fatih/color"

	"github.com/nasfiles/nasfilesapi"
)

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// APIHandler is the type of function that should be written and used with Wrap
// to simulate standard HandlerFunc's functions,but with access to the config
type APIHandler func(w http.ResponseWriter, r *http.Request, c *nasfilesapi.Config) (int, error)

// Wrap function defines and returns a HTTP handler function with access
// to the config and all services
func Wrap(h APIHandler, c *nasfilesapi.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			code int
			err  error
		)

		defer func() {
			// if the return codes are 0 and no errors
			// it means the important data has been written into
			// the response body and the request is finished
			if code == 0 && err == nil {
				return
			}

			res := &response{
				Code: code,
			}

			if code != 0 {
				w.WriteHeader(code)
			}

			// Write JSON into response body
			data, e := json.MarshalIndent(res, "", "\t")
			if e != nil {
				color.Red(e.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)

			return
		}()

		code, err = h(w, r, c)
	}
}
