package util

import (
	"fmt"
	"strconv"
)

// RecursiveAppend is used in templates to give a way to concisely reference nested fields.
type RecursiveAppend[T any] struct {
	Idx  int
	Keys []T
	Root string
	Data map[string]any
}

func (r RecursiveAppend[T]) IsLast() bool {
	return r.Idx >= len(r.Keys)-1
}

func (r RecursiveAppend[T]) IterVar() string {
	if r.Idx > 0 {
		return r.varFor(r.Idx - 1)
	}

	return r.Root
}

func (r RecursiveAppend[T]) Current() T {
	return r.Keys[r.Idx]
}

func (r RecursiveAppend[T]) NextIterVar() string {
	return r.varFor(r.Idx)
}

func (r RecursiveAppend[T]) varFor(i int) string {
	return "item_" + strconv.Itoa(i)
}

func (r RecursiveAppend[T]) Next() (RecursiveAppend[T], error) {
	if r.IsLast() {
		return r, fmt.Errorf("end of recursive append keys")
	}

	r.Idx++
	return r, nil
}

func (r RecursiveAppend[T]) WithRoot(root string) RecursiveAppend[T] {
	r.Root = root
	return r
}

func (r *RecursiveAppend[T]) SetData(key string, value any) error {
	if r.Data == nil {
		r.Data = make(map[string]any)
	}

	r.Data[key] = value

	// Returns an error so it can be used in templates
	return nil
}
