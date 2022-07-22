package testingutil

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestResourceFilter_Matches(t *testing.T) {
	tests := []struct {
		name     string
		filter   ResourceFilter
		resource model.Resource
		want     bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "regionMatch",
			want: true,
			filter: ResourceFilter{
				Region: "foo",
			},
			resource: model.Resource{
				Region: "foo",
				Type:   "bar",
			},
		},
		{
			name: "regionMismatch",
			filter: ResourceFilter{
				Region: "foo",
			},
			resource: model.Resource{
				Region: "bar",
			},
		},
		{
			name: "typeMatch",
			want: true,
			filter: ResourceFilter{
				Type: "bar",
			},
			resource: model.Resource{
				Region: "foo",
				Type:   "bar",
			},
		},
		{
			name: "typeMismatch",
			filter: ResourceFilter{
				Type: "bar",
			},
			resource: model.Resource{
				Region: "foo",
				Type:   "spam",
			},
		},
		{
			name: "tagsEmptyMatch",
			want: true,
			filter: ResourceFilter{
				Tags: model.Tags{},
			},
			resource: model.Resource{},
		},
		{
			name: "tagsEmptyMismatch",
			filter: ResourceFilter{
				Tags: model.Tags{},
			},
			resource: model.Resource{
				Tags: model.Tags{
					{Key: "Foo", Value: "Bar"},
				},
			},
		},
		{
			name: "tagsKeyOnlyMatch",
			want: true,
			filter: ResourceFilter{
				Tags: model.Tags{
					{
						Key: "Foo",
					},
				},
			},
			resource: model.Resource{
				Tags: model.Tags{
					{Key: "Foo", Value: "bar"},
					{Key: "Spam", Value: "ham"},
				},
			},
		},
		{
			name: "tagsKeyOnlyMismatch",
			filter: ResourceFilter{
				Tags: model.Tags{
					{
						Key: "Foo",
					},
				},
			},
			resource: model.Resource{
				Tags: model.Tags{
					{Key: "Spam", Value: "ham"},
				},
			},
		},
		{
			name: "tagsMatch",
			want: true,
			filter: ResourceFilter{
				Tags: model.Tags{
					{
						Key:   "Foo",
						Value: "Bar",
					},
				},
			},
			resource: model.Resource{
				Tags: model.Tags{
					{Key: "Foo", Value: "Bar"},
					{Key: "Spam", Value: "ham"},
				},
			},
		},
		{
			name: "tagsMismatch",
			filter: ResourceFilter{
				Tags: model.Tags{
					{
						Key:   "Foo",
						Value: "bar",
					},
				},
			},
			resource: model.Resource{
				Tags: model.Tags{
					{Key: "Foo", Value: "foo"},
					{Key: "Spam", Value: "ham"},
				},
			},
		},
		{
			name: "rawMatch",
			want: true,
			filter: ResourceFilter{
				RawData: map[string]any{
					"foo": "bar",
				},
			},
			resource: model.Resource{
				RawData: []byte(`{"foo":"bar"}`),
			},
		},
		{
			name: "rawMissingMismatch",
			filter: ResourceFilter{
				RawData: map[string]any{
					"foo": "bar",
				},
			},
			resource: model.Resource{
				RawData: []byte(`{"spam":"ham"}`),
			},
		},
		{
			name: "rawMismatch",
			filter: ResourceFilter{
				RawData: map[string]any{
					"foo": "spam",
				},
			},
			resource: model.Resource{
				RawData: []byte(`{"foo":"ham"}`),
			},
		},
		{
			name: "displayIdPrefixMatchNoDisplay",
			want: true,
			filter: ResourceFilter{
				DisplayIdPrefix: "foo-",
			},
			resource: model.Resource{
				Id: "foo-1",
			},
		},
		{
			name: "displayIdPrefixMismatchNoDisplay",
			filter: ResourceFilter{
				DisplayIdPrefix: "bar",
			},
			resource: model.Resource{
				Id: "ba",
			},
		},
		{
			name: "displayIdPrefixMatch",
			want: true,
			filter: ResourceFilter{
				DisplayIdPrefix: "foo-",
			},
			resource: model.Resource{
				Id:        "bar-1",
				DisplayId: "foo-2",
			},
		},
		{
			name: "displayIdPrefixMismatch",
			filter: ResourceFilter{
				DisplayIdPrefix: "bar",
			},
			resource: model.Resource{
				Id:        "bar",
				DisplayId: "foo-2",
			},
		},
		{
			name: "accountIdMatch",
			want: true,
			filter: ResourceFilter{
				AccountId: "foo",
			},
			resource: model.Resource{
				AccountId: "foo",
				Id:        "spam",
			},
		},
		{
			name: "accountIdMismatch",
			filter: ResourceFilter{
				AccountId: "foo",
			},
			resource: model.Resource{
				AccountId: "bar",
				Id:        "spam",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.filter.Matches(test.resource); got != test.want {
				t.Errorf("ResourceFilter.Matches() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestResourceFilter_Filter(t *testing.T) {
	in := []model.Resource{
		{
			Id:     "A",
			Type:   "Foo",
			Region: "bar",
			Tags: model.Tags{
				{
					Key:   "Foo",
					Value: "Bar",
				},
				{
					Key:   "Spam",
					Value: "Hi",
				},
			},
		},
		{
			Id:     "B",
			Type:   "Foo",
			Region: "spam",
			Tags: model.Tags{
				{
					Key:   "Foo",
					Value: "Bar",
				},
				{
					Key:   "Spam",
					Value: "Hi",
				},
			},
		},
		{
			Id:   "C",
			Type: "",
		},
	}

	filter := ResourceFilter{
		Type: "Foo",
		Tags: model.Tags{
			{
				Key:   "Foo",
				Value: "Bar",
			},
			{
				Key: "Spam",
			},
		},
	}

	expectedIds := []string{"A", "B"}

	actual := filter.Filter(in)
	var actualIds []string
	for _, resource := range actual {
		actualIds = append(actualIds, resource.Id)
	}

	assert.ElementsMatch(t, expectedIds, actualIds)

}

func TestResourceFilter_String(t *testing.T) {
	tests := []struct {
		name   string
		filter ResourceFilter
		want   string
	}{
		{
			name: "empty",
			want: "ResourceFilter{}",
		},
		{
			name: "emptyTags",
			filter: ResourceFilter{
				Tags: model.Tags{},
			},
			want: "ResourceFilter{Tags=[]}",
		},
		{
			name: "full",
			filter: ResourceFilter{
				AccountId:       "TestAccount",
				Type:            "A",
				Region:          "B",
				DisplayIdPrefix: "foo-",
				Tags: model.Tags{
					{
						Key: "Foo",
					},
					{
						Key:   "Spam",
						Value: "ham",
					},
				},
				RawData: map[string]any{
					"apple": 1,
					"fool":  "took",
				},
			},
			want: "ResourceFilter{AccountId=TestAccount, Type=A, Region=B, DisplayIdPrefix=foo-, Tags[Foo], Tags[Spam]=ham, RawData={apple=1, fool=took}}",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.filter.String()
			assert.Equal(t, test.want, actual)
		})
	}
}

func TestResourceFilter_matchers_unqiue(t *testing.T) {
	names := make(map[string]struct{})
	f := ResourceFilter{}
	for _, matcher := range f.matchers() {
		if _, has := names[matcher.name]; has {
			t.Errorf("duplicate matcher name: %s", matcher.name)
		}

		names[matcher.name] = struct{}{}
	}
}

func TestResourceFilter_PartialFilter(t *testing.T) {
	tests := []struct {
		name      string
		want      map[string][]string
		filter    ResourceFilter
		resources []model.Resource
	}{
		{
			name: "multiple",
			want: map[string][]string{
				"AccountId": {"foo"},
				"Type":      {"foo", "bar"},
				"Region":    {"bar"},
			},
			filter: ResourceFilter{
				AccountId: "a",
				Type:      "b",
				Region:    "c",
			},
			resources: []model.Resource{
				{
					Id:        "foo",
					AccountId: "a",
					Type:      "b",
					Region:    "d",
				},
				{
					Id:        "bar",
					AccountId: "e",
					Type:      "b",
					Region:    "c",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.filter.PartialFilter(test.resources)
			actualKeys := maps.Keys(actual)
			expectedKeys := maps.Keys(test.want)

			assert.ElementsMatch(t, expectedKeys, actualKeys)

			for name, expectedIds := range test.want {
				actualResources, has := actual[name]
				if !has {
					// Already errored on above ElementsMatch
					continue
				}

				var actualIds []string
				for _, resource := range actualResources {
					actualIds = append(actualIds, resource.Id)
				}

				assert.ElementsMatchf(t, expectedIds, actualIds, "expected ids on matcher %s to match", name)
			}
		})
	}
}

func TestResourceFilter_Match_rawDataPanic(t *testing.T) {
	f := ResourceFilter{
		RawData: map[string]any{
			"foo": "bar",
		},
	}

	r := model.Resource{
		Id:      "spam",
		RawData: []byte("{"),
	}

	assert.PanicsWithError(t, "cannot parse model.Resource.RawData: spam", func() {
		f.Matches(r)
	})
}
