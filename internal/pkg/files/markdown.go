package files

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	ExtensionMarkdown = "md"
)

type Markdown struct {
	Directory *Directory
	File      string
	Extension string
}

func NewMarkdownFile(path string, options ...Option) (*Markdown, error) {
	// create the directory object
	directory, err := NewDirectory(filepath.Dir(path), options...)
	if err != nil {
		return nil, err
	}

	// validate the file path
	fileParts := strings.Split(filepath.Base(path), ".")
	if len(fileParts) == 0 {
		return nil, fmt.Errorf("invalid file path: [%s]", path)
	} else if len(fileParts) > 2 {
		return nil, fmt.Errorf("invalid file path: [%s]", path)
	}

	file := &Markdown{
		Directory: directory,
		File:      fileParts[0],
		Extension: ExtensionMarkdown,
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
func (file *Markdown) Path() string {
	return Path(file.Directory, file.File, file.Extension)
}

func (file *Markdown) Write(data []byte, permissions fs.FileMode, options ...Option) error {
	return Write(file.Path(), data, permissions, options...)
}
