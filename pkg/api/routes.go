package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/command"
)

func SetupMiddlewares(group *gin.RouterGroup) {
	if command.Opts.Debug {
		group.Use(requestInspectMiddleware())
	}

}

func SetupRoutes(router *gin.Engine) {
	root := router.Group(command.Opts.Prefix)

	root.GET("/", gin.WrapH(GetHome(command.Opts.Prefix)))
	root.GET("/static/*path", gin.WrapH(GetAssets(command.Opts.Prefix)))

	api := root.Group("/api")
	SetupMiddlewares(api)

	api.GET("/info", GetInfo)

	// mock api serving static files (temporary)
	mock_files, err := ioutil.ReadDir("./static/mock")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range mock_files {
		basename := f.Name()
		name := strings.TrimSuffix(basename, filepath.Ext(basename))
		relativePath := fmt.Sprintf("/mock/api/%s", name)
		filePath := fmt.Sprintf("./static/mock/%s", basename)
		log.Printf("[MOCK] %s -> %s", relativePath, filePath)
		router.StaticFile(relativePath, filePath)
	}
}
