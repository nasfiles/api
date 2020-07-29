package webdav

import (
	"golang.org/x/net/webdav"

	"github.com/nasfiles/nasfilesapi"
)

//UserHandler spawns a WebDAV HTTP handler for his folder
func UserHandler(u *nasfilesapi.User) webdav.Handler {
	var h webdav.Handler
	h = webdav.Handler{}

	return h
}
