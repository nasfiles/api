package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/nasfiles/api"
	"github.com/nasfiles/api/bolt"

	boltdb "github.com/boltdb/bolt"
	"github.com/fatih/color"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse flags
	cfgPath := flag.String("config", "", "JSON config file")
	development := flag.Bool("development", false, "Development mode")

	// Server address
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "3000", "Server port")
	secure := flag.Bool("secure", false, "Secure connection")

	// Database and Storage
	DBpath := flag.String("db", "nasfiles.db", "Database path")
	storage := flag.String("storage", ".", "Storage base folder")

	flag.Parse()

	var cfg *api.Config
	var db *boltdb.DB
	var err error
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

		db, err = boltdb.Open(*DBpath, 0600, nil)
		if err != nil {
			color.HiRed("Couldn't open database.")
			log.Fatalf("Couldn't open database: %v", err)
		}

		// Initialize buckets if the database file was just created
		bolt.Setup(db)

		cfg = &api.Config{
			Development: *development,
			Host:        *host,
			Port:        portInt,
			Secure:      *secure,
			Services: &api.Services{
				User: &bolt.UserService{
					DB: db,
				},
			},
			StorageRoot: *storage,
		}
	}

	// Print config values
	cfg.Log()

	bolt.Dump(db)

	// Start HTTP WebDav Server
	color.HiCyan("\n\nStarting WebDAV server...")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil); err != nil {
		color.HiRed("Couldn't start HTTP WebDAV server.")
		log.Fatalf("Error starting HTTP WebDAV server: %v", err)
	}
}
