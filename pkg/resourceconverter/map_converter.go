package resourceconverter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
)

type MapConverter struct {
	ResourceFactory ResourceFactory
	IdField         string
	DisplayIdField  string
	TagField        TagField
}

func (mc *MapConverter) ToResource(ctx context.Context, x any, tags model.Tags) (model.Resource, error) {
	if mc.ResourceFactory == nil {
		panic(errors.New("expected ResourceFactory to be set"))
	}

	resource := mc.ResourceFactory()

	xKind := reflect.TypeOf(x).Kind()
	if xKind != reflect.Map {
		return model.Resource{}, fmt.Errorf("invalid format %v, expected map", xKind)
	}
	marshaledMap, err := json.Marshal(x)
	if err != nil {
		return model.Resource{}, err
	}

	resource.RawData = marshaledMap

	var xConverted map[string]any
	err = json.Unmarshal(marshaledMap, &xConverted)
	if err != nil {
		return model.Resource{}, err
	}

	// get the id field
	id, ok := xConverted[mc.IdField]
	if !ok {
		return model.Resource{}, fmt.Errorf("could not find id field %v in map %v", mc.IdField, xConverted)
	}
	resource.Id = fmt.Sprint(id)

	if mc.DisplayIdField != "" {
		displayId, ok := xConverted[mc.DisplayIdField]
		if !ok {
			return model.Resource{}, fmt.Errorf("could not find display id field %v in map %v", mc.DisplayIdField, xConverted)
		}
		resource.DisplayId = fmt.Sprint(displayId)
	}

	// generate tags field
	if !mc.TagField.IsZero() {
		//use field
		tagsValue, ok := xConverted[mc.TagField.Name]
		if !ok {
			return model.Resource{}, fmt.Errorf("could not find tag field '%v' for type '%v", mc.TagField.Name, resource.Type)
		}

		tags = append(tags, getTags(reflect.ValueOf(tagsValue), mc.TagField)...)
	}

	resource.Tags = append(resource.Tags, tags...)

	return resource, nil
}
