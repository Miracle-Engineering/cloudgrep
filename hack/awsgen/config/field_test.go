package config

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestField_Decoding_string(t *testing.T) {
	var actual Field
	in := `foo`
	expected := Field{Name: "foo"}

	err := yaml.Unmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestField_Decoding_mapping(t *testing.T) {
	var actual Field
	in := `{name: foo, sliceType: string}`
	expected := Field{Name: "foo", SliceType: "string"}

	err := yaml.Unmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestField_Decoding_sequence(t *testing.T) {
	var actual Field
	in := "- foo\n- bar"
	err := yaml.Unmarshal([]byte(in), &actual)
	assert.ErrorContains(t, err, "unexpected node kind")

	expected := Field{}

	assert.Equal(t, expected, actual)
}

func TestField_Decoding_noName(t *testing.T) {
	var actual Field
	in := `{sliceType: string}`
	expected := Field{}

	err := yaml.Unmarshal([]byte(in), &actual)
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

func TestTagField_Zero(t *testing.T) {
	var f *TagField
	assert.True(t, f.Zero())

	f = &TagField{}
	assert.True(t, f.Zero())

	f.Field = NestedField{Field{}}
	assert.False(t, f.Zero())
}
