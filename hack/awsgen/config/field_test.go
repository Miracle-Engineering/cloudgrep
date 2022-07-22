package config

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_Decoding_string(t *testing.T) {
	var actual Field
	in := `foo`
	expected := Field{Name: "foo"}

	err := yamlUnmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestField_Decoding_mapping(t *testing.T) {
	var actual Field
	in := `{name: foo, sliceType: string}`
	expected := Field{Name: "foo", SliceType: "string"}

	err := yamlUnmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestField_Decoding_mappingInvalid(t *testing.T) {
	in := `{name: foo, pointer: string}`

	err := yamlUnmarshal([]byte(in), &Field{})
	assert.ErrorContains(t, err, "cannot unmarshal !!str `string` into bool")
}

func TestField_Decoding_sequence(t *testing.T) {
	in := "- foo\n- bar"
	err := yamlUnmarshal([]byte(in), &Field{})
	assert.ErrorContains(t, err, "unexpected node kind")
}

func TestField_Decoding_noName(t *testing.T) {
	var actual Field
	in := `{sliceType: string}`
	expected := Field{}

	err := yamlUnmarshal([]byte(in), &actual)
	assert.ErrorContains(t, err, "missing \"name\"")
	assert.Equal(t, expected, actual)
}

func TestField_Zero(t *testing.T) {
	var f *Field
	assert.True(t, f.Zero())

	f = &Field{}
	assert.True(t, f.Zero())

	f.Name = "Foo"
	assert.False(t, f.Zero())
}

func TestTagFieldValidStyles(t *testing.T) {
	for _, val := range TagFieldValidStyles() {
		assert.Regexp(t, regexp.MustCompile("^[a-z]+$"), val)
	}
}

func TestNestedField_Last(t *testing.T) {
	nf := NestedField{}
	assert.Zero(t, nf.Last())

	nf = append(nf, Field{Name: "foo"}, Field{Name: "bar"})
	last := nf.Last()
	assert.Equal(t, "bar", last.Name)
}

func TestNestedField_Decoding_string(t *testing.T) {
	in := "foo"
	expected := NestedField{Field{Name: "foo"}}

	var actual NestedField
	err := yamlUnmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestNestedField_Decoding_map(t *testing.T) {
	in := "name: foo"

	err := yamlUnmarshal([]byte(in), &NestedField{})
	assert.ErrorContains(t, err, "unexpected node kind")
}

func TestNestedField_Decoding_list(t *testing.T) {
	in := "- foo\n- name: bar\n  pointer: true"
	expected := NestedField{
		Field{Name: "foo"},
		Field{Name: "bar", Pointer: true},
	}

	var actual NestedField
	err := yamlUnmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTagField_Zero(t *testing.T) {
	var f *TagField
	assert.True(t, f.Zero())

	f = &TagField{}
	assert.True(t, f.Zero())

	f.Field = NestedField{Field{}}
	assert.False(t, f.Zero())
}
