package api

import (
	"github.com/run-x/cloudgrep/static"
	"net/http"
)

// GetHome renders the home page
func GetHome(prefix string) http.Handler {
	if prefix != "" && prefix != "/" {
		prefix = "/" + prefix
	}
	return http.StripPrefix(prefix, http.FileServer(http.FS(static.Static)))
}

func GetAssets(prefix string) http.Handler {
	if prefix != "" && prefix != "/" {
		prefix = "/" + prefix + "static/"
	} else {
		prefix = "/static/"
	}
	return http.StripPrefix(prefix, http.FileServer(http.FS(static.Static)))
}
