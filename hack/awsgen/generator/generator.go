package generator

import (
	"fmt"
	"go/format"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/writer"
)

type Generator struct {
	Format      bool
	LineNumbers bool
}

func (g *Generator) Generate(w writer.Writer, cfg config.Config) error {
	for _, service := range cfg.Services {
		text := g.generateService(service)
		err := g.writeFile(w, service.Name, text)
		if err != nil {
			return fmt.Errorf("cannot generate %s: %w", service.Name, err)
		}
	}

	text := g.generateRegisterFile(cfg.Services)
	err := g.writeFile(w, "register", text)
	if err != nil {
		return fmt.Errorf("cannot generate registration file: %w", err)
	}

	return nil
}

func (g Generator) writeFile(w writer.Writer, name, text string) error {
	if g.Format {
		formatted, err := format.Source([]byte(text))
		if err != nil {
			return fmt.Errorf("cannot format text: %w", err)
		}

		text = string(formatted)
	}

	if g.LineNumbers {
		text = linenumbers(text)
	}

	err := w.WriteFile(name, []byte(text))
	if err != nil {
		return fmt.Errorf("cannot write %s: %w", name, err)
	}

	return nil
}
