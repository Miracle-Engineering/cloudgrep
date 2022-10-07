package util

import (
	"reflect"
	"strings"
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestImportSet_Add(t *testing.T) {
	var set ImportSet
	assert.Empty(t, set)

	set.Add(Import{Path: "foo", As: "bar"})
	assert.Len(t, set, 1)
	assert.NotNil(t, set)

	assert.Contains(t, set, Import{Path: "foo", As: "bar"})
}

func TestImportSet_Get(t *testing.T) {
	var set ImportSet
	assert.Empty(t, set.Get())

	set.AddPath("foo")
	set.AddPath("bar")

	expected := []Import{
		{Path: "foo"},
		{Path: "bar"},
	}

	actual := set.Get()
	assert.ElementsMatch(t, expected, actual)
}

func TestImportSet_Merge(t *testing.T) {
	var first ImportSet
	first.AddPath("foo")
	first.AddPath("bar")
	require.Len(t, first, 2)

	var second ImportSet
	second.AddPath("foo")
	second.AddPath("spam")
	require.Len(t, second, 2)

	first.Merge(second)
	assert.Len(t, first, 3)
	assert.Len(t, second, 2)

	expected := []Import{
		{Path: "foo"},
		{Path: "bar"},
		{Path: "spam"},
	}
	assert.ElementsMatch(t, expected, first.Get())
}

func TestGroupImports(t *testing.T) {
	input := []Import{
		{Path: "os"},
		{Path: "runtime/debug"},
		{Path: "foo.com/bar", As: "foo"},
		{Path: "bar.com/foo"},
		{Path: "github.com/juandiegopalomino/cloudgrep/pkg/util"},
	}

	expected := GroupedImports{
		StandardLib: []Import{
			{Path: "os"},
			{Path: "runtime/debug"},
		},
		ThirdParty: []Import{
			{Path: "foo.com/bar", As: "foo"},
			{Path: "bar.com/foo"},
		},
		Module: []Import{
			{Path: "github.com/juandiegopalomino/cloudgrep/pkg/util"},
		},
	}

	actual := GroupImports(input)

	assert.Equal(t, expected, actual)
	assert.Equal(t, 5, actual.Len())

	groups := actual.Groups()
	assert.Equal(t, expected.StandardLib, groups[0])
	assert.Equal(t, expected.ThirdParty, groups[1])
	assert.Equal(t, expected.Module, groups[2])
}

func TestSortImports(t *testing.T) {
	input := []Import{
		{Path: "os"},
		{Path: "strings", As: "a"},
		{Path: "runtime/debug"},
		{Path: "runtime", As: "z"},
		{Path: "foo.com/bar", As: "foo"},
		{Path: "bar.com/foo"},
	}

	expected := []Import{
		{Path: "bar.com/foo"},
		{Path: "foo.com/bar", As: "foo"},
		{Path: "os"},
		{Path: "runtime", As: "z"},
		{Path: "runtime/debug"},
		{Path: "strings", As: "a"},
	}

	actual := slices.Clone(input)
	SortImports(actual)

	assert.Equal(t, expected, actual)
}

func TestModulePrefix(t *testing.T) {
	var i provider.Provider
	typ := reflect.TypeOf(&i).Elem()
	pkg := typ.PkgPath()
	assert.True(t, strings.HasPrefix(pkg, modulePrefix))
}
