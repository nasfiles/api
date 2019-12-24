package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/nasfiles/api"
)

func main() {
	// parse flags
	cfgPath := flag.String("config", "", "JSON config file")

	development := flag.Bool("development", false, "Development mode")
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "3000", "Server port")
	secure := flag.Bool("secure", false, "Secure connection")
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
		}
	}

	fmt.Println(cfg.Host)
	cfg.Stats()

}
