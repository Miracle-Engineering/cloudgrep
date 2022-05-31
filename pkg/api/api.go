package api

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/version"
)

func StartWebServer(ctx context.Context, cfg config.Config, logger *zap.Logger, ds datastore.Datastore) {
	router := gin.Default()

	if logger.Core().Enabled(zap.DebugLevel) {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	SetupRoutes(router, cfg, logger, ds)

	fmt.Println("Starting server...")
	go func() {
		err := router.Run(fmt.Sprintf("%v:%v", cfg.Web.Host, cfg.Web.Port))
		if err != nil {
			fmt.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
}

// Info renders the system information
func Info(c *gin.Context) {
	successResponse(c, gin.H{
		"version":    version.Version,
		"go_version": version.GoVersion,
		"git_sha":    version.GitCommit,
		"build_time": version.BuildTime,
	})
}

// Resource retrieves a resource by its id
func Resource(c *gin.Context) {
	ds := c.MustGet("datastore").(datastore.Datastore)
	id := c.GetString("id")
	if id == "" {
		badRequest(c, fmt.Errorf("missing required parameter 'id'"))
		return
	}
	resource, err := ds.GetResource(c, id)
	if err != nil {
		badRequest(c, err)
	}
	if resource == nil {
		notFoundf(c, "can't find resource with id '%v'", id)
		return
	}
	c.JSON(200, resource)
}

// Resources retrieves the cloud resources matching the query parameters
func Resources(c *gin.Context) {
	ds := c.MustGet("datastore").(datastore.Datastore)

	//the body contains the query
	body, err := c.GetRawData()
	if err != nil {
		badRequest(c, err)
		return
	}
	resources, err := ds.GetResources(c, body)
	if err != nil {
		badRequest(c, err)
		return
	}
	c.JSON(200, resources)
}

// Stats provides stats about stored data
func Stats(c *gin.Context) {
	ds := c.MustGet("datastore").(datastore.Datastore)
	stats, err := ds.Stats(c)
	if err != nil {
		badRequest(c, err)
		return
	}
	c.JSON(200, stats)
}

// Fields Return the list of fields available for filtering the resources
func Fields(c *gin.Context) {
	ds := c.MustGet("datastore").(datastore.Datastore)
	fields, err := ds.GetFields(c)
	if err != nil {
		badRequest(c, err)
		return
	}
	c.JSON(200, fields)
}
