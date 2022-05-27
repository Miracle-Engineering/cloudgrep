package testingutil

func FilterFunc[T any](in []T, predicate func(T) bool) []T {
	out := make([]T, 0, len(in))
	for _, val := range in {
		if predicate(val) {
			out = append(out, val)
		}
	}

	return out
}
