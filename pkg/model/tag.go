package model

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

type Tag struct {
	ResourceId string `json:"-" gorm:"primaryKey"`
	Key        string `json:"key" gorm:"primaryKey"`
	Value      string `json:"value"`
	//when used as a filter indicates to look for resources without this tag
	Exclude bool `json:"-"`
}

type Tags []Tag

func (t Tag) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("key", t.Key)
	enc.AddString("value", t.Value)
	return nil
}

func (t Tags) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, tag := range t {
		if err := enc.AppendObject(tag); err != nil {
			return err
		}
	}
	return nil
}

func newTags(m map[string]string) Tags {
	var tags []Tag
	for k, str := range m {
		//a filter on a tag can have multiple values
		values := strings.Split(str, ",")
		for _, v := range values {
			tags = append(tags, Tag{Key: k, Value: v})
		}
	}
	return tags
}

//DistinctKeys returns number of distinct keys - if one tag is specified with two values - it counts as 1 distinct
func (t Tags) DistinctKeys() int {
	counter := make(map[string]int)
	for _, tag := range t {
		counter[tag.Key]++
	}
	return len(counter)
}

//clean removes unexported fields
func (t Tag) clean() Tag {
	return Tag{
		Key:   t.Key,
		Value: t.Value,
	}
}
func (t Tags) Clean() Tags {
	var tags Tags
	for _, tag := range t {
		tags = append(tags, tag.clean())
	}
	return tags
}

func (t Tags) Empty() bool {
	return len(t) == 0
}
