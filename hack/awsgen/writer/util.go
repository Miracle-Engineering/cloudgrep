package writer

import "io"

const FileNamePrefix = "zz_"
const FileNameSuffix = ".go"

func multiWrite(s io.Writer, items ...[]byte) error {
	for _, item := range items {
		if _, err := s.Write(item); err != nil {
			return err
		}
	}

	return nil
}

func fileName(name string) string {
	return FileNamePrefix + name + FileNameSuffix
}
