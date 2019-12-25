package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fatih/color"
	"github.com/nasfiles/api"

	"golang.org/x/net/webdav"
)

func main() {
	// parse flags
	cfgPath := flag.String("config", "", "JSON config file")
	development := flag.Bool("development", false, "Development mode")

	// Server address
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "3000", "Server port")
	secure := flag.Bool("secure", false, "Secure connection")

	// Storage
	storage := flag.String("storage", ".", "Storage base folder")

	flag.Parse()

	var cfg *api.Config
	if len(*cfgPath) > 0 {
		c, err := loadConfig(*cfgPath)
		if err != nil {
			log.Panic(err)
		}

		cfg = &api.Config{
			Development: c.Development,
			Host:        c.Host,
			Port:        c.Port,
			Secure:      c.Secure,
			StorageRoot: c.StorageRoot,
		}
	} else {
		portInt, _ := strconv.Atoi(*port)

		cfg = &api.Config{
			Development: *development,
			Host:        *host,
			Port:        portInt,
			Secure:      *secure,
			StorageRoot: *storage,
		}
	}

	cfg.Log()

	// WebDAV configuration
	srv := &webdav.Handler{
		FileSystem: webdav.Dir(*storage),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
			}
		},
	}
	http.Handle("/", srv)

	color.HiCyan("\n\nStarting WebDAV server...")
	// Start HTTP WebDav Server
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil); err != nil {
		color.HiRed("Couldn't start HTTP WebDAV server.")
		log.Fatalf("Error starting HTTP WebDAV server: %v", err)
	}
}
