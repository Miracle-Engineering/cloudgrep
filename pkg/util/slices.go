package util

func Chunks[T any](xs []T, chunkSize int) [][]T {
	divided := make([][]T, (len(xs)+chunkSize-1)/chunkSize)
	if len(xs) == 0 {
		return divided
	}
	prev := 0
	i := 0
	till := len(xs) - chunkSize
	for prev < till {
		next := prev + chunkSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}
	divided[i] = xs[prev:]
	return divided
}
