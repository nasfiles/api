package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	boltdb "github.com/boltdb/bolt"
	"github.com/fatih/color"

	"github.com/nasfiles/nasfilesapi"
	"github.com/nasfiles/nasfilesapi/bolt"
	h "github.com/nasfiles/nasfilesapi/http"
)

var cfgPath string
var dbPath string
var dumpDB bool

var cfg *nasfilesapi.Config
var db *boltdb.DB

func parseFlags() {
	flag.StringVar(&cfgPath, "config", "config.json", "JSON config file")
	flag.StringVar(&dbPath, "db", "nasfilesapi.db", "Database file")
	flag.BoolVar(&dumpDB, "dump", false, "Dump database")

	flag.Parse()
}

func main() {
	// Run the api server with the maximum number of cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse flags
	parseFlags()

	var err error
	c, err := loadConfig(cfgPath)
	if err != nil {
		log.Panic(err)
	}

	// Open database file
	db, err = boltdb.Open(dbPath, 0600, nil)
	if err != nil {
		color.HiRed("Couldn't open database file.")
		log.Fatalf("Couldn't open database file: %v", err)
	}

	// Initialize buckets if the database file was just created
	bolt.Setup(db)

	cfg = &nasfilesapi.Config{
		Development: c.Development,
		Host:        c.Host,
		Port:        c.Port,
		Secure:      c.Secure,

		PrivateKey: c.PrivateKey,

		StorageRoot: c.StorageRoot,
		Services: &nasfilesapi.Services{
			Auth: &bolt.AuthService{
				DB: db,
			},
			User: &bolt.UserService{
				DB: db,
			},
		},
	}

	// Create a storage folder if it doesn't exist
	if _, err := os.Stat(cfg.StorageRoot); os.IsNotExist(err) {
		err := os.Mkdir(cfg.StorageRoot, 0755)
		if err != nil {
			color.HiRed("Couldn't create storage folder at %s.\n-")
			log.Fatalf("Couldn't storage folder: %v", err)
		}
	}

	// Print config values
	cfg.Log()

	// If the dump flag is used, dump the database info and don't spawn the http server
	if dumpDB {
		err := bolt.Dump(db, true)
		if err != nil {
			color.HiRed("Couldn't dump database information.\n")
			log.Fatalf("Couldn't dump database information: %v", err)
		}

		return
	}

	h.Serve(cfg)
}
