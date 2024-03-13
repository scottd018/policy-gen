package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/scottd018/go-utils/pkg/directory"
)

type Directory struct {
	Path string
}

func NewDirectory(path string, options ...Option) (*Directory, error) {
	// trim the trailing / if we have one
	directoryPath := path
	if directoryPath != "/" {
		directoryPath = strings.TrimSuffix(path, "/")
		directoryPath = filepath.Clean(directoryPath)
	}

	// if we have requested pre-existing directory validation, check to ensure
	// it is valid
	if hasOption(WithPreExistingDirectory, options...) {
		exists, err := directory.Exists(directoryPath)
		if err != nil {
			return nil, fmt.Errorf("invalid directory path [%s] - %w", path, err)
		}

		if !exists {
			return nil, fmt.Errorf("missing directory path [%s]", path)
		}
	}

	return &Directory{
		Path: directoryPath,
	}, nil
}

// ListFilePaths lists file paths within a directory.
func (dir *Directory) ListFilePaths(recursive bool) ([]string, error) {
	paths := []string{}

	// if recursive was requested walk the file path and get any regular files
	if recursive {
		err := filepath.Walk(dir.Path, func(path string, info os.FileInfo, err error) error {
			// return any errors
			if err != nil {
				return err
			}

			// skip any non-regular files
			if !info.Mode().IsRegular() {
				return nil
			}

			paths = append(paths, path)

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("unable to recursively collect file paths - %w", err)
		}

		return paths, nil
	}

	// read only the files within this flat directory structure
	files, err := os.ReadDir(dir.Path)
	if err != nil {
		return paths, fmt.Errorf("unable to list files for directory [%s] - %w", dir.Path, err)
	}

	for path := range files {
		if files[path].Type().IsRegular() {
			paths = append(paths, files[path].Name())
		}
	}

	return paths, nil
}
