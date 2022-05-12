package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/static"
)

func StartWebServer(ctx context.Context, cfg config.Config, ds datastore.Datastore) {
	router := gin.Default()

	if cfg.Logging.IsDev() {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	SetupRoutes(router, cfg, ds)

	fmt.Println("Starting server...")
	go func() {
		err := router.Run(fmt.Sprintf("%v:%v", cfg.Web.Host, cfg.Web.Port))
		if err != nil {
			fmt.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
}

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
		"version":    Version,
		"go_version": GoVersion,
		"git_sha":    GitCommit,
		"build_time": BuildTime,
	})
}

// GetResources retrieves the cloud resources matching the query parameters
func GetResources(c *gin.Context) {
	datastore := c.MustGet("datastore").(datastore.Datastore)
	filter := model.EmptyFilter()
	if f, ok := c.Get("filter"); ok {
		filter = f.(model.Filter)
	}
	resources, err := datastore.GetResources(c, filter)
	if err != nil {
		badRequest(c, err)
		return
	}
	c.JSON(200, resources)
}

// Stats provides stats about stored data
func Stats(c *gin.Context) {
	datastore := c.MustGet("datastore").(datastore.Datastore)
	stats, err := datastore.Stats(c)
	if err != nil {
		badRequest(c, err)
		return
	}
	c.JSON(200, stats)
}
