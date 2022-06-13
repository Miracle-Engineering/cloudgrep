package resourceconverter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/run-x/cloudgrep/pkg/model"
	"reflect"
)

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

type ReflectionConverter struct {
	ResourceType string
	TagField     TagField
	IdField      string
	Region       string
}

func (rc *ReflectionConverter) ToResource(ctx context.Context, x any, tags model.Tags) (model.Resource, error) {
	t := reflect.TypeOf(x)

	// get the id field
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
		return model.Resource{}, fmt.Errorf("could not find id field '%v' for type '%v", rc.IdField, rc.ResourceType)
	}

	// generate tags field
	if tags == nil {
		//use field
		tagsValue := reflect.ValueOf(x).FieldByName(rc.TagField.Name)
		if !tagsValue.IsValid() {
			return model.Resource{}, fmt.Errorf("Could not find tag field '%v' for type '%v", rc.TagField.Name, rc.ResourceType)
		}
		tags = getTags(tagsValue, rc.TagField)
	}
	marshaledStruct, err := json.Marshal(x)
	if err != nil {
		return model.Resource{}, err
	}
	return model.Resource{
		Id:      id,
		Region:  rc.Region,
		Type:    rc.ResourceType,
		RawData: marshaledStruct,
		Tags:    tags,
	}, nil
}
