package writer

// Writer is a mechanism to persist generated Go files
type Writer interface {
	// WriteFile writes a named file, adding a prefix and suffix to the name.
	WriteFile(name string, contents []byte) error

	// Clean removes any generated files from previous executions that WriteFile has not touched.
	Clean() error
}
