package resourceconverter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/run-x/cloudgrep/pkg/model"
)

type MapConverter struct {
	ResourceType   string
	TagField       TagField
	IdField        string
	DisplayIdField string
	Region         string
}

func (mc *MapConverter) ToResource(ctx context.Context, x any, tags model.Tags) (model.Resource, error) {
	xKind := reflect.TypeOf(x).Kind()
	if xKind != reflect.Map {
		return model.Resource{}, fmt.Errorf("invalid format %v, expected map", xKind)
	}
	marshaledMap, err := json.Marshal(x)
	if err != nil {
		return model.Resource{}, err
	}
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
	idString := fmt.Sprintf("%v", id)

	var displayIdString string
	if mc.DisplayIdField != "" {
		displayId, ok := xConverted[mc.DisplayIdField]
		if !ok {
			return model.Resource{}, fmt.Errorf("could not find display id field %v in map %v", mc.DisplayIdField, xConverted)
		}
		displayIdString = fmt.Sprintf("%v", displayId)
	}

	// generate tags field
	if !mc.TagField.IsZero() {
		//use field
		tagsValue, ok := xConverted[mc.TagField.Name]
		if !ok {
			return model.Resource{}, fmt.Errorf("could not find tag field '%v' for type '%v", mc.TagField.Name, mc.ResourceType)
		}

		tags = getTags(reflect.ValueOf(tagsValue), mc.TagField)
	}

	if err != nil {
		return model.Resource{}, err
	}
	return model.Resource{
		Id:        idString,
		DisplayId: displayIdString,
		Region:    mc.Region,
		Type:      mc.ResourceType,
		RawData:   marshaledMap,
		Tags:      tags,
	}, nil
}
