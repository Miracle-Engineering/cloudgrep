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

//go:embed resource_events.json
var embedResourceEvents []byte

//go:embed engine_status.json
var embedEngineStatus []byte

type TestEngineStatusResourceEvents struct {
	ResourceEvents []model.ResourceEvent `json:"resourceEvents"`
	EngineStatus   model.EngineStatus    `json:"engineStatus"`
}

func GetResources(t *testing.T) []*model.Resource {
	resources := []*model.Resource{}
	err := json.Unmarshal(embedResources, &resources)
	assert.NoError(t, err)
	assert.Positive(t, len(resources))
	return resources
}

func GetResourceEvents(t *testing.T) []model.ResourceEvent {
	resourceEvents := model.ResourceEvents{}
	err := json.Unmarshal(embedResourceEvents, &resourceEvents)
	assert.NoError(t, err)
	return resourceEvents
}

func GetEngineStatusesResourceEvents(t *testing.T) []TestEngineStatusResourceEvents {
	var engineStatusesResourceEvents []TestEngineStatusResourceEvents
	err := json.Unmarshal(embedEngineStatus, &engineStatusesResourceEvents)
	assert.NoError(t, err)
	return engineStatusesResourceEvents
}
