package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Development bool

	Host string
	Port int

	Secure bool

	StorageRoot string
}

func loadConfig(path string) (*config, error) {
	c := &config{}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	parser := json.NewDecoder(f)
	err = parser.Decode(&c)

	return c, nil
}
