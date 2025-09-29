package fileutil

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ValidExtensions = map[string]bool{
	".mvd": true,
}

func ReadMVD(file string) ([]byte, error) {
	ext := filepath.Ext(file)

	switch ext {
	case ".gz":
		innerExt := filepath.Ext(strings.TrimSuffix(file, ext))
		if !ValidExtensions[innerExt] {
			return nil, fmt.Errorf("unsupported extension %q", innerExt)
		}
		return readGzipFile(file)
	default:
		if !ValidExtensions[ext] {
			return nil, fmt.Errorf("unsupported extension %q", ext)
		}
		return os.ReadFile(file)
	}
}

func readGzipFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	return io.ReadAll(gzr)
}
