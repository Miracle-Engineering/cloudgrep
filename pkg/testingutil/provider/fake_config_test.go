package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeProviderResourceConfig_RunCount(t *testing.T) {
	config := FakeProviderResourceConfig{}

	assert.Zero(t, config.RunCount())
	config.addRun()
	assert.Equal(t, 1, config.RunCount())
}
