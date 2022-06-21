package regions

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/require"
)

func mustRegion(t testing.TB, raw string) Region {
	region, err := regionForRaw(raw)
	require.NoError(t, err)
	return region
}

func regionIds(regions []Region) []string {
	return testingutil.SliceConvertFunc(regions, func(region Region) string {
		return region.ID()
	})
}
