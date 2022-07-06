package testingutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/run-x/cloudgrep/pkg/model"
)

type ResourceFilter struct {
	Region  string
	Type    string
	Tags    model.Tags
	RawData map[string]any
}

func (f ResourceFilter) String() string {
	var parts []string

	if f.Type != "" {
		parts = append(parts, fmt.Sprintf("Type=%s", f.Type))
	}

	if f.Region != "" {
		parts = append(parts, fmt.Sprintf("Region=%s", f.Region))
	}

	if f.Tags != nil && len(f.Tags) == 0 {
		parts = append(parts, "Tags=[]")
	} else {
		for _, tag := range f.Tags {
			if tag.Value == "" {
				parts = append(parts, fmt.Sprintf("Tags[%s]", tag.Key))
			} else {
				parts = append(parts, fmt.Sprintf("Tags[%s]=%s", tag.Key, tag.Value))
			}
		}
	}

	if len(f.RawData) > 0 {
		rawParts := make([]string, 0, len(f.RawData))

		for key, val := range f.RawData {
			rawParts = append(rawParts, fmt.Sprintf("%s=%v", key, val))
		}
		//sorting ensures consistent output for testing
		sort.Strings(rawParts)
		parts = append(parts, fmt.Sprintf("RawData={%s}", strings.Join(rawParts, ", ")))
	}

	fields := strings.Join(parts, ", ")
	return fmt.Sprintf("ResourceFilter{%s}", fields)
}

func (f ResourceFilter) Matches(resource model.Resource) bool {
	tests := resourceFilterMatchFuncs()

	for _, test := range tests {
		match, tested := test(f, resource)
		if tested && !match {
			return false
		}
	}

	return true
}

func (f ResourceFilter) Filter(in []model.Resource) []model.Resource {
	out := make([]model.Resource, 0, len(in))

	for _, resource := range in {
		if f.Matches(resource) {
			out = append(out, resource)
		}
	}

	return out
}

// PartialFilter makes it easy to debug where the filter isn't matching resources
func (f ResourceFilter) PartialFilter(in []model.Resource) map[string][]model.Resource {
	output := make(map[string][]model.Resource)

	tests := resourceFilterMatchFuncs()
	for name, test := range tests {
		active := false
		var resources []model.Resource
		for _, res := range in {
			matched, tested := test(f, res)
			if !tested {
				continue
			}

			active = true
			if matched {
				resources = append(resources, res)
			}
		}

		if active {
			output[name] = resources
		}
	}

	return output
}

type resourceFilterTestFunc = func(ResourceFilter, model.Resource) (bool, bool)

func resourceFilterMatchFuncs() map[string]resourceFilterTestFunc {
	return map[string]resourceFilterTestFunc{
		"region":  resourceFilterMatchRegion,
		"type":    resourceFilterMatchType,
		"tags":    resourceFilterMatchTags,
		"rawData": resourceFilterMatchRawData,
	}
}

func resourceFilterMatchRegion(f ResourceFilter, r model.Resource) (bool, bool) {
	if f.Region == "" {
		return false, false
	}

	return f.Region == r.Region, true
}

func resourceFilterMatchType(f ResourceFilter, r model.Resource) (bool, bool) {
	if f.Type == "" {
		return false, false
	}

	return f.Type == r.Type, true
}

func resourceFilterMatchTags(f ResourceFilter, r model.Resource) (bool, bool) {
	if f.Tags == nil {
		return false, false
	}

	// Treat empty slice as special "no tags" filter
	if len(f.Tags) == 0 {
		return len(r.Tags) == 0, true
	}

	tagMap := make(map[string]string)
	for _, tag := range r.Tags {
		tagMap[tag.Key] = tag.Value
	}

	for _, tag := range f.Tags {
		val, has := tagMap[tag.Key]
		if !has {
			return false, true
		}

		if tag.Value == "" {
			continue
		}

		if strings.TrimSpace(val) != strings.TrimSpace(tag.Value) {
			return false, true
		}
	}

	return true, true
}

func resourceFilterMatchRawData(f ResourceFilter, r model.Resource) (bool, bool) {
	if len(f.RawData) == 0 {
		return false, false
	}

	var raw map[string]any
	err := json.Unmarshal(r.RawData, &raw)
	if err != nil {
		panic(fmt.Errorf("cannot pase model.Resource.RawData: %s", r.Id))
	}

	for key, val := range f.RawData {
		rawVal, has := raw[key]
		if !has {
			return false, true
		}

		if !reflect.DeepEqual(val, rawVal) {
			return false, true
		}
	}

	return true, true
}
