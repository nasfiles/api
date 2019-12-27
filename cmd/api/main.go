package main

import (
	"flag"
	"log"
	"runtime"
	"strconv"

	"github.com/nasfiles/api"
	"github.com/nasfiles/api/bolt"
	h "github.com/nasfiles/api/http"

	boltdb "github.com/boltdb/bolt"
	"github.com/fatih/color"
)

func main() {
	// run the api server with the maximum number of cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse flags
	cfgPath := flag.String("config", "", "JSON config file")
	development := flag.Bool("development", false, "Development mode")
	dumpDB := flag.Bool("dump", false, "Dump database")

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

	// if the dump flag is used, dump the database info and don't spawn the http server
	if *dumpDB {
		bolt.Dump(db, true)
		return
	}

	h.Serve(cfg)
}
