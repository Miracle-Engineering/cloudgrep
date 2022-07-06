package api

import (
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
)

func SetupRoutes(router *gin.Engine, cfg config.Config, logger *zap.Logger, ds datastore.Datastore, engineF EngineFunc) {
	root := router.Group(cfg.Web.Prefix)

	root.GET("/", gin.WrapH(GetHome(cfg.Web.Prefix)))
	root.GET("/static/*path", gin.WrapH(GetAssets(cfg.Web.Prefix)))

	api := root.Group("/api")
	setupMiddlewares(api, cfg, logger, ds, engineF)

	healthz := root.Group("/healthz")
	healthz.Use(setSharedObjects(map[string]interface{}{
		"datastore": ds,
	}))
	healthz.GET("", Healthz)

	api.GET("/info", Info)
	api.GET("/resource", Resource)
	api.GET("/resources", Resources)
	api.POST("/resources", Resources)
	api.GET("/stats", Stats)
	api.GET("/enginestatus", EngineStatus)
	api.POST("/refresh", Refresh)
}
