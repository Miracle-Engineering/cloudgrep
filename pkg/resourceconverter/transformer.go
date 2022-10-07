package resourceconverter

import (
	"context"
	"fmt"
	"strconv"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/util"
)

// TransformFunc is a function that receives a raw SDK value and uses it to mutate the passed model.Resource in some way.
type TransformFunc[T any] func(context.Context, T, *model.Resource) error

// TransformResourceFunc is a function that modifies a model.Resource in some way
type TransformResourceFunc func(context.Context, *model.Resource) error

// Transformers is a constructed sequence of TransformFunc and TransformResourceFunc to modify created model.Resources.
type Transformers[T any] struct {
	// *robots in disguise*
	entries []transformerEntry[T]
}

type transformerEntry[T any] struct {
	name string
	f    TransformFunc[T]
}

// Add adds new TransformFuncs to the list, automatically assigning names to each.
func (t *Transformers[T]) Add(funcs ...TransformFunc[T]) {
	for _, f := range funcs {
		entry := transformerEntry[T]{
			name: strconv.Itoa(len(t.entries)),
			f:    f,
		}
		t.entries = append(t.entries, entry)
	}
}

// AddResource adds new TransformResourceFuncs to the list, automatically assigning names to each.
func (t *Transformers[T]) AddResource(funcs ...TransformResourceFunc) {
	for _, f := range funcs {
		resourceFunc := f
		genericFunc := func(ctx context.Context, _ T, res *model.Resource) error {
			return resourceFunc(ctx, res)
		}
		t.Add(genericFunc)
	}
}

// AddNamed adds the given TransformFunc to the list under the specified name.
func (t *Transformers[T]) AddNamed(name string, f TransformFunc[T]) {
	entry := transformerEntry[T]{
		name: name,
		f:    f,
	}
	t.entries = append(t.entries, entry)
}

// AddNamedResource adds the given TransformResourceFunc to the list under the specified name.
func (t *Transformers[T]) AddNamedResource(name string, f TransformResourceFunc) {
	genericFunc := func(ctx context.Context, _ T, res *model.Resource) error {
		return f(ctx, res)
	}

	t.AddNamed(name, genericFunc)
}

// AddTags is a convienience function to add a tag func as a transformer
func (t *Transformers[T]) AddTags(f TagFunc[T]) {
	t.AddNamed("tags", TagTransformer(f))
}

// Apply applies all transform funcs in order to the specified raw SDK value and model.Resource.
func (t Transformers[T]) Apply(ctx context.Context, raw T, resource *model.Resource) error {
	for _, entry := range t.entries {
		err := entry.f(ctx, raw, resource)
		if err != nil {
			err = fmt.Errorf("transformer[%s] failed to apply: %w", entry.name, err)
			return util.AddStackTrace(err)
		}
	}

	return nil
}

// TagFunc is a function that returns the tags for a given SDK value.
type TagFunc[T any] func(context.Context, T) (model.Tags, error)

// TagTransformer converts a TagFunc[T] into a TransformFunc[T]
func TagTransformer[T any](f TagFunc[T]) TransformFunc[T] {
	return func(ctx context.Context, raw T, resource *model.Resource) error {
		tags, err := f(ctx, raw)
		if err != nil {
			return err
		}

		resource.Tags = append(resource.Tags, tags...)
		return nil
	}
}
