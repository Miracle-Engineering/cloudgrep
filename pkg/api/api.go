package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/version"
)

type EngineFunc func(context.Context) error

func StartWebServer(ctx context.Context, cfg config.Config, logger *zap.Logger, ds datastore.Datastore, engineF EngineFunc) {
	router := gin.Default()

	if logger.Core().Enabled(zap.DebugLevel) {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	SetupRoutes(router, cfg, logger, ds, engineF)

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
		"version":   version.Version,
		"goVersion": version.GoVersion,
		"gitSha":    version.GitCommit,
		"buildTime": version.BuildTime,
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
	var body []byte
	//the body contains the query
	if c.Request.Body != nil {
		var err error
		body, err = c.GetRawData()
		if err != nil {
			badRequest(c, err)
			return
		}
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

//EngineStatus returns the status of the engine
func EngineStatus(c *gin.Context) {
	ds := c.MustGet("datastore").(datastore.Datastore)
	status, err := ds.EngineStatus(c)
	if err != nil {
		badRequest(c, err)
		return
	}
	c.JSON(200, status)
}

//Refresh trigger the engine to fetch the resources
func Refresh(c *gin.Context) {

	ds := c.MustGet("datastore").(datastore.Datastore)
	logger := c.MustGet("logger").(*zap.Logger)
	status, err := ds.GetEngineStatus(c)
	if err != nil {
		badRequest(c, err)
		return
	}
	if status.Status == model.EngineStatusFetching {
		errorResponse(c, http.StatusAccepted, fmt.Errorf("engine is already running"))
		return
	}
	engineFunc := c.MustGet("engineFunc").(EngineFunc)
	//run this process async - the UI will call EngineStatus api to check the result
	go func() {
		err := engineFunc(c)
		logger.Sugar().Errorw("some error(s) when running the provider engine", "error", err)
	}()
	c.Status(http.StatusOK)

}
