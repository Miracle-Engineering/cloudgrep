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
	AccountId       string
	Region          string
	Type            string
	DisplayIdPrefix string
	Tags            model.Tags
	RawData         map[string]any
}

func (f ResourceFilter) String() string {
	var parts []string

	for _, matcher := range f.matchers() {
		if !matcher.present() {
			continue
		}

		parts = append(parts, matcher.stringer())
	}

	fields := strings.Join(parts, ", ")
	return fmt.Sprintf("ResourceFilter{%s}", fields)
}

func (f ResourceFilter) Matches(resource model.Resource) bool {
	for _, matcher := range f.matchers() {
		if !matcher.present() {
			continue
		}

		if !matcher.match(resource) {
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

// PartialFilter returns a more detailed filtering of the resources, with filter field processed separately.
// This can aid in debugging to determine why a resource isn't matching a given filter.
func (f ResourceFilter) PartialFilter(in []model.Resource) map[string][]model.Resource {
	output := make(map[string][]model.Resource)

	for _, matcher := range f.matchers() {
		if !matcher.present() {
			continue
		}

		var resources []model.Resource
		for _, res := range in {
			if !matcher.match(res) {
				continue
			}

			resources = append(resources, res)
		}

		output[matcher.name] = resources
	}

	return output
}

type resourceFilterMatcher struct {
	name     string
	present  func() bool
	stringer func() string
	match    func(model.Resource) bool
}

func (f ResourceFilter) matchers() []resourceFilterMatcher {
	p := func(present bool) func() bool {
		return func() bool {
			return present
		}
	}
	s := func(format string, val any) func() string {
		return func() string {
			return fmt.Sprintf(format, val)
		}
	}
	return []resourceFilterMatcher{
		{
			name:     "AccountId",
			present:  p(f.AccountId != ""),
			stringer: s("AccountId=%s", f.AccountId),
			match:    func(r model.Resource) bool { return f.AccountId == r.AccountId },
		},
		{
			name:     "Type",
			present:  p(f.Type != ""),
			stringer: s("Type=%s", f.Type),
			match:    func(r model.Resource) bool { return f.Type == r.Type },
		},
		{
			name:     "Region",
			present:  p(f.Region != ""),
			stringer: s("Region=%s", f.Region),
			match:    func(r model.Resource) bool { return f.Region == r.Region },
		},
		{
			name:     "DisplayIdPrefix",
			present:  p(f.DisplayIdPrefix != ""),
			stringer: s("DisplayIdPrefix=%s", f.DisplayIdPrefix),
			match:    func(r model.Resource) bool { return strings.HasPrefix(r.EffectiveDisplayId(), f.DisplayIdPrefix) },
		},
		{
			name:    "Tags",
			present: p(f.Tags != nil), // Empty non-nil slice has special meaning
			stringer: func() string {
				if len(f.Tags) == 0 {
					return "Tags=[]"
				}

				var parts []string
				for _, tag := range f.Tags {
					if tag.Value == "" {
						parts = append(parts, fmt.Sprintf("Tags[%s]", tag.Key))
					} else {
						parts = append(parts, fmt.Sprintf("Tags[%s]=%s", tag.Key, tag.Value))
					}
				}

				return strings.Join(parts, ", ")
			},
			match: func(r model.Resource) bool {
				// Treat empty slice as special "no tags" filter
				if len(f.Tags) == 0 {
					return len(r.Tags) == 0
				}

				tagMap := make(map[string]string)
				for _, tag := range r.Tags {
					tagMap[tag.Key] = tag.Value
				}

				for _, tag := range f.Tags {
					val, has := tagMap[tag.Key]
					if !has {
						return false
					}

					if tag.Value == "" {
						continue
					}

					if strings.TrimSpace(val) != strings.TrimSpace(tag.Value) {
						return false
					}
				}

				return true
			},
		},
		{
			name:    "RawData",
			present: p(len(f.RawData) > 0),
			stringer: func() string {
				rawParts := make([]string, 0, len(f.RawData))

				for key, val := range f.RawData {
					pair := fmt.Sprintf("%s=%v", key, val)
					rawParts = append(rawParts, pair)
				}
				//sorting ensures consistent output for testing
				sort.Strings(rawParts)

				data := strings.Join(rawParts, ", ")
				return fmt.Sprintf("RawData={%s}", data)
			},
			match: func(r model.Resource) bool {
				var raw map[string]any
				err := json.Unmarshal(r.RawData, &raw)
				if err != nil {
					panic(fmt.Errorf("cannot parse model.Resource.RawData: %s", r.Id))
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

				return true
			},
		},
	}
}
