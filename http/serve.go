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

	// API subrouter
	api := r.PathPrefix("/api").Subrouter()

	// Auth
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", APIWrapper(AuthLogin, cfg)).Methods("POST", "OPTIONS")

	// Users
	users := api.PathPrefix("/users").Subrouter()
	users.HandleFunc("/", APIWrapper(UserAdd, cfg)).Methods("POST", "OPTIONS")
	users.HandleFunc("/users/{uid}", APIWrapper(UserGet, cfg)).Methods("GET", "OPTIONS")
	users.HandleFunc("/users/{uid}", APIWrapper(UserDelete, cfg)).Methods("DELETE", "OPTIONS")

	// WebDAV
	r.PathPrefix("/").Handler(mustLogin(cfg))

	// Start HTTP WebDav Server
	color.HiCyan("Starting WebDAV server...")
	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), r); err != nil {
		log.Fatalf("Error starting HTTP WebDAV server: %v", err)
	}
}
