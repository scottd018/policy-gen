package files

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	ExtensionJSON = "json"
)

type JSON struct {
	Directory *Directory
	File      string
	Extension string
}

func NewJSONFile(path string, options ...Option) (*JSON, error) {
	// create the directory object
	directory, err := NewDirectory(filepath.Dir(path), options...)
	if err != nil {
		return nil, fmt.Errorf("unable to create directory object for file path [%s]- %w", path, err)
	}

	// validate the file path
	fileParts := strings.Split(filepath.Base(path), ".")
	if len(fileParts) != 1 {
		return nil, fmt.Errorf("invalid file path: [%s]", path)
	}

	return &JSON{
		Directory: directory,
		File:      fileParts[0],
		Extension: ExtensionJSON,
	}, nil
}

func (file *JSON) Path() string {
	return Path(file.Directory, file.File, file.Extension)
}

func (file *JSON) Write(object any, permissions fs.FileMode, options ...Option) error {
	force := hasOption(WithOverwrite, options...)

	// convert struct to json
	jsonData, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal json for file [%s] - %w", file.Path(), err)
	}

	// check if the file already exists
	if _, err = os.Stat(file.Path()); os.IsNotExist(err) {
		// write the file
		err = os.WriteFile(file.Path(), jsonData, permissions)
		if err != nil {
			return fmt.Errorf("unable to write file [%s] - %w", file.Path(), err)
		}

		return nil
	}

	// write the file only if force is requested
	if force {
		if err := os.WriteFile(file.Path(), jsonData, permissions); err != nil {
			return fmt.Errorf("unable to write file [%s] - %w", file.Path(), err)
		}

		return nil
	}

	return fmt.Errorf("cannot write file [%s] - file already exists", file.Path())
}
