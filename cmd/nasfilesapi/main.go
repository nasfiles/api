package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	boltdb "github.com/boltdb/bolt"
	"github.com/fatih/color"
	"golang.org/x/net/webdav"

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

func loadUsers(c *nasfilesapi.Config) error {
	// Fetch all users
	users, err := c.Services.User.GetsAll()
	if err != nil {
		return err
	}

	// Create WebDAV handler and save it in the memory
	for _, user := range *users {
		user.Handler = &webdav.Handler{
			FileSystem: webdav.Dir(path.Join(c.StorageRoot, user.Username)),
			LockSystem: webdav.NewMemLS(),
			Logger: func(r *http.Request, err error) {
				if err != nil {
					log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
				} else {
					log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
				}
			},
		}

		c.Users[user.Username] = user
	}

	return nil
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

		Auth:       c.Auth,
		PrivateKey: c.PrivateKey,

		StorageRoot: c.StorageRoot,
		Users:       map[string]nasfilesapi.User{},
		DefaultHandler: &webdav.Handler{
			FileSystem: webdav.Dir(c.StorageRoot),
			LockSystem: webdav.NewMemLS(),
			Logger: func(r *http.Request, err error) {
				if err != nil {
					log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
				} else {
					log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
				}
			},
		},

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

	// Load users into the memory
	loadUsers(cfg)

	// Start API
	h.Serve(cfg)
}
