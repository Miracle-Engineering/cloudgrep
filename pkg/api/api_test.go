package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func PrepareApiUnitTest(t *testing.T) (*zap.Logger, datastore.Datastore, *gin.Engine) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	datastoreConfigs := config.Datastore{
		Type:           "sqlite",
		DataSourceName: "file::memory:",
	}
	cfg, err := config.GetDefault()
	assert.NoError(t, err)
	cfg.Datastore = datastoreConfigs

	ds, err := datastore.NewDatastore(ctx, cfg, zaptest.NewLogger(t))
	assert.NoError(t, err)

	router := gin.Default()
	SetupRoutes(router, cfg, logger, ds)
	return logger, ds, router
}

func TestHealthRoute(t *testing.T) {
	_, _, router := PrepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"status\":\"All good!\"}", w.Body.String())
}
