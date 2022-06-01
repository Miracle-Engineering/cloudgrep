package testdata

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
)

//go:embed resources.json
var embedResources []byte

//go:embed engine_status.json
var embedEngineStatus []byte

func GetResources(t *testing.T) []*model.Resource {
	resources := []*model.Resource{}
	err := json.Unmarshal(embedResources, &resources)
	assert.NoError(t, err)
	assert.Positive(t, len(resources))
	return resources
}

func GetEngineStatus(t *testing.T) []model.EngineStatus {
	engineStatus := []model.EngineStatus{}
	err := json.Unmarshal(embedEngineStatus, &engineStatus)
	assert.NoError(t, err)
	return engineStatus
}
