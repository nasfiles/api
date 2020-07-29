package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"strconv"

	boltdb "github.com/boltdb/bolt"
	"github.com/fatih/color"

	"github.com/nasfiles/nasfilesapi"
	"github.com/nasfiles/nasfilesapi/bolt"
	h "github.com/nasfiles/nasfilesapi/http"
)

var cfgPath string
var development bool
var dumpDB bool

var host string
var port string
var secure bool

var dbPath string
var storageRoot string

var cfg *nasfilesapi.Config
var db *boltdb.DB

func parseFlags() {
	flag.StringVar(&cfgPath, "config", "config.json", "JSON config file")
	flag.BoolVar(&development, "development", false, "Development mode")
	flag.BoolVar(&dumpDB, "dump", false, "Dump database")

	// Server address
	flag.StringVar(&host, "host", "localhost", "Server host")
	flag.StringVar(&port, "port", "3000", "Server port")
	flag.BoolVar(&secure, "secure", false, "Secure connection")

	// Database and Storage
	flag.StringVar(&dbPath, "db", "nasfiles.db", "Database path")
	flag.StringVar(&storageRoot, "storage", "storage", "Storage base folder")

	flag.Parse()
}

func main() {
	// Run the api server with the maximum number of cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse flags
	parseFlags()

	var err error
	if len(cfgPath) > 0 {
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

			StorageRoot: c.StorageRoot,
			Services: &nasfilesapi.Services{
				User: &bolt.UserService{
					DB: db,
				},
			},
		}
	} else {
		portInt, _ := strconv.Atoi(port)

		// Open database file
		db, err = boltdb.Open(dbPath, 0600, nil)
		if err != nil {
			color.HiRed("Couldn't open database file.")
			log.Fatalf("Couldn't open database file: %v", err)
		}

		// Initialize buckets if the database file was just created
		bolt.Setup(db)

		cfg = &nasfilesapi.Config{
			Development: development,
			Host:        host,
			Port:        portInt,
			Secure:      secure,

			StorageRoot: storageRoot,
			Services: &nasfilesapi.Services{
				User: &bolt.UserService{
					DB: db,
				},
			},
		}
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
