package registry

import "reflect"

var registered []registryEntry

type registryEntry struct {
}

func register() {

}

type Registry struct {
	registrations map[string][]TypeRegistration
}

func (r *Registry) RegisterService(name string, types []TypeRegistration) {
	if r.registrations == nil {
		r.registrations = make(map[string][]TypeRegistration)
	}

	r.registrations[name] = types
}

type TypeRegistration struct {
	Name          string
	APIType       reflect.Type
	IDField       reflect.StructField
	IgnoredFields []string
	TagField      reflect.StructField // TODO: Make this more comprehensive
}
