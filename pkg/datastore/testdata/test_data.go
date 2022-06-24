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

func GetResources(t *testing.T) []*model.Resource {
	resources := []*model.Resource{}
	err := json.Unmarshal(embedResources, &resources)
	assert.NoError(t, err)
	assert.Positive(t, len(resources))
	return resources
}
