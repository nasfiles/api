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

// APIWrapper function defines and returns a HTTP handler function with access
// to the config and all services
func APIWrapper(h APIHandler, c *nasfilesapi.Config) http.HandlerFunc {
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

func mustLogin(c *nasfilesapi.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.Auth {
			// Gets the correct user for this request.
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "401 Unauthorized", 401)

				return
			}

			credentials := &nasfilesapi.AuthGuess{
				Username: username,
				Password: password,
			}

			// Authenticate user
			err := c.Services.Auth.Login(credentials)
			if err != nil {
				color.HiRed("Passwords don't match")
				http.Error(w, "403 Unauthorized", 403)

				return
			}

			user, ok := c.Users[username]
			if !ok {
				color.HiRed("No handler found")
				http.Error(w, "Not authorized", 401)

				return
			}

			// Excerpt from RFC4918, section 9.4:
			//
			// 		GET, when applied to a collection, may return the contents of an
			//		"index.html" resource, a human-readable view of the contents of
			//		the collection, or something else altogether.
			//
			// Get, when applied to collection, will return the same as PROPFIND method.
			// if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/") {
			// 	info, err := user.Handler.FileSystem.Stat(context.TODO(), strings.TrimPrefix(r.URL.Path, "/"))
			// 	if err == nil && info.IsDir() {
			// 		r.Method = "PROPFIND"

			// 		if r.Header.Get("Depth") == "" {
			// 			r.Header.Add("Depth", "1")
			// 		}
			// 	}
			// }

			c.Users[user.Username].Handler.ServeHTTP(w, r)
			return
		}

		// Our middleware logic goes here...
		c.DefaultHandler.ServeHTTP(w, r)
	})
}
