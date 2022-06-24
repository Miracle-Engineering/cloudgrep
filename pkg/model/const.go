package model

const (
	//FieldMissing used in a query, means that the resource should not have this field defined
	FieldMissing = "(missing)"
	//NullValue used in a query, means that the resource should have this field defined
	FieldPresent = "(not null)"

	//CountValueIgnored means that the current value is ignored in the current query
	CountValueIgnored = "-"

	//name of the field groups as shown in API
	FieldGroupCore = "core"
	FieldGroupTags = "tags"

	//event status as shown in API
	EventStatusFetching string = "fetching"
	EventStatusFailed   string = "failed"
	EventStatusSuccess  string = "success"
	EventStatusLoaded   string = "loaded"

	//event type as shown in API
	EventTypeEngine   string = "engine"
	EventTypeProvider string = "provider"
	EventTypeResource string = "resource"
)
