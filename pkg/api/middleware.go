package api

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/model"
	"go.uber.org/zap"
)

func setupMiddlewares(group *gin.RouterGroup, cfg config.Config, logger *zap.Logger, ds datastore.Datastore) {
	if logger.Core().Enabled(zap.DebugLevel) {
		group.Use(logAllQueryParams(cfg, logger), logAllRequests(cfg, logger))
	}
	group.Use(setDatastore(ds))
	group.Use(setParams(cfg, logger))
}

// Middleware to print out request parameters and body for debugging
func logAllQueryParams(cfg config.Config, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Request.ParseForm()
		logger.Sugar().Debug("Request params:", "error", err, "params", c.Request.Form)
	}
}

func logAllRequests(cfg config.Config, logger *zap.Logger) gin.HandlerFunc {
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	return ginzap.Ginzap(logger, time.RFC3339, true)
}

func setDatastore(ds datastore.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("datastore", ds)
		c.Next()
	}
}

func setParams(cfg config.Config, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter model.Filter
		tags := c.QueryMap("tags")
		if len(tags) > 0 {
			filter = model.NewFilter(tags)
			c.Set("filter", filter)
		}
		id := c.Query("id")
		c.Set("id", id)
		logger.Sugar().Debugw("Request params:",
			zap.Object("filter", filter),
			zap.String("id", id),
		)
	}
}
