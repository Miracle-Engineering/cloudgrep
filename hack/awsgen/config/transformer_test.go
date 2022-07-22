package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformer_Decoding_string(t *testing.T) {
	var actual Transformer
	in := `foo`
	expected := Transformer{Expr: "foo"}

	err := yamlUnmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTransformer_Decoding_mappingExprOnly(t *testing.T) {
	var actual Transformer
	in := `{expr: bar}`
	expected := Transformer{Expr: "bar"}

	err := yamlUnmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTransformer_Decoding_mapping(t *testing.T) {
	var actual Transformer
	in := `{name: foo, expr: bar}`
	expected := Transformer{Name: "foo", Expr: "bar"}

	err := yamlUnmarshal([]byte(in), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestTransformer_Decoding_mappingInvalid(t *testing.T) {
	in := `{name: [1, 2], expr: bar}`

	err := yamlUnmarshal([]byte(in), &Transformer{})
	assert.ErrorContains(t, err, "cannot unmarshal !!seq into string")
}

func TestTransformer_Decoding_mappingMissingExpr(t *testing.T) {
	in := `{name: foo}`

	err := yamlUnmarshal([]byte(in), &Transformer{})
	assert.ErrorContains(t, err, "missing \"expr\"")
}

func TestTransformer_Decoding_invalidKind(t *testing.T) {
	in := `[1]`

	err := yamlUnmarshal([]byte(in), &Transformer{})
	assert.ErrorContains(t, err, "unexpected node kind")
}

func TestTransformer_Expression_concrete(t *testing.T) {
	tr := Transformer{Expr: "foo[Bar]"}

	assert.Equal(t, "foo[Bar]", tr.Expression("Spam"))
}

func TestTransformer_Expression_generic(t *testing.T) {
	tr := Transformer{Expr: "foo[%type]"}

	assert.Equal(t, "foo[Spam]", tr.Expression("Spam"))
}
