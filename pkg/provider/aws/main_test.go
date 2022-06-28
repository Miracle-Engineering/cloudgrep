package aws

import (
	"encoding/json"
	"os"
	"path"
	"sort"
	"testing"

	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
	"golang.org/x/exp/maps"
)

func TestMain(m *testing.M) {
	code := m.Run()
	writeStats()
	os.Exit(code)
}

//generate the list of supported resources
func writeStats() {
	stats := testprovider.FetchStats()
	testedResources := maps.Keys(stats)

	if len(testedResources) == 0 {
		//the integration tests didn't run - do not generate the file
		return
	}

	sort.Strings(testedResources)

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	statsPath := path.Join(cwd, "zz_integration_stats.json")

	f, err := os.Create(statsPath)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(testedResources)
	if err != nil {
		panic(err)
	}
}
