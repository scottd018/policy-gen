package policy

import (
	"fmt"

	"github.com/scottd018/policy-gen/internal/pkg/files"
)

// File is a generic interface which represents an object that is a policy file.
type File interface {
	FromMarkers([]Marker)
}

// FileGenerator is a generic interface which represents an object that generates
// files from a set of policy markers.
type FileGenerator interface {
	GenerateFile(string, []Marker) (*files.File, error)
	GetDirectory() *files.Directory
	ToPolicyMarkerMap(markers []Marker) (MarkerMap, error)
}

// ToFiles generates a set of files mapped to their content based on a given
// set of input markers.
func ToFiles(markers []Marker, generator FileGenerator) ([]*files.File, error) {
	// generate a marker map from our given set of markers
	markerMap, err := generator.ToPolicyMarkerMap(markers)
	if err != nil {
		return nil, fmt.Errorf("unable to generate policy marker map - %w", err)
	}

	// create a new policy file for each unique key in the markersByFile map
	policyFiles := []*files.File{}

	for filename, markers := range markerMap {
		file, err := generator.GenerateFile(filename, markers)
		if err != nil {
			return nil, fmt.Errorf("unable to generate file from markers for path [%s] - %w", filename, err)
		}

		policyFiles = append(policyFiles, file)
	}

	return policyFiles, nil
}
