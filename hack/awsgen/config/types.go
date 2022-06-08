package config

// Root is the schema for the root config.yaml file, which is currently just a list of services
type Root struct {
	// Services is a list of service names (where each is a file with a basename of `NAME.yaml`)
	Services []string `yaml:"services"`
}

// Config contains the aggregated service configuration after parsing.
type Config struct {
	Services []Service
}

// Service defines the configuration for a specific AWS service
type Service struct {
	// Name is the canonical initialized name of the service, in lower case.
	// The name is used for constructing the type identifiers of each resource
	Name string

	// ServicePackage is the name of the service package used by the github.com/aws/aws-sdk-go-v2 module.
	// For example, for the `elb` service, the ServicePackage is `elasticloadbalancingv2`.
	// Defaults to Name.
	ServicePackage string `yaml:"servicePackage"`

	// Global controls whether or not all types in this service default to global, but can be overriden on a per-type basis.
	// A global service is one where resources are not defined in a specific region.
	Global bool `yaml:"global"`

	// Types is the configuration for each type of resource
	Types []Type `yaml:"types"`
}

type Type struct {
	// Name is the name of this resource type.
	// It has the same formatting requirements as an exported Go identifier.
	Name string `yaml:"name"`

	// Global overrides the Service.Global setting
	Global *bool `yaml:"global"`

	// ListAPI contains the configuration used for listing all resources of this type.
	ListAPI ListAPI `yaml:"listApi"`

	// GetTagsAPI contains the configuration used for pulling the tags for each resource of this type.
	// It is not required if the ListAPI is able to pull tags itself.
	GetTagsAPI GetTagsAPI `yaml:"getTagsApi"`
}

// ListAPI is configuration for calling an AWS API to list resources of a type
type ListAPI struct {
	// Call is the AWS API call to make within the service
	Call string `yaml:"call"`

	// InputOverrides stores configuration for setting input fields
	InputOverrides InputOverrides `yaml:"inputOverrides"`

	// Pagination should be set to true if this API has pagination support.
	Pagination bool `yaml:"pagination"`

	// OutputKey sets the "path" to the list of resources within the API response.
	// Each item must be a valid Go identifier.
	// TODO: This should be a NestedField instead of a string slice, to support pointers and non-slices.
	OutputKey []string `yaml:"outputKey"`

	// IDField points to the field within each resource struct that stores the ID of that resource
	IDField Field `yaml:"id"`

	// Tags configures where to look for tags on this resource.
	// Required if the getTagsApi is not configured.
	Tags *TagField `yaml:"tags"`
}

// GetTags is configuration for calling an AWS API to get the tags on a particular resource
type GetTagsAPI struct {
	// ResourceType is the name of the Go struct type within the services "types" package that represents this resource.
	// This type must be what is returned by the list API.
	ResourceType string `yaml:"type"`

	// Call is the AWS API call to make within the service
	Call string `yaml:"call"`

	// InputIDField is the field within the API call's input where we put the resource's ID (from the ListAPI.IDField)
	InputIDField Field `yaml:"inputIDField"`

	// InputOverrides stores configuration for setting input fields
	InputOverrides InputOverrides `yaml:"inputOverrides"`

	// Tags defines where tags are present in the API call output.
	Tags *TagField `yaml:"tags"`

	// AllowedAPIErrorCodes is a list of error codes to ignore when making the API call (treating the resource as having no tags)
	AllowedAPIErrorCodes []string `yaml:"allowedApiErrorCodes"`
}

// TagField defines where tags can be found and how they are accessed.
type TagField struct {
	// Field defines where the tags are located.
	Field NestedField `yaml:"field"`

	/*
		Style sets the style the API uses to store the tags.
		The supported styles are "map" and "struct".
		"map" stores the tags in a map[string]string mapping.
		"struct" stores the tags in something that looks like []struct{Key string; Value string}
	*/
	Style string `yaml:"style"`

	// Pointer sets whether the keys and values are pointers to strings instead of just plain strings.
	// Only supported with the "struct" style.
	Pointer bool `yaml:"pointer"`

	// Key sets the name of the struct field that holds the tag's key.
	Key string `yaml:"key"`

	// Value sets the name of the struct field that holds the tag's value.
	Value string `yaml:"value"`
}

// NestedField is a list of nested Fields, where each successive field is a subfield of the previous.
type NestedField []Field

// Field holds a reference to a field in a struct
type Field struct {
	// Name is the identifier of the field. It must be a valid exported Go identifier.
	Name string `yaml:"name"`

	// SliceType sets the type of the slice this field is, if it is a slice. If not a slice, must be the empty string.
	SliceType string `yaml:"sliceType"`

	// Pointer controls if this field is a pointer that must be dereferenced.
	// Cannot be used with SliceType
	Pointer bool `yaml:"pointer"`
}

// InputOverrides configures functions to call to set values in the API input struct.
type InputOverrides struct {
	// FieldFuncs is a mapping of input field names to function names for setting single fields on the input struct.
	// Each function is called without any arguments, and its return type must be the type of the field.
	FieldFuncs map[string]string `yaml:"fieldFuncs"`

	// FullFuncs is a list functions to call with a pointer to the input struct.
	// Each function must have a return type of an error.
	FullFuncs []string `yaml:"fullFuncs"`
}
