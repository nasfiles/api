package nasfilesapi

import (
	"fmt"

	"github.com/fatih/color"
	"golang.org/x/net/webdav"

	"github.com/nasfiles/nasfilesapi/terminal"
)

// Services is a struct to tie all the services into a unified struct
type Services struct {
	User UserService
	Auth AuthService
}

// Config is a struct which contains all working components of the NAS File API
type Config struct {
	Development bool
	Host        string
	Port        int
	Secure      bool

	Auth       bool
	PrivateKey string

	StorageRoot string
	Users       map[string]User

	// API/Database service
	Services *Services

	// WebDAV default service
	DefaultHandler *webdav.Handler
}

// Log prints all the coniguration values the API is running under
func (c *Config) Log() {
	width := terminal.Size()
	if width == 0 {
		width = 60
	}

	// Start printing configuration values
	color.HiYellow("Configuration")

	// Beginning separator
	terminal.LineSeparator("-", color.New(color.FgHiCyan), width)

	// Development mode
	fmt.Printf("Development mode: ")
	terminal.YesNoColored(c.Development)

	// Host
	fmt.Printf("Host: ")
	color.HiRed("%s:%d", c.Host, c.Port)

	// Secure
	fmt.Printf("Secure: ")
	terminal.YesNoColored(c.Secure)

	// Storage
	fmt.Printf("Storage path: ")
	color.HiBlue(c.StorageRoot)

	// Ending separator
	terminal.LineSeparator("-", color.New(color.FgHiCyan), width)
}
