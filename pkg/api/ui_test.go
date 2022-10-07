package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestHomeRoute(t *testing.T) {
	m := prepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	m.router.ServeHTTP(w, req)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, w.Body.Len() > 0)
}

func TestInfoRoute(t *testing.T) {
	m := prepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/info", nil)
	m.router.ServeHTTP(w, req)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, body["version"] == version.Version)
	assert.True(t, body["goVersion"] == version.GoVersion)
	assert.True(t, body["gitSha"] == version.GitCommit)
	assert.True(t, body["buildTime"] == version.BuildTime)
}
