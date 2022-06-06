package config

func (f NestedField) Empty() bool {
	return len(f) == 0
}
