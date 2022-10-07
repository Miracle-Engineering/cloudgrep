package demo

import (
	"context"
	"os"
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/datastore"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestDemoData(t *testing.T) {
	//load the demo data
	cfg, err := GetDemoConfig()
	require.NoError(t, err)
	defer os.Remove(cfg.Datastore.DataSourceName)

	//check a temporary file is used
	require.Contains(t, cfg.Datastore.DataSourceName, "cloudgrepdemodb")

	//demo doesn't fetch any data
	require.True(t, cfg.Datastore.SkipRefresh)

	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	ds, err := datastore.NewDatastore(ctx, cfg, logger)
	require.NoError(t, err)

	//test some demo data
	result, err := ds.GetResources(ctx, nil)
	require.NoError(t, err)

	//core attribute are present
	field := result.FieldGroups.FindField("core", "region")
	require.Equal(t, "region", field.Name)

	//the teams attribute have some of the expected values used in the demo
	field = result.FieldGroups.FindField("tags", "team")
	for _, val := range []string{"billing", "consumer", "marketplace"} {
		require.Equal(t, val, field.Values.Find(val).Value)
	}

	field = result.FieldGroups.FindField("tags", "managed-by")
	for _, val := range []string{"cloudformation", "terraform", "(missing)"} {
		require.Equal(t, val, field.Values.Find(val).Value)
	}

	//test a filter
	query := `{
  "filter":{
    "tags.team": "consumer"
  }
}`
	result, err = ds.GetResources(ctx, []byte(query))
	require.NoError(t, err)
	require.Greater(t, result.Count, 10)
}
