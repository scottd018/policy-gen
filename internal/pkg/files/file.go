package files

import (
	"fmt"
	"io/fs"
	"os"
)

// Path returns the string value of a file path given a directory, file name
// and an extension.
func Path(dir *Directory, file, extension string) string {
	return fmt.Sprintf("%s/%s.%s", dir.Path, file, extension)
}

// Write writes data to a file.
func Write(path string, data []byte, permissions fs.FileMode, options ...Option) error {
	// write the file if it does not already exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, data, permissions); err != nil {
			return fmt.Errorf("unable to write file [%s] - %w", path, err)
		}

		return nil
	}

	force := hasOption(WithOverwrite, options...)

	// write the file only if force is requested
	if force {
		if err := os.WriteFile(path, data, permissions); err != nil {
			return fmt.Errorf("unable to write file [%s] - %w", path, err)
		}

		return nil
	}

	return fmt.Errorf("cannot write file [%s] - file already exists and force not requested", path)
}
