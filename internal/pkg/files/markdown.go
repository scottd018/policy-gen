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

	return &Markdown{
		Directory: directory,
		File:      fileParts[0],
		Extension: ExtensionMarkdown,
	}, nil
}

func (file *Markdown) Path() string {
	return Path(file.Directory, file.File, file.Extension)
}

func (file *Markdown) Write(data []byte, permissions fs.FileMode, options ...Option) error {
	force := hasOption(WithOverwrite, options...)

	// check if the file already exists
	if _, err := os.Stat(file.Path()); os.IsNotExist(err) {
		// create the file
		markdownFile, err := os.Create(file.Path())
		if err != nil {
			return fmt.Errorf("unable to create file [%s] - %w", file.Path(), err)
		}
		defer markdownFile.Close()

		// write the file
		if err := os.WriteFile(file.Path(), data, permissions); err != nil {
			return fmt.Errorf("unable to write file [%s] - %w", file.Path(), err)
		}

		return nil
	}

	// write the file only if force is requested
	if force {
		if err := os.WriteFile(file.Path(), data, permissions); err != nil {
			return fmt.Errorf("unable to write file [%s] - %w", file.Path(), err)
		}

		return nil
	}

	return fmt.Errorf("cannot write file [%s] - file already exists", file.Path())
}
