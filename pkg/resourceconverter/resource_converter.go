package resourceconverter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/run-x/cloudgrep/pkg/model"
)

type ResourceConverter interface {
	ToResource(context.Context, any, *model.Resource) error
}

type TagField struct {
	//how to fetch the tags
	//field name
	Name string `yaml:"name"`
	//name of key attribute
	Key string `yaml:"key"`
	//name of value attribute
	Value string `yaml:"value"`
}

func (f TagField) IsZero() bool {
	return f.Name == ""
}

type ReflectionConverter struct {
	ResourceType string
	TagField     TagField
	IdField      string
}

func (rc *ReflectionConverter) ToResource(ctx context.Context, x any, res *model.Resource) error {
	if res == nil {
		panic("unexpected nil model.Resource")
	}

	t := reflect.TypeOf(x)

	if err := rc.loadId(t, x, res); err != nil {
		return err
	}

	if err := rc.loadTags(x, res); err != nil {
		return err
	}

	if err := rc.loadRaw(x, res); err != nil {
		return err
	}

	return nil
}

func (rc *ReflectionConverter) loadId(t reflect.Type, x any, res *model.Resource) error {
	if rc.IdField == "" {
		return nil
	}

	var id string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		if name == rc.IdField {
			fieldPtrRef := reflect.ValueOf(x).FieldByName(name)
			fieldRef := reflect.Indirect(fieldPtrRef)
			if !fieldRef.IsZero() {
				id = fmt.Sprintf("%v", fieldRef.Interface())
			}
			break
		}
	}

	if id == "" {
		return fmt.Errorf("could not find id field '%v' for type '%v", rc.IdField, rc.ResourceType)
	}

	res.Id = id

	return nil
}

func (rc *ReflectionConverter) loadTags(x any, res *model.Resource) error {
	if rc.TagField.IsZero() {
		return nil
	}

	tagsValue := reflect.ValueOf(x).FieldByName(rc.TagField.Name)
	if !tagsValue.IsValid() {
		return fmt.Errorf("Could not find tag field '%v' for type '%v", rc.TagField.Name, rc.ResourceType)
	}
	tags := getTags(tagsValue, rc.TagField)

	res.Tags = tags

	return nil
}

func (rc *ReflectionConverter) loadRaw(x any, res *model.Resource) error {
	marshaledStruct, err := json.Marshal(x)

	if err != nil {
		return err
	}

	res.RawData = marshaledStruct
	return nil
}
