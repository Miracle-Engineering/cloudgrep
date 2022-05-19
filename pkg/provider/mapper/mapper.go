package mapper

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// limit the depth when using reflection to generate the property list
const maxRecursion = 5

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
	//the method is found at runtime using Impl
	Method *reflect.Value
	//optional: method implementating fetching the tags, only needed if there is not a field containing tags
	TagImpl string `yaml:"tagImpl"`
	//the method is found at runtime using TagImpl
	TagMethod *reflect.Value
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

func New(config []byte, logger zap.Logger, providerValue reflect.Value) (Mapper, error) {
	var configStruct Config
	err := yaml.Unmarshal(config, &configStruct)
	if err != nil {
		return Mapper{}, err
	}
	return new(configStruct, logger, providerValue)
}

func new(config Config, logger zap.Logger, providerValue reflect.Value) (Mapper, error) {
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
		//find the implementation method
		mapping.Method = findImplMethod(providerValue, mapping.Impl)
		//find the implementation method for tags
		if mapping.TagImpl != "" {
			mapping.TagMethod = findImplMethod(providerValue, mapping.TagImpl)
		}
		mapper.Mappings[mapping.Type] = mapping

	}
	if len(mapper.Mappings) == 0 {
		return Mapper{}, fmt.Errorf("no mapping loaded for '%T'", providerValue.Interface())
	}
	return mapper, nil
}

func findImplMethod(v reflect.Value, impl string) *reflect.Value {
	//find the implementation method
	method := v.MethodByName(impl)
	if reflect.ValueOf(method).IsZero() {
		panic(fmt.Errorf("could not find a method called '%v' on '%T'", impl, v.Interface()))
	}
	//check return type is (slice, error)
	if t := method.Type(); t.NumOut() != 2 || t.Out(0).Kind().String() != "slice" || t.Out(1).Name() != "error" {
		panic(fmt.Errorf("method %v has invalid return type, expecting ([]any, error)", impl))
	}
	return &method
}

//ToResource generate a Resource by using reflection
// all fields will become properties
// if there is a Tags field, this will become Tags
func (m Mapper) ToResource(ctx context.Context, x any, region string) (model.Resource, error) {

	t := reflect.TypeOf(x)
	// key is package name + Type.name to prevent duplicated keys
	key := fmt.Sprintf("%v/%v", path.Dir(t.PkgPath()), t.String())
	mapping, found := m.Mappings[key]
	if !found {
		return model.Resource{}, fmt.Errorf("could not find a mapping definition for type '%v'", key)
	}

	//generate properties
	var properties []model.Property
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		value := reflect.ValueOf(x).FieldByName(name)
		if field.IsExported() {
			properties = append(properties, getProperties(name, value, mapping.IgnoredFields, maxRecursion)...)
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
		return model.Resource{}, fmt.Errorf("could not find id field '%v' for type '%v", mapping.IdField, key)
	}

	// generate tags field
	var tags model.Tags
	if mapping.TagMethod == nil {
		//use field
		tagsValue := reflect.ValueOf(x).FieldByName(mapping.TagField.Name)
		tags = getTags(tagsValue, mapping.TagField)
	} else {
		//call the defined method
		result := mapping.TagMethod.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(x)})
		for _, v := range result {
			switch v.Kind() {
			case reflect.Slice:
				// got the tags
				for i := 0; i < v.Len(); i++ {
					tag := v.Index(i).Interface().(model.Tag)
					tags = append(tags, tag)
				}
			case reflect.Interface:
				// check if an error was returned
				err, ok := v.Interface().(error)
				if ok {
					//an error getting tags is not blocking, return the resource without tags
					if !strings.Contains(err.Error(), "404") {
						//404 means no tags - ignore this error code
						m.logger.Sugar().Error(err)
					}
				}
			default:
				return model.Resource{}, fmt.Errorf("method '%v' has the wrong return type", mapping.TagImpl)
			}
		}
	}

	return model.Resource{
		Id:         id,
		Region:     region,
		Type:       mapping.ResourceType,
		Properties: properties,
		Tags:       tags,
	}, nil
}

//FetchResources calls the implementation method on each Mapping and returns the resources
//this method can return both resources and error - if it has partially worked
func (m Mapper) FetchResources(ctx context.Context, region string) ([]*model.Resource, error) {
	var resources []*model.Resource
	var errors error
	for _, mapping := range m.Mappings {
		new_resources, err := m.fetchResources(ctx, mapping, region)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
		resources = append(resources, new_resources...)
	}
	//if we have no resource and some errors, only return the first error
	//most likely a global issue such as auth or connectivity and would be too verbose
	if len(resources) == 0 && errors != nil {
		if merr, ok := errors.(*multierror.Error); ok {
			errors = merr.Errors[0]
		}
	}
	return resources, errors
}

func (m Mapper) fetchResources(ctx context.Context, mapping Mapping, region string) ([]*model.Resource, error) {
	m.logger.Sugar().Infow("Fetching resources",
		zap.String("ResourceType", mapping.ResourceType),
		zap.String("Region", region),
	)
	var resources []*model.Resource
	// call the method to fetch the resources
	result := mapping.Method.Call([]reflect.Value{reflect.ValueOf(ctx)})
	//generate a error message to avoid duplication in code below
	for _, v := range result {
		switch v.Kind() {
		case reflect.Slice:
			//convert all slice elements to resources
			for i := 0; i < v.Len(); i++ {
				any := v.Index(i).Interface()
				resource, err := m.ToResource(ctx, any, region)
				if err != nil {
					return nil, fmt.Errorf("error converting %v result slice to resource: %w", mapping.Impl, err)
				}
				resources = append(resources, &resource)
			}
		case reflect.Interface:
			// an error was returned
			err, ok := v.Interface().(error)
			if ok {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("method '%v' has the wrong return type", mapping.Impl)
		}
	}

	m.logger.Sugar().Infow("Fetched resources",
		zap.String("ResourceType", mapping.ResourceType),
		zap.String("Region", region),
		zap.Int("Count", len(resources)),
	)
	return resources, nil
}

func getProperties(name string, v reflect.Value, ignoredFields []string, maxRecursion int) []model.Property {

	if util.Contains(ignoredFields, name) || maxRecursion <= 0 {
		//ignore this field
		return nil
	}

	emptyProp := []model.Property{{Name: name, Value: ""}}

	switch v.Kind() {
	case reflect.Invalid:
		//ignore this field
		return nil
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
	switch v.Kind() {
	case reflect.Invalid:
		return nil
	case reflect.Interface, reflect.Ptr:
		if v.IsZero() {
			//empty pointer
			return nil
		}
		//display pointer value
		return getTags(v.Elem(), tagField)
	case reflect.Slice:
		if v.IsZero() {
			//empty slice
			return nil
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
		return nil
	}

}

func getPtrVal(v reflect.Value) reflect.Value {
	if v.IsValid() && !v.IsZero() && v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}
