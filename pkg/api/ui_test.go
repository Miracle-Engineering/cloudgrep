package api

import (
	"encoding/json"
	"github.com/run-x/cloudgrep/pkg/version"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeRoute(t *testing.T) {
	_, _, router := prepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, w.Body.Len() > 0)
}

func TestInfoRoute(t *testing.T) {
	_, _, router := prepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/info", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, body["version"] == version.Version)
	assert.True(t, body["go_version"] == version.GoVersion)
	assert.True(t, body["git_sha"] == version.GitCommit)
	assert.True(t, body["build_time"] == version.BuildTime)
}
