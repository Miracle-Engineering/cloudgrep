package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
)

func SetupMiddlewares(group *gin.RouterGroup, cfg config.Config, ds datastore.Datastore) {
	if cfg.Logging.IsDev() {
		group.Use(logAllQueryParams(cfg), logAllRequests(cfg))
	}
	group.Use(setDatastore(ds))
}

func SetupRoutes(router *gin.Engine, cfg config.Config, ds datastore.Datastore) {
	root := router.Group(cfg.Web.Prefix)

	root.GET("/", gin.WrapH(GetHome(cfg.Web.Prefix)))
	root.GET("/static/*path", gin.WrapH(GetAssets(cfg.Web.Prefix)))

	api := root.Group("/api")
	SetupMiddlewares(api, cfg, ds)

	api.GET("/info", GetInfo)
	api.GET("/resources", GetResources)

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
