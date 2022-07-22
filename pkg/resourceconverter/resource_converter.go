package resourceconverter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/run-x/cloudgrep/pkg/model"
)

// ResourceFactory is a factory function configured on a ResourceConverter to create a new resource
// with the standard, known fields filed, namely type, region, and account ID.
type ResourceFactory func() model.Resource

type ResourceConverter interface {
	ToResource(context.Context, any, model.Tags) (model.Resource, error)
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
	ResourceFactory ResourceFactory
	IdField         string
	DisplayIdField  string
	TagField        TagField
}

func (rc *ReflectionConverter) ToResource(ctx context.Context, x any, tags model.Tags) (model.Resource, error) {
	if rc.ResourceFactory == nil {
		panic(errors.New("expected ResourceFactory to be set"))
	}

	resource := rc.ResourceFactory()

	resource.Id = rc.findId(rc.IdField, x)
	if resource.Id == "" {
		return model.Resource{}, fmt.Errorf("could not find id field '%v' for type '%v'", rc.IdField, resource.Type)
	}

	if err := rc.loadDisplayId(x, &resource); err != nil {
		return model.Resource{}, err
	}

	// generate tags field
	if !rc.TagField.IsZero() {
		//use field
		tagsValue := reflect.ValueOf(x).FieldByName(rc.TagField.Name)
		if !tagsValue.IsValid() {
			return model.Resource{}, fmt.Errorf("could not find tag field '%v' for type '%v'", rc.TagField.Name, resource.Type)
		}
		tags = append(tags, getTags(tagsValue, rc.TagField)...)
	}

	resource.Tags = append(resource.Tags, tags...)

	marshaledStruct, err := json.Marshal(x)
	if err != nil {
		return model.Resource{}, err
	}

	resource.RawData = marshaledStruct

	return resource, nil
}

func (rc *ReflectionConverter) loadDisplayId(x any, resource *model.Resource) error {
	if rc.DisplayIdField == "" {
		return nil
	}

	id := rc.findId(rc.DisplayIdField, x)

	if id == "" {
		return fmt.Errorf("could not find display id field '%v' for type '%v'", rc.DisplayIdField, resource.Type)
	}

	resource.DisplayId = id

	return nil
}

func (rc *ReflectionConverter) findId(fieldName string, x any) string {
	v := reflect.ValueOf(x)
	if v.IsZero() {
		return ""
	}

	fieldPtrRef := v.FieldByName(fieldName)
	if !fieldPtrRef.IsValid() || fieldPtrRef.IsZero() {
		// Field not on x
		return ""
	}

	fieldRef := reflect.Indirect(fieldPtrRef)
	if !fieldRef.IsValid() || fieldRef.IsZero() {
		return ""
	}

	return fmt.Sprint(fieldRef.Interface())
}
