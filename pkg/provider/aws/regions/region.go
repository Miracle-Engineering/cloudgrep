package regions

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

// Region holds details on a region as used by the aws.Provider.
type Region struct {
	region *endpoints.Region
}

var _ fmt.Stringer = Region{}

func (r Region) String() string {
	return r.ID()
}

// ID returns the identifier for the region, or "global" for the global region.
func (r Region) ID() string {
	if r.IsGlobal() {
		return "global"
	}

	return r.region.ID()
}

// IsGlobal returns true if this Region refers to the global region, false otherwise.
func (r Region) IsGlobal() bool {
	return r.region == nil
}

// IsServiceSupported returns true if the service, identified by its endpoint, is supported in this region.
// Assumes all services are supported in the global region
func (r Region) IsServiceSupported(serviceEndpointID string) bool {
	if r.IsGlobal() {
		return true
	}

	services := r.base().Services()
	_, has := services[serviceEndpointID]

	return has
}

func (r Region) base() endpoints.Region {
	if r.IsGlobal() {
		return globalEndpointRegion()
	}

	return *r.region
}

func regionsFromStrings(rawRegions []string) ([]Region, error) {
	regions := make([]Region, 0, len(rawRegions))
	for _, raw := range rawRegions {
		region, err := regionForRaw(raw)
		if err != nil {
			return nil, err
		}

		regions = append(regions, region)
	}

	return regions, nil
}

func globalEndpointRegion() endpoints.Region {
	return officialRegions["us-east-1"]
}

func regionForRaw(raw string) (Region, error) {
	if raw == Global {
		return Region{}, nil
	}

	baseRegion, has := officialRegions[raw]
	if !has {
		return Region{}, fmt.Errorf("invalid region: %s", raw)
	}

	return Region{region: &baseRegion}, nil
}
