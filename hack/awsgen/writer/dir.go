package writer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/exp/slices"
)

type dirWriter struct {
	path    string
	written []string
}

func NewDirWriter(dir string) (Writer, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read dir %s: %w", dir, err)
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", dir)
	}

	return &dirWriter{path: dir}, nil
}

func (w *dirWriter) WriteFile(name string, contents []byte) error {
	filePath := w.pathFor(name)

	fmt.Printf("Writing %s\n", filePath)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("cannot open %s for writing: %w", filePath, err)
	}

	defer f.Close()
	w.written = append(w.written, name)

	size, err := f.Write(contents)
	if err != nil {
		return fmt.Errorf("cannot write to %s: %w", filePath, err)
	}

	err = f.Truncate(int64(size))
	if err != nil {
		return fmt.Errorf("cannot truncate %s to %d: %w", filePath, size, err)
	}

	return nil
}

func (w *dirWriter) Clean() error {
	var err error

	dir, err := ioutil.ReadDir(w.path)
	if err != nil {
		return fmt.Errorf("cannot list contents of %s: %w", w.path, err)
	}

	var removeErrors error
	for _, info := range dir {
		if info.IsDir() {
			continue
		}

		name := info.Name()
		filePath := path.Join(w.path, name)

		if !strings.HasPrefix(name, FileNamePrefix) {
			continue
		}

		if !strings.HasSuffix(name, FileNameSuffix) {
			continue
		}

		name = strings.TrimPrefix(name, FileNamePrefix)
		name = strings.TrimSuffix(name, FileNameSuffix)

		if slices.Contains(w.written, name) {
			continue
		}

		fmt.Printf("Removing %s\n", filePath)
		err = os.Remove(filePath)
		if err != nil {
			removeErrors = multierror.Append(err, fmt.Errorf("cannot remove %s: %w", filePath, err))
		}
	}

	w.written = nil

	return removeErrors
}

func (w *dirWriter) pathFor(name string) string {
	return path.Join(w.path, fileName(name))
}
