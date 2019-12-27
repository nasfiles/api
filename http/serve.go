package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fatih/color"
	"github.com/nasfiles/api"

	"github.com/gorilla/mux"
)

//Serve sets all routes with mux and spawns the HTTP server
func Serve(cfg *api.Config) {
	r := mux.NewRouter()

	r.HandleFunc("/users", Wrap(UserAdd, cfg)).Methods("POST")
	r.HandleFunc("/users/{uid}", Wrap(UserGet, cfg)).Methods("GET")
	r.HandleFunc("/users/{uid}", Wrap(UserDelete, cfg)).Methods("DELETE")

	// Start HTTP WebDav Server
	color.HiCyan("\n\nStarting WebDAV server...")
	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), r); err != nil {
		log.Fatalf("Error starting HTTP WebDAV server: %v", err)
	}
}
