package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/datastore/testdata"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func prepareApiUnitTest(t *testing.T) (*zap.Logger, datastore.Datastore, *gin.Engine) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	datastoreConfigs := config.Datastore{
		Type:           "sqlite",
		DataSourceName: "file::memory:",
	}
	cfg, err := config.GetDefault()
	require.NoError(t, err)
	cfg.Datastore = datastoreConfigs

	ds, err := datastore.NewDatastore(ctx, cfg, zaptest.NewLogger(t))
	require.NoError(t, err)

	router := gin.Default()
	SetupRoutes(router, cfg, logger, ds)

	//write the resources
	resources := testdata.GetResources(t)
	require.NotZero(t, len(resources))
	require.NoError(t, ds.WriteResources(ctx, resources))
	return logger, ds, router
}

func TestStatsRoute(t *testing.T) {
	_, _, router := prepareApiUnitTest(t)
	path := "/api/stats"

	t.Run("SomeResources", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)
		router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var body *model.Stats
		err := json.Unmarshal(w.Body.Bytes(), &body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, body.ResourcesCount, 3)
	})
}

func TestResourcesRoute(t *testing.T) {
	_, _, router := prepareApiUnitTest(t)
	path := "/api/resources"

	t.Run("SomeResources", func(t *testing.T) {
		resources := testdata.GetResources(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var body model.Resources
		err = json.Unmarshal(w.Body.Bytes(), &body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, len(body), 3)
		model.AssertEqualsResources(t, body, resources)
	})
}

func TestResourceRoute(t *testing.T) {
	_, _, router := prepareApiUnitTest(t)
	path := "/api/resource"

	t.Run("MissingParam", func(t *testing.T) {
		var body map[string]interface{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Equal(t, body["status"], float64(http.StatusBadRequest))
		require.Equal(t, body["error"], "missing required parameter 'id'")
	})

	t.Run("EmptyParam", func(t *testing.T) {
		var body map[string]interface{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("id", "")
		req.URL.RawQuery = q.Encode()
		router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Equal(t, body["status"], float64(http.StatusBadRequest))
		require.Equal(t, body["error"], "missing required parameter 'id'")
	})

	t.Run("UnknownParam", func(t *testing.T) {
		var body map[string]interface{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("id", "blah")
		req.URL.RawQuery = q.Encode()
		router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusNotFound, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Equal(t, body["status"], float64(http.StatusNotFound))
		require.Equal(t, body["error"], "can't find resource with id 'blah'")
	})

	t.Run("ValidParam", func(t *testing.T) {
		var actualResource model.Resource
		resources := testdata.GetResources(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("id", resources[0].Id)
		req.URL.RawQuery = q.Encode()
		router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &actualResource))
		model.AssertEqualsResource(t, actualResource, *resources[0])
	})
}
