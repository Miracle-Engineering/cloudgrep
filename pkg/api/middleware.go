package api

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
)

// Middleware to print out request parameters and body for debugging
func logAllQueryParams(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Request.ParseForm()
		cfg.Logging.Logger.Sugar().Infow("Request params:", "error", err, "params", c.Request.Form)
	}
}

func logAllRequests(cfg config.Config) gin.HandlerFunc {
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	return ginzap.Ginzap(cfg.Logging.Logger, time.RFC3339, true)
}
