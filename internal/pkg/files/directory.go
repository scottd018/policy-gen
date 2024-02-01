package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

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

// CollectData recursively collects all data as an []byte from a given directory.
func (dir *Directory) CollectData() ([]byte, error) {
	var result []byte

	err := filepath.Walk(dir.Path, func(filePath string, info os.FileInfo, err error) error {
		// return any errors
		if err != nil {
			return err
		}

		// skip directories and symlinks
		if info.IsDir() || (info.Mode()&os.ModeSymlink) == os.ModeSymlink {
			return nil
		}

		// read in the file
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("unable to read file [%s] - %w", filePath, err)
		}

		// only append text file content
		if utf8.Valid(fileContent) {
			result = append(result, fileContent...)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to recursively collect file input - %w", err)
	}

	return result, nil
}
