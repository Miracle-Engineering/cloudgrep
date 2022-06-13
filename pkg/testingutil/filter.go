package testingutil

import (
	"encoding/json"
	"fmt"
	"reflect"
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

	if f.Region != "" {
		parts = append(parts, fmt.Sprintf("Region=%s", f.Region))
	}

	if f.Type != "" {
		parts = append(parts, fmt.Sprintf("Type=%s", f.Type))
	}

	for _, tag := range f.Tags {
		parts = append(parts, fmt.Sprintf("Tags[%s]=%s", tag.Key, tag.Value))
	}

	if len(f.RawData) > 0 {
		rawParts := make([]string, 0, len(f.RawData))

		for key, val := range f.RawData {
			rawParts = append(rawParts, fmt.Sprintf("%s=%v", key, val))
		}

		parts = append(parts, fmt.Sprintf("RawData={%s}", strings.Join(rawParts, ", ")))
	}

	fields := strings.Join(parts, ", ")
	return fmt.Sprintf("ResourceFilter{%s}", fields)
}

func (f ResourceFilter) Matches(resource model.Resource) bool {
	if f.Region != "" {
		if resource.Region != f.Region {
			return false
		}
	}

	if f.Type != "" {
		if resource.Type != f.Type {
			return false
		}
	}

	// Treat empty slice different from nil
	if f.Tags != nil {
		// Treat empty slice as special "no tags" filter
		if len(f.Tags) == 0 {
			if len(resource.Tags) != 0 {
				return false
			}
		} else {
			tagMap := make(map[string]string)
			for _, tag := range resource.Tags {
				tagMap[tag.Key] = tag.Value
			}

			for _, tag := range f.Tags {
				val, has := tagMap[tag.Key]
				if !has {
					return false
				}

				if strings.TrimSpace(val) != strings.TrimSpace(tag.Value) {
					return false
				}
			}
		}
	}

	if len(f.RawData) > 0 {
		var raw map[string]any
		err := json.Unmarshal(resource.RawData, &raw)
		if err != nil {
			panic(fmt.Errorf("cannot pase model.Resource.RawData: %s", resource.Id))
		}

		for key, val := range f.RawData {
			rawVal, has := raw[key]
			if !has {
				return false
			}

			if !reflect.DeepEqual(val, rawVal) {
				return false
			}
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
