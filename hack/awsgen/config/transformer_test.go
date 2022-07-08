package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestTransformer_Decoding_string(t *testing.T) {
	var actual Transformer
	in := `foo`
	expected := Transformer{Expr: "foo"}

	err := yaml.Unmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTransformer_Decoding_mappingExprOnly(t *testing.T) {
	var actual Transformer
	in := `{expr: bar}`
	expected := Transformer{Expr: "bar"}

	err := yaml.Unmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTransformer_Decoding_mapping(t *testing.T) {
	var actual Transformer
	in := `{name: foo, expr: bar}`
	expected := Transformer{Name: "foo", Expr: "bar"}

	err := yaml.Unmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTransformer_Decoding_mappingInvalid(t *testing.T) {
	in := `{name: [1, 2], expr: bar}`

	err := yaml.Unmarshal([]byte(in), &Field{})
	assert.ErrorContains(t, err, "cannot unmarshal !!seq into string")
}
