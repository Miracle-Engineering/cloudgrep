package config

import (
	"fmt"
	"strings"
)

func (c *Config) Validate() []error {
	var errs []error

	for _, service := range c.Services {
		svcErrs := service.Validate()
		errs = append(errs, svcErrs...)
	}

	return errs
}

func AggregateValidationErrors(errs []error) error {
	var count int // Handle nil errs

	buf := strings.Builder{}
	for _, err := range errs {
		count++

		buf.WriteString("- ")
		buf.WriteString(err.Error())
		buf.WriteString("\n")
	}

	if count == 0 {
		return nil
	}

	plural := "error"
	if count > 1 {
		plural = "errors"
	}

	return fmt.Errorf("encountered %d validation %s:\n%s", count, plural, buf.String())
}
