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
	if len(fileParts) == 0 {
		return nil, fmt.Errorf("invalid file path: [%s]", path)
	} else if len(fileParts) > 2 {
		return nil, fmt.Errorf("invalid file path: [%s]", path)
	}

	file := &JSON{
		Directory: directory,
		File:      fileParts[0],
		Extension: ExtensionJSON,
	}

	// ensure we are not pointing at a directory path
	dirInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return file, nil
	}

	// check if it is actually a directory
	if dirInfo.IsDir() {
		return nil, fmt.Errorf("invalid file path; path is directory: [%s]", path)
	}

	return file, nil
}

// Path returns the full path value for a file.
func (file *JSON) Path() string {
	return Path(file.Directory, file.File, file.Extension)
}

// Write writes the JSON data, from an object, to a file.
func (file *JSON) Write(object any, permissions fs.FileMode, options ...Option) error {
	// convert object to json
	jsonData, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshal json for file [%s] - %w", file.Path(), err)
	}

	return Write(file.Path(), jsonData, permissions, options...)
}
