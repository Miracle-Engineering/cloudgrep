package api

import (
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
)

func SetupRoutes(router *gin.Engine, cfg config.Config, logger *zap.Logger, ds datastore.Datastore) {
	root := router.Group(cfg.Web.Prefix)

	root.GET("/", gin.WrapH(GetHome(cfg.Web.Prefix)))
	root.GET("/static/*path", gin.WrapH(GetAssets(cfg.Web.Prefix)))

	api := root.Group("/api")
	setupMiddlewares(api, cfg, logger, ds)

	api.GET("/info", Info)
	api.GET("/resource", Resource)
	api.GET("/resources", Resources)
	api.GET("/stats", Stats)
	api.GET("/fields", Fields)

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
