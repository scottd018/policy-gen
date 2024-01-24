package input

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func Collect(path string) ([]byte, error) {
	var result []byte

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
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
