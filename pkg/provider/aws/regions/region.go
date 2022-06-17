package regions

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

type Region struct {
	region *endpoints.Region
}

var _ fmt.Stringer = Region{}

func (r Region) String() string {
	return r.ID()
}

func (r Region) ID() string {
	if r.region == nil {
		return "global"
	}

	return r.region.ID()
}

func (r Region) IsGlobal() bool {
	return r.region == nil
}

func (r Region) IsServiceSupported(name string) bool {
	if r.IsGlobal() {
		return true
	}

	services := r.base().Services()
	_, has := services[name]

	return has
}

func (r Region) base() endpoints.Region {
	if r.region == nil {
		return globalEndpointRegion()
	}

	return *r.region
}

func regionsFromStrings(rawRegions []string) []Region {
	regions := make([]Region, 0, len(rawRegions))
	for _, raw := range rawRegions {
		if raw == Global {
			regions = append(regions, Region{})
			continue
		}

		eregion, has := officialRegions[raw]
		if !has {
			panic(fmt.Errorf("invalid region: %s", raw))
		}

		regions = append(regions, Region{region: &eregion})
	}

	return regions
}

func globalEndpointRegion() endpoints.Region {
	return officialRegions["us-east-1"]
}
