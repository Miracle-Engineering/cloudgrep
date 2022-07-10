package config

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func yamlUnmarshal(in []byte, out any) error {
	r := bytes.NewReader(in)
	d := yaml.NewDecoder(r)
	d.KnownFields(true)

	return d.Decode(out)
}
