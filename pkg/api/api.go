package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/command"
	"github.com/run-x/cloudgrep/static"
)

// GetHome renderes the home page
func GetHome(prefix string) http.Handler {
	if prefix != "" {
		prefix = "/" + prefix
	}
	return http.StripPrefix(prefix, http.FileServer(http.FS(static.Static)))
}

func GetAssets(prefix string) http.Handler {
	if prefix != "" {
		prefix = "/" + prefix + "static/"
	} else {
		prefix = "/static/"
	}
	return http.StripPrefix(prefix, http.FileServer(http.FS(static.Static)))
}

// GetInfo renders the system information
func GetInfo(c *gin.Context) {
	successResponse(c, gin.H{
		"version":    command.Version,
		"go_version": command.GoVersion,
		"git_sha":    command.GitCommit,
		"build_time": command.BuildTime,
	})
}
