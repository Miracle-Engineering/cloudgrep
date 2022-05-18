package model

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

type Tag struct {
	ResourceId ResourceId `json:"-" gorm:"primaryKey"`
	Key        string     `json:"key" gorm:"primaryKey"`
	Value      string     `json:"value"`
	//when used as a filter indicates to look for resources without this tag
	Exclude bool `json:"-" gorm:"-"`
}

//TagInfo provide information about a searched tag, including count and values found
type TagInfo struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
	//number of resources found with this tag
	Count int `json:"count"`
	//list the resource with this tags
	ResourceIds []ResourceId `json:"resourceIds"`
}

type Tags []Tag
type TagInfos []*TagInfo

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

func newTags(m map[string]string, excludes []string) Tags {
	var tags Tags
	for k, str := range m {
		//a filter on a tag can have multiple values
		values := strings.Split(str, ",")
		if len(k) == 0 {
			continue
		}
		for _, v := range values {
			tags = append(tags, Tag{Key: k, Value: v})
		}
	}
	for _, tag := range excludes {
		if len(tag) == 0 {
			continue
		}
		tags = append(tags, Tag{Key: tag, Value: "*", Exclude: true})
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

//Find finds a TagInfo by Key, return nil if not found
func (ts TagInfos) Find(key string) *TagInfo {
	for _, t := range ts {
		if t.Key == key {
			return t
		}
	}
	return nil
}
