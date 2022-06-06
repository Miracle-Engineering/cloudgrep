package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestFieldDecoding_string(t *testing.T) {
	var actual Field
	in := `foo`
	expected := Field{Name: "foo"}

	err := yaml.Unmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFieldDecoding_mapping(t *testing.T) {
	var actual Field
	in := `{name: foo, slice: string}`
	expected := Field{Name: "foo", SliceType: "string"}

	err := yaml.Unmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFieldDecoding_sequence(t *testing.T) {
	var actual Field
	in := "- foo\n- bar"
	err := yaml.Unmarshal([]byte(in), &actual)
	assert.ErrorContains(t, err, "unexpected node kind")

	expected := Field{}

	assert.Equal(t, expected, actual)
}

func TestFieldDecoding_noName(t *testing.T) {
	var actual Field
	in := `{slice: true}`
	expected := Field{}

	err := yaml.Unmarshal([]byte(in), &actual)
	assert.ErrorContains(t, err, "missing \"name\"")
	assert.Equal(t, expected, actual)
}
