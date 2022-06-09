package model

import (
	"go.uber.org/zap/zapcore"
)

type Tag struct {
	ResourceId string `json:"-" gorm:"primaryKey"`
	Key        string `json:"key" gorm:"primaryKey"`
	Value      string `json:"value"`
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
func (t Tags) Find(key string) *Tag {
	for _, tag := range t {
		if tag.Key == key {
			return &tag
		}
	}
	return nil
}

//Delete deletes a tag from the list by it's key
func (t Tags) Delete(key string) Tags {
	var result Tags
	for _, tag := range t {
		if tag.Key != key {
			result = append(result, tag)
		}
	}
	return result
}
