package mapper

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/run-x/cloudgrep/pkg/model"
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
			mapping.TagMethod = findTagMethod(providerValue, mapping.TagImpl)
		}
		mapper.Mappings[mapping.Type] = mapping

	}
	if len(mapper.Mappings) == 0 {
		return Mapper{}, fmt.Errorf("no mapping loaded for '%T'", providerValue.Interface())
	}
	return mapper, nil
}

//ToResource generate a Resource by using reflection
// all fields will be added to raw data
// if there is a Tags field, this will become Tags
func (m Mapper) ToResource(ctx context.Context, x any, region string) (model.Resource, error) {

	t := reflect.TypeOf(x)
	// key is package name + Type.name to prevent duplicated keys
	key := fmt.Sprintf("%v/%v", path.Dir(t.PkgPath()), t.String())
	mapping, found := m.Mappings[key]
	if !found {
		return model.Resource{}, fmt.Errorf("could not find a mapping definition for type '%v'", key)
	}

	// get the id field
	var id string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		if name == mapping.IdField {
			fieldPtrRef := reflect.ValueOf(x).FieldByName(name)
			fieldRef := reflect.Indirect(fieldPtrRef)
			if !fieldRef.IsZero() {
				id = fmt.Sprintf("%v", fieldRef.Interface())
			}
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
	marshaledStruct, err := json.Marshal(x)
	if err != nil {
		return model.Resource{}, err
	}
	return model.Resource{
		Id:      id,
		Region:  region,
		Type:    mapping.ResourceType,
		RawData: marshaledStruct,
		Tags:    tags,
	}, nil
}

//FetchResources calls the implementation method on each Mapping and returns the resources
//this method can return both resources and error - if it has partially worked
func (m Mapper) FetchResources(ctx context.Context, region string) ([]*model.Resource, error) {
	var resources []*model.Resource
	var errors error

	errorsChan := make(chan error)
	resourceListChan := make(chan []*model.Resource)
	for _, mapping := range m.Mappings {
		go func(mapping Mapping) {
			new_resources, err := m.fetchResources(ctx, mapping, region)
			errorsChan <- err
			resourceListChan <- new_resources
		}(mapping)
	}

	errorChan := make(chan error)
	go func() {
		var mainError error
		for range m.Mappings {
			err := <-errorsChan
			if err != nil {
				mainError = multierror.Append(mainError, err)
			}
		}
		errorChan <- mainError
	}()

	resourcesChan := make(chan []*model.Resource)
	go func() {
		var mainResourceList []*model.Resource
		for range m.Mappings {
			new_resources := <-resourceListChan
			mainResourceList = append(mainResourceList, new_resources...)
		}
		resourcesChan <- mainResourceList
	}()

	errors = <-errorChan
	resources = <-resourcesChan
	//if we have no resource and some errors, only return the first error
	//most likely a global issue such as auth or connectivity and would be too verbose
	if len(resources) == 0 && errors != nil {
		if merr, ok := errors.(*multierror.Error); ok {
			errors = merr.Errors[0]
		}
	}
	return resources, errors
}

//fetchResources fetches the resources for a mapping and a region
//note this method can return some resources and an error, if it has partially worked
func (m Mapper) fetchResources(ctx context.Context, mapping Mapping, region string) ([]*model.Resource, error) {
	m.logger.Sugar().Infow("Fetching resources async",
		zap.String("ResourceType", mapping.ResourceType),
		zap.String("Region", region),
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fetchFuncType := mapping.Method.Type()
	outChanType := fetchFuncType.In(1)
	bothOutChanType := reflect.ChanOf(reflect.BothDir, outChanType.Elem())
	bothOutChan := reflect.MakeChan(bothOutChanType, 0)
	outChan := bothOutChan.Convert(outChanType)

	var resources []*model.Resource
	doneCh := make(chan struct{})
	var recvErr error

	go func() {
		defer close(doneCh)

		cases := []reflect.SelectCase{
			{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ctx.Done()),
			},
			{
				Dir:  reflect.SelectRecv,
				Chan: bothOutChan,
			},
		}

		for {
			chosen, recv, ok := reflect.Select(cases)
			if chosen == 0 {
				recvErr = ctx.Err()
				return
			}

			// chosen == 1
			if !ok {
				// Channel closed, so we are done
				return
			}

			val := recv.Interface()
			resource, err := m.ToResource(ctx, val, region)
			if err != nil {
				//store error and keep processing other resources
				recvErr = multierror.Append(recvErr,
					fmt.Errorf(
						"error converting %v result slice to resource: %w", mapping.Impl, err),
				)
				continue
			}
			resources = append(resources, &resource)
		}
	}()

	args := []reflect.Value{reflect.ValueOf(ctx), outChan}

	results := mapping.Method.Call(args)
	err := results[0].Interface()
	if err != nil {
		return nil, err.(error)
	}

	bothOutChan.Close()

	// Wait for all resources to be read
	<-doneCh

	if recvErr != nil {
		return resources, recvErr
	}

	return resources, nil
}
