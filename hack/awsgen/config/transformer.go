package config

import (
	"strings"

	"gopkg.in/yaml.v3"
)

const TransformerTypePlaceholder = "%type"

var _ yaml.Unmarshaler = &Transformer{}

func (t *Transformer) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.MappingNode {
		return t.decodeMappingNode(value)
	}
	if value.Kind == yaml.ScalarNode {
		return t.decodeScalarNode(value)
	}

	return &yaml.TypeError{Errors: []string{
		"unexpected node kind",
	}}
}

func (t *Transformer) decodeScalarNode(value *yaml.Node) error {
	var expr string
	err := value.Decode(&expr)
	if err != nil {
		return err
	}

	*t = Transformer{Expr: expr}
	return nil
}

func (t *Transformer) decodeMappingNode(value *yaml.Node) error {
	// This must match the type def of Transformer
	// (but importantly, this type is not a yaml.Unmarshaler to avoid infinite recursion)
	transformer := struct {
		Name         string `yaml:"name"`
		Expr         string `yaml:"expr"`
		ForceGeneric bool   `yaml:"generic"`
	}{}

	err := value.Decode(&transformer)
	if err != nil {
		return err
	}

	if transformer.Expr == "" {
		return &yaml.TypeError{Errors: []string{
			"missing \"expr\"",
		}}
	}

	*t = transformer
	return nil
}

func (t Transformer) IsGeneric() bool {
	return t.ForceGeneric || strings.Contains(t.Expr, TransformerTypePlaceholder)
}

func (t Transformer) Expression(genericType string) string {
	if t.IsGeneric() {
		return strings.ReplaceAll(t.Expr, TransformerTypePlaceholder, genericType)
	}

	return t.Expr
}
