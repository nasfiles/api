package webdav

import (
	"github.com/nasfiles/api"

	"golang.org/x/net/webdav"
)

//UserHandler spawns a WebDAV HTTP handler for his folder
func UserHandler(u *api.User) webdav.Handler {
	var h webdav.Handler
	h = webdav.Handler{}

	return h
}
