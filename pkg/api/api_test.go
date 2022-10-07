package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juandiegopalomino/cloudgrep/pkg/config"
	"github.com/juandiegopalomino/cloudgrep/pkg/datastore"
	"github.com/juandiegopalomino/cloudgrep/pkg/datastore/testdata"
	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func prepareApiUnitTest(t *testing.T) *mockApi {
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
	resources := testdata.GetResources(t)
	mockApi := mockApi{
		ctx:       ctx,
		router:    router,
		ds:        ds,
		resources: resources,
	}
	SetupRoutes(router, cfg, logger, ds, mockApi.runEngine)

	//write the resources
	require.NotZero(t, len(resources))
	require.NoError(t, mockApi.runEngine(ctx))
	return &mockApi
}

type mockApi struct {
	ctx       context.Context
	router    *gin.Engine
	ds        datastore.Datastore
	resources model.Resources
	//if set calling the engine will return it
	engineErr error
}

// runEngine simulates running the engine by writing resources to the datastore
func (m *mockApi) runEngine(ctx context.Context) error {
	err := m.ds.WriteEvent(ctx, model.NewEngineEventStart())
	if err != nil {
		return err
	}
	if m.engineErr == nil {
		err = m.ds.WriteResources(ctx, m.resources)
	} else {
		err = m.engineErr
	}
	//for testing async simulate a longer run by waiting
	time.Sleep(10 * time.Millisecond)

	defer func() {
		err := m.ds.WriteEvent(ctx, model.NewEngineEventEnd(err))
		if err != nil {
			log.Default().Println(err.Error())
		}
	}()

	if err != nil {
		return err
	}

	return nil
}

func (m *mockApi) waitForEngine(status string) {
	for {
		statusResp, _ := m.ds.EngineStatus(m.ctx)
		if statusResp.Status == status {
			break
		}
		time.Sleep(time.Millisecond)
	}
}

func TestStatsRoute(t *testing.T) {

	m := prepareApiUnitTest(t)
	path := "/api/stats"

	t.Run("SomeResources", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var body *model.Stats
		err := json.Unmarshal(w.Body.Bytes(), &body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, body.ResourcesCount, 3)
	})
}

func TestResourcesRoute(t *testing.T) {
	m := prepareApiUnitTest(t)
	path := "/api/resources"

	t.Run("SomeResources", func(t *testing.T) {
		resources := testdata.GetResources(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var body model.ResourcesResponse
		err = json.Unmarshal(w.Body.Bytes(), &body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, len(body.Resources), 3)
		require.Equal(t, body.Count, 3)
		testingutil.AssertEqualsResources(t, body.Resources, resources)
	})
}

func TestResourcesPostRoute(t *testing.T) {
	m := prepareApiUnitTest(t)
	path := "/api/resources"

	all_resources := testdata.GetResources(t)
	resourceInst1 := all_resources[0]  //i-123 team:infra, release tag, tag region:us-west-2
	resourceInst2 := all_resources[1]  //i-124 team:dev, no release tag
	resourceBucket := all_resources[2] //s3 bucket without tags

	t.Run("FilterSearch", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := strings.NewReader(`{
  "filter":{
    "$or":[
      {
        "tags.team":"infra"
      },
      {
        "tags.team":"dev"
      }
    ]
  }
}`)
		req, err := http.NewRequest("POST", path, body)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var response model.ResourcesResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, 2, response.Count)
		testingutil.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, response.Resources)

		//check a few fields
		fields := response.FieldGroups
		testingutil.AssertEqualsField(t, model.Field{
			Name:  "team",
			Count: 2,
			Values: model.FieldValues{
				&model.FieldValue{Value: "infra", Count: "1"},
				&model.FieldValue{Value: "dev", Count: "1"},
			}}, *fields.FindField("tags", "team"))
		testingutil.AssertEqualsField(t, model.Field{
			Name:  "type",
			Count: 2,
			Values: model.FieldValues{
				&model.FieldValue{Value: "test.Instance", Count: "2"},
				&model.FieldValue{Value: "s3.Bucket", Count: "-"},
			}}, *fields.FindField("core", "type"))

	})

	t.Run("FilterEmpty", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := strings.NewReader(`{
  "filter":{ }
}`)
		req, err := http.NewRequest("POST", path, body)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var response model.ResourcesResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		testingutil.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2, resourceBucket}, response.Resources)
	})

	t.Run("NoBody", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := strings.NewReader(``)
		req, err := http.NewRequest("POST", path, body)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var response model.ResourcesResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		testingutil.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2, resourceBucket}, response.Resources)
	})
}

