package regions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegion_String_usw1(t *testing.T) {
	region := mustRegion(t, "us-west-1")

	assert.Equal(t, "us-west-1", region.String())
}

func TestRegion_String_global(t *testing.T) {
	region := mustRegion(t, "global")

	assert.Equal(t, "global", region.String())
}

func TestRegion_IsServiceSupported_global(t *testing.T) {
	region := mustRegion(t, "global")

	assert.True(t, region.IsServiceSupported("foobar"))
}

const simpleDbEndpointID = "sdb"

func TestRegion_IsServiceSupported_yes(t *testing.T) {
	region := mustRegion(t, "us-east-1")

	assert.True(t, region.IsServiceSupported(simpleDbEndpointID))
}

func TestRegion_IsServiceSupported_no(t *testing.T) {
	region := mustRegion(t, "us-east-2")

	assert.False(t, region.IsServiceSupported(simpleDbEndpointID))
}

func TestRegionsFromStrings_empty(t *testing.T) {
	regions, err := regionsFromStrings(nil)
	assert.NoError(t, err)
	assert.Empty(t, regions)
}

func TestRegionsFromStrings_invalid(t *testing.T) {
	regions, err := regionsFromStrings([]string{"foo"})
	assert.ErrorContains(t, err, "invalid region: foo")
	assert.Empty(t, regions)
}

func TestRegionsFromStrings_valid(t *testing.T) {
	regions, err := regionsFromStrings([]string{"us-east-1", "us-west-2", "global"})
	assert.NoError(t, err)
	assert.Len(t, regions, 3)

	assert.Equal(t, "us-east-1", regions[0].ID())
	assert.Equal(t, "us-west-2", regions[1].ID())
	assert.Equal(t, "global", regions[2].ID())
}
