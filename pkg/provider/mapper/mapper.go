package mapper

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

//Mapper defines rules on how to map a Go Type to a model.Resource
type Mapper struct {
	logger *zap.Logger
	//key is Go Type name, value is mapping configuration
	Mappings map[string]Mapping
}

type Mapping struct {
	Type string `yaml:"type"`
	//optional: if the Resource Type is different that the Golang type
	ResourceType string `yaml:"resourceType"`
	//optional: if the ID field is not called 'Id'
	IdField string `default:"Id" yaml:"idField"`
	//optional: override the TagField from MappingConfig
	TagField TagField `yaml:"tagField"`
	//optional: override the IgnoredFields from MappingConfig
	IgnoredFields []string `yaml:"ignoredFields"`
	//method implementating fetching the resources, the provider must implement it
	Impl string `yaml:"impl"`
}

type Config struct {
	TagField      TagField  `yaml:"tagField"`
	IgnoredFields []string  `yaml:"ignoredFields"`
	Mappings      []Mapping `yaml:"mappings"`
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

func LoadConfig(data []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func New(config Config, logger zap.Logger, providerValue reflect.Value) (Mapper, error) {
	logger.Sugar().Infow("Loading Mappings", zap.String("provider", fmt.Sprintf("%T", providerValue.Interface())))
	tagField := config.TagField
	ignoredFields := config.IgnoredFields
	mapper := Mapper{logger: &logger}
	mapper.Mappings = make(map[string]Mapping)
	for _, mapping := range config.Mappings {
		if mapping.Type == "" {
			return Mapper{}, fmt.Errorf("Mapping %+v is mising 'Type'", mapping)
		}
		if mapping.ResourceType == "" {
			mapping.ResourceType = mapping.Type
		}
		if mapping.IdField == "" {
			mapping.IdField = "Id"
		}
		if mapping.TagField == (TagField{}) {
			//use top level config
			mapping.TagField = tagField
		}
		if len(mapping.IgnoredFields) == 0 {
			//use top level config
			mapping.IgnoredFields = ignoredFields
		}
		//always ignore Tags - they have their own model
		mapping.IgnoredFields = append(mapping.IgnoredFields, mapping.TagField.Name)

		logger.Sugar().Debugf("Loading Mapping %+v", mapping)
		//find the implementation method
		method := providerValue.MethodByName(mapping.Impl)
		if reflect.ValueOf(method).IsZero() {
			return Mapper{}, fmt.Errorf("could not find a method called '%v' on '%T'", mapping.Impl, providerValue.Interface())
		}

		mapper.Mappings[mapping.Type] = mapping

	}
	if len(mapper.Mappings) == 0 {
		return Mapper{}, fmt.Errorf("no mapping loaded for '%T'", providerValue.Interface())
	}
	return mapper, nil
}

//ToRessource generate a Resource by using reflection
// all fields will become properties
// if there is a Tags field, this will become Tags
func (m Mapper) ToRessource(x any, region string) (model.Resource, error) {

	t := reflect.TypeOf(x)
	// key is package name + Type.name to prevent duplicated keys
	key := fmt.Sprintf("%v/%v", path.Dir(t.PkgPath()), t.String())
	mapping, found := m.Mappings[key]
	if !found {
		return model.Resource{}, fmt.Errorf("could not find a mapping definition for type '%v'", t.String())
	}

	var properties []model.Property
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		value := reflect.ValueOf(x).FieldByName(name)
		if field.IsExported() {
			properties = append(properties, getProperties(name, value, mapping.IgnoredFields, 3)...)
		}
	}

	// generate id field
	var id string
	for _, p := range properties {
		if p.Name == mapping.IdField {
			id = p.Value
			break
		}
	}
	if id == "" {
		return model.Resource{}, fmt.Errorf("could not find id field '%v' for type '%T", mapping.IdField, t)
	}

	// generate tags field
	tagsValue := reflect.ValueOf(x).FieldByName(mapping.TagField.Name)
	tags := append(tags, getTags(tagsValue, mapping.TagField)...)

	return model.Resource{
		Id:         id,
		Region:     region,
		Type:       mapping.ResourceType,
		Properties: properties,
		Tags:       tags,
	}, nil
}

//FetchResources calls the implementation method on the Mapping and returns the resources
func (m Mapper) FetchResources(ctx context.Context, mapping Mapping, providerValue reflect.Value, region string) ([]*model.Resource, error) {

	m.logger.Sugar().Infow("Fetching resources",
		zap.String("provider", fmt.Sprintf("%T", providerValue.Interface())),
		zap.String("ResourceType", mapping.ResourceType),
		zap.String("Region", region),
	)
	var resources []*model.Resource
	// call the method to fetch the resources
	method := providerValue.MethodByName(mapping.Impl)
	result := method.Call([]reflect.Value{reflect.ValueOf(ctx)})
	//generate a error message to avoid duplication in code below
	errorMessage := fmt.Sprintf("Method '%T%v' has the wrong return type", providerValue.Interface(), mapping.Impl)
	for _, v := range result {
		switch v.Kind() {
		case reflect.Slice:
			//convert all slice elements to resources
			for i := 0; i < v.Len(); i++ {
				any := v.Index(i).Interface()
				resource, err := m.ToRessource(any, region)
				if err != nil {
					return []*model.Resource{}, err
				}
				resources = append(resources, &resource)
			}
		case reflect.Interface:
			// an error was returned
			err, ok := v.Interface().(error)
			if ok {
				return []*model.Resource{}, err
			}
		default:
			return []*model.Resource{}, fmt.Errorf(errorMessage)
		}
	}

	m.logger.Sugar().Infow("Fetched resources",
		zap.String("provider", fmt.Sprintf("%T", providerValue.Interface())),
		zap.String("ResourceType", mapping.ResourceType),
		zap.String("Region", region),
		zap.Int("Count", len(resources)),
	)
	return resources, nil
}

func getProperties(name string, v reflect.Value, ignoredFields []string, maxRecursion int) []model.Property {

	if util.Contains(ignoredFields, name) || maxRecursion <= 0 {
		//ignore this field
		return []model.Property{}
	}

	emptyProp := []model.Property{{Name: name, Value: ""}}

	switch v.Kind() {
	case reflect.Invalid:
		//ignore this field
		return []model.Property{}
	case reflect.Interface, reflect.Ptr:
		if v.IsZero() {
			//empty pointer
			return emptyProp
		}
		//display pointer value
		return getProperties(name, v.Elem(), ignoredFields, maxRecursion)
	case reflect.Slice:
		if v.IsZero() {
			//empty slice
			return emptyProp
		}
		//return a distinct property for each slice element
		//ex: Subnets=[a,b] -> Subnets[0]=a Subnets[1]=b
		var props []model.Property
		for i := 0; i < v.Len(); i++ {
			props = append(props,
				getProperties(
					fmt.Sprintf("%v[%d]", name, i), v.Index(i), ignoredFields, maxRecursion-1)...)
		}
		return props
	case reflect.Struct:
		defaultFormatVal := fmt.Sprintf("%v", v)
		if !strings.Contains(defaultFormatVal, "{") {
			// this looks like a custom format, use it
			// ex: a datetime might have a nice formatting already implemented
			return []model.Property{{Name: name, Value: defaultFormatVal}}
		}
		//return a distinct property for each struct element
		//ex: IamInstanceProfile{Arn:, Id:} -> IamInstanceProfile["Arn"]=... IamInstanceProfile["Id"]=...
		t := v.Type()
		var props []model.Property
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			props = append(props,
				getProperties(
					fmt.Sprintf("%v[%v]", name, field.Name), v.Field(i), ignoredFields, maxRecursion-1)...)
		}
		return props
	default:
		return []model.Property{{Name: name, Value: fmt.Sprintf("%v", v)}}
	}
}

func getTags(v reflect.Value, tagField TagField) []model.Tag {
	noTags := []model.Tag{}

	switch v.Kind() {
	case reflect.Invalid:
		return noTags
	case reflect.Interface, reflect.Ptr:
		if v.IsZero() {
			//empty pointer
			return noTags
		}
		//display pointer value
		return getTags(v.Elem(), tagField)
	case reflect.Slice:
		if v.IsZero() {
			//empty slice
			return noTags
		}
		//return a distinct Tag for each slice element
		//ex: Tags=[a,b] -> Tag=a Tag=b
		var tags []model.Tag
		for i := 0; i < v.Len(); i++ {
			tags = append(tags,
				getTags(v.Index(i), tagField)...)
		}
		return tags
	case reflect.Struct:
		//expects a Key and Value
		key := getPtrVal(v.FieldByName(tagField.Key))
		value := getPtrVal(v.FieldByName(tagField.Value))
		keyStr := fmt.Sprintf("%v", key)
		valStr := fmt.Sprintf("%v", value)
		// we have a tag
		return []model.Tag{{Key: keyStr, Value: valStr}}
	default:
		return noTags
	}

}

func getPtrVal(v reflect.Value) reflect.Value {
	if v.IsValid() && !v.IsZero() && v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}