func TestResourceRoute(t *testing.T) {
	m := prepareApiUnitTest(t)
	path := "/api/resource"

	t.Run("MissingParam", func(t *testing.T) {
		var body map[string]interface{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
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
		m.router.ServeHTTP(w, req)
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
		m.router.ServeHTTP(w, req)
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
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &actualResource))
		testingutil.AssertEqualsResource(t, actualResource, *resources[0])
	})
}

func TestResourceFieldsRoute(t *testing.T) {
	m := prepareApiUnitTest(t)
	path := "/api/resources"

	t.Run("Standard", func(t *testing.T) {
		var response model.ResourcesResponse
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		fields := response.FieldGroups
		require.Equal(t, len(fields), 2)
		//check number of groups
		require.Equal(t, 2, len(fields))
		//check fields by group
		require.Equal(t, 3, len(fields.FindGroup("core").Fields))
		require.Equal(t, 10, len(fields.FindGroup("tags").Fields))
	})
}

func TestRefreshPostRoute(t *testing.T) {
	refreshPath := "/api/refresh"
	engineStatusPath := "/api/enginestatus"

	t.Run("Success", func(t *testing.T) {
		m := prepareApiUnitTest(t)
		//trigger a refresh
		record := httptest.NewRecorder()
		req, err := http.NewRequest("POST", refreshPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)
		require.Equal(t, "", record.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, record.Code)
		//wait for engine to complete
		m.waitForEngine("fetching")
		m.waitForEngine("success")

		//check that the status api returns success
		record = httptest.NewRecorder()
		req, err = http.NewRequest("GET", engineStatusPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)
		require.Equal(t, "application/json; charset=utf-8", record.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, record.Code)
		body := make(map[string]interface{})
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &body))
		require.Equal(t, "success", body["status"])
		require.Equal(t, "", body["error"])
	})

	t.Run("Engine Error", func(t *testing.T) {
		//test an error while refreshing
		m := prepareApiUnitTest(t)
		//configure an error in the engine
		m.engineErr = fmt.Errorf("There was an engine error")
		record := httptest.NewRecorder()
		req, err := http.NewRequest("POST", refreshPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)

		//check engine was run
		require.Equal(t, "", record.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, record.Code)

		//wait for engine to complete
		m.waitForEngine("fetching")
		m.waitForEngine("failed")

		//check that the status api returns the error
		record = httptest.NewRecorder()
		req, err = http.NewRequest("GET", engineStatusPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)
		require.Equal(t, "application/json; charset=utf-8", record.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, record.Code)
		body := make(map[string]interface{})
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &body))
		require.Equal(t, "failed", body["status"])
		require.Equal(t, "There was an engine error", body["error"])
	})

	t.Run("Engine Already Running", func(t *testing.T) {
		m := prepareApiUnitTest(t)
		record := httptest.NewRecorder()
		req, err := http.NewRequest("POST", refreshPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)
		require.Equal(t, http.StatusOK, record.Code)
		//the engine is still running at this point
		m.waitForEngine("fetching")

		//refresh again - it should respond "engine already running"
		record = httptest.NewRecorder()
		req, err = http.NewRequest("POST", refreshPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)
		require.Equal(t, http.StatusAccepted, record.Code)
		body := make(map[string]interface{})
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &body))
		require.Equal(t, float64(202), body["status"])
		require.Equal(t, "engine is already running", body["error"])

		//check that the status api returns "fetching"
		record = httptest.NewRecorder()
		req, err = http.NewRequest("GET", engineStatusPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)
		require.Equal(t, "application/json; charset=utf-8", record.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, record.Code)
		body = make(map[string]interface{})
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &body))
		require.Equal(t, "fetching", body["status"])
		require.Equal(t, "", body["error"])
	})
}
