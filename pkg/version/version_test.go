package version

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Kind of straight forward, just running this for now
	buildConfig := Get()
	assert.Equal(t, buildConfig.GoVersion, "testing")
}
