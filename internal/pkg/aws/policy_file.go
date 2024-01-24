package aws

import (
	"encoding/json"
	"fmt"
	"os"
)

type PolicyFiles map[string]*PolicyDocument

// Write writes the policy files.
func (files PolicyFiles) Write(path string, force bool) error {
	for file, document := range files {
		// convert struct to json
		jsonData, err := json.MarshalIndent(document, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal json for file [%s] - %w", file, err)
		}

		// create a file with the key as the file name
		filename := fmt.Sprintf("%s/%s.json", path, file)

		// check if the file already exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			// write the file
			err = os.WriteFile(filename, jsonData, 0644)
			if err != nil {
				return fmt.Errorf("unable to write file [%s] - %w", filename, err)
			}
		}

		// write the file only if force is requested
		if force {
			err = os.WriteFile(filename, jsonData, 0644)
			if err != nil {
				return fmt.Errorf("unable to write file [%s] - %w", filename, err)
			}

			return nil
		}

		return fmt.Errorf("unable to write file [%s]; use --force if you wish to overwrite", filename)
	}

	return nil
}
