package generator

import (
	"fmt"
	"go/format"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/writer"
)

// Generator generates the AWS provider code based on the configuration.
type Generator struct {
	// Format configures whether or not the Go code is passed through `go fmt` before being written.
	// Disabling this is useful for diagnosing issues with the generator.
	Format bool

	// LineNumbers causes line number comments to be added to the beginning of each line that is written.
	// It is intended to aid in debugging the generator, as the line number comments may not be completely accurate depending
	// on which generator is being used.
	LineNumbers bool
}

// Generate generates all functions and writes them to the Writer.
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

// writeFile acts as a convienence function for file writing that automatically applies formatting and line numbers as required.
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
