package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fatih/color"
	"github.com/gorilla/mux"

	"github.com/nasfiles/nasfilesapi"
)

// Serve sets all routes with mux and spawns the HTTP server
func Serve(cfg *nasfilesapi.Config) {
	r := mux.NewRouter()

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", Wrap(AuthLogin, cfg)).Methods("POST", "OPTIONS")

	r.HandleFunc("/users", Wrap(UserAdd, cfg)).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/{uid}", Wrap(UserGet, cfg)).Methods("GET", "OPTIONS")
	r.HandleFunc("/users/{uid}", Wrap(UserDelete, cfg)).Methods("DELETE", "OPTIONS")

	// Start HTTP WebDav Server
	color.HiCyan("\n\nStarting WebDAV server...")
	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), r); err != nil {
		log.Fatalf("Error starting HTTP WebDAV server: %v", err)
	}
}
