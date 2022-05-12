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

func NewFilter(tags map[string]string) Filter {
	return Filter{Tags: newTags(tags)}
}
