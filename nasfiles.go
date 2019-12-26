package api

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

//Config ...
type Config struct {
	Development bool

	Host string
	Port int

	Secure bool

	StorageRoot string

	Services *Services
}

//Log prints all the coniguration values the API is running under
func (c *Config) Log() {
	// get terminal size
	var width int = 120

	x, err := terminal.Width()
	if err != nil {
		log.Println("Couldn't get the terminal width")
	}
	width = int(x)

	// Start printing configuration values
	color.HiYellow("Configuration")

	// Beginning separator
	LineSeparator("-", color.New(color.FgHiCyan), width)

	// Development mode
	fmt.Printf("Development mode: ")
	YesNoColored(c.Development)

	// Host
	fmt.Printf("Host: ")
	color.HiRed("%s:%d", c.Host, c.Port)

	// Secure
	fmt.Printf("Secure: ")
	YesNoColored(c.Secure)

	// Storage
	fmt.Printf("Storage path: ")
	color.HiBlue(c.StorageRoot)

	// Ending separator
	LineSeparator("-", color.New(color.FgHiCyan), width)
}

//Services is a struct to tie all the services into a unified struct
type Services struct {
	User UserService
}
