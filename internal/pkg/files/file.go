package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	ExtensionJSON     = "json"
	ExtensionMarkdown = "md"

	ModePolicyFile   = 0600
	ModeDocumentFile = 0600
)

var (
	ErrFileMissingContent = errors.New("file is missing content")
)

// File represents a file object.
type File struct {
	Directory *Directory
	File      string
	Content   []byte
}

// NewFile creates a new instance of a file object.
func NewFile(path string, options ...Option) (*File, error) {
	// create the directory object
	directory, err := NewDirectory(filepath.Dir(path), options...)
	if err != nil {
		return nil, fmt.Errorf("unable to create directory object for file path [%s]- %w", path, err)
	}

	file := &File{
		Directory: directory,
		File:      path,
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

// PolicyFilePath returns the full path for a policy file given a directory and file key name.
func PolicyFilePath(directory *Directory, key string) string {
	return fmt.Sprintf("%s/%s.%s", directory.Path, key, ExtensionJSON)
}

// DocumentationFilePath returns the full path for a documentation file given a directory and file key name.
func DocumentationFilePath(directory *Directory, key string) string {
	return fmt.Sprintf("%s/%s.%s", directory.Path, key, ExtensionMarkdown)
}

// Path returns the full path value for a file.
func (file *File) Path() string {
	return fmt.Sprintf("%s/%s", file.Directory.Path, file.File)
}

// Write writes data to a file.
func (file *File) Write(permissions fs.FileMode, options ...Option) error {
	// return an error if we have no content to write
	if file.Content == nil || len(file.Content) == 0 {
		return ErrFileMissingContent
	}

	// write the file if it does not already exist
	if _, err := os.Stat(file.File); os.IsNotExist(err) {
		if err := os.WriteFile(file.File, file.Content, permissions); err != nil {
			return fmt.Errorf("unable to write file [%s] - %w", file.File, err)
		}

		return nil
	}

	force := hasOption(WithOverwrite, options...)

	// write the file only if force is requested
	if force {
		if err := os.WriteFile(file.File, file.Content, permissions); err != nil {
			return fmt.Errorf("unable to write file [%s] - %w", file.File, err)
		}

		return nil
	}

	return fmt.Errorf("cannot write file [%s] - file already exists and force not requested", file.File)
}
