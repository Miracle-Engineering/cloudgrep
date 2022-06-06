package util

import (
	"fmt"
	"strconv"
)

type RecursiveAppend struct {
	Idx  int
	Keys []string
	Root string
	Data map[string]any
}

func (r RecursiveAppend) IsLast() bool {
	return r.Idx >= len(r.Keys)-1
}

func (r RecursiveAppend) IterVar() string {
	if r.Idx > 0 {
		return r.varFor(r.Idx - 1)
	}

	return r.Root
}

func (r RecursiveAppend) Current() string {
	return r.Keys[r.Idx]
}

func (r RecursiveAppend) NextIterVar() string {
	return r.varFor(r.Idx)
}

func (r RecursiveAppend) varFor(i int) string {
	return "item_" + strconv.Itoa(i)
}

func (r RecursiveAppend) Next() (RecursiveAppend, error) {
	if r.IsLast() {
		return r, fmt.Errorf("end of recursive append keys")
	}

	r.Idx++
	return r, nil
}

func (r RecursiveAppend) WithRoot(root string) RecursiveAppend {
	r.Root = root
	return r
}

func (r *RecursiveAppend) SetData(key string, value any) error {
	if r.Data == nil {
		r.Data = make(map[string]any)
	}

	r.Data[key] = value

	return nil
}
