package policy

import (
	"fmt"

	"github.com/scottd018/policy-gen/internal/pkg/files"
)

// ToFiles generates a set of files mapped to their content based on a given
// set of input markers.
func ToFiles(markers []Marker, generator DocumentGenerator) ([]*files.File, error) {
	// generate a marker map from our given set of markers
	markerMap, err := generator.ToPolicyMarkerMap(markers)
	if err != nil {
		return nil, fmt.Errorf("unable to generate policy marker map - %w", err)
	}

	// create a new policy file for each unique key in the markersByFile map
	policyFiles := []*files.File{}

	for filename, markers := range markerMap {
		document, err := generator.ToDocument(markers)
		if err != nil {
			return nil, fmt.Errorf("unable to create document from markers for path [%s] - %w", filename, err)
		}

		file, err := document.ToFile(filename)
		if err != nil {
			return nil, fmt.Errorf("unable to create file from document for path [%s] - %w", filename, err)
		}

		policyFiles = append(policyFiles, file)
	}

	return policyFiles, nil
}
