package testingutil

// FilterFunc filters a slice based on the predicate function
func FilterFunc[T any](in []T, predicate func(T) bool) []T {
	out := make([]T, 0, len(in))
	for _, val := range in {
		if predicate(val) {
			out = append(out, val)
		}
	}

	return out
}

// Unique returns a new slice with duplicate elements removed and the order preserved.
func Unique[T comparable](in []T) []T {
	out := make([]T, 0, len(in))
	seen := make(map[T]struct{})

	for _, val := range in {
		if _, has := seen[val]; has {
			continue
		}

		seen[val] = struct{}{}
		out = append(out, val)
	}

	return out
}
