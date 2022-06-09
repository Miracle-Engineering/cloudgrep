package writer

// FakeWriter is a writer that is useful during tests
type FakeWriter struct {
	Files map[string]string
}

var _ Writer = &FakeWriter{}

func NewFakeWriter() *FakeWriter {
	w := &FakeWriter{}
	_ = w.Clean()

	return w
}

func (w *FakeWriter) WriteFile(name string, contents []byte) error {
	if w.Files == nil {
		w.Files = make(map[string]string)
	}

	w.Files[name] = string(contents)
	return nil
}

func (w *FakeWriter) Clean() error {
	w.Files = make(map[string]string)
	return nil
}
