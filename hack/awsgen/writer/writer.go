package writer

type Writer interface {
	WriteFile(name string, contents []byte) error
	Clean() error
}
