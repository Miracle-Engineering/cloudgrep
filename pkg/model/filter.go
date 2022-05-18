package model

import (
	"go.uber.org/zap/zapcore"
)

type Filter struct {
	Tags Tags
}

func EmptyFilter() Filter {
	return Filter{}
}

func (f Filter) IsEmpty() bool {
	return len(f.Tags) == 0
}

func (f Filter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return enc.AddArray("tags", f.Tags)
}

func NewFilter(tags map[string]string, excludeTags []string) Filter {
	return Filter{Tags: newTags(tags, excludeTags)}
}

//TagIncludes returns the tags to include
func (f Filter) TagsInclude() Tags {
	var tags Tags
	for _, tag := range f.Tags {
		if !tag.Exclude {
			tags = append(tags, tag)
		}
	}
	return tags
}

//TagsExclude returns the tags to exclude
func (f Filter) TagsExclude() Tags {
	var tags Tags
	for _, tag := range f.Tags {
		if tag.Exclude {
			tags = append(tags, tag)
		}
	}
	return tags
}
