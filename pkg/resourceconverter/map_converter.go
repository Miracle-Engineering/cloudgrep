package resourceconverter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/run-x/cloudgrep/pkg/model"
)

type MapConverter struct {
	ResourceType string
	TagField     TagField
	IdField      string
}

func (mc *MapConverter) ToResource(ctx context.Context, x any, res *model.Resource) error {
	xKind := reflect.TypeOf(x).Kind()
	if xKind != reflect.Map {
		return fmt.Errorf("invalid format %v, expected map", xKind)
	}
	marshaledMap, err := json.Marshal(x)
	if err != nil {
		return err
	}
	var xConverted map[string]any
	err = json.Unmarshal(marshaledMap, &xConverted)
	if err != nil {
		return err
	}

	res.RawData = marshaledMap

	// get the id field
	id, ok := xConverted[mc.IdField]
	if !ok {
		return fmt.Errorf("could not find id field %v in map %v", mc.IdField, xConverted)
	}
	idString := fmt.Sprintf("%v", id)
	res.Id = idString

	// generate tags field
	if !mc.TagField.IsZero() {
		//use field
		tagsValue, ok := xConverted[mc.TagField.Name]
		if !ok {
			return fmt.Errorf("Could not find tag field '%v' for type '%v", mc.TagField.Name, mc.ResourceType)
		}

		res.Tags = getTags(reflect.ValueOf(tagsValue), mc.TagField)
	}

	return nil
}
