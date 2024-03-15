package policy

import (
	"fmt"

	"github.com/scottd018/policy-gen/internal/pkg/files"
)

const (
	MarkerPrefix = "policy-gen"
)

// Marker is a generic interface which represents a marker within a file.
type Marker interface {
	// for policies
	Definition() string
	Validate() error
	GetName() string
	WithDefault()

	// for documentation
	EffectColumn() string
	PermissionColumn() string
	ReasonColumn() string
	ResourceColumn() string
}

// MarkerMap is a map of a string to a set of markers.  In this case the string represents
// a file name where the markers will be used to generate content in a file.
type MarkerMap map[string][]Marker

// ToPolicyFiles generates a set of files mapped to their content based on a given
// set of input markers.
func ToPolicyFiles(markers []Marker, generator FileGenerator) ([]*files.File, error) {
	// markersByFile collects all of the markers that belong to a particular file.
	markersByFile := map[string][]Marker{}

	// collect all of the markers that belong to a particular file and then store
	// them in the markersByFile map.
	for _, marker := range markers {
		// ensure default values for the marker
		marker.WithDefault()

		// generate a full file path path as the unique key for our markersByFile map
		path := files.PolicyFilePath(generator.GetDirectory(), marker.GetName())

		// if the map is nil, add the marker to the array
		if markersByFile[path] == nil {
			markersByFile[path] = []Marker{marker}

			continue
		}

		// if the array is flat, this marker as the first in the array
		if len(markersByFile[path]) == 0 {
			markersByFile[path] = []Marker{marker}

			continue
		}

		// append the marker to the current list of markers
		markersByFile[path] = append(markersByFile[path], marker)
	}

	// create a new policy file for each unique key in the markersByFile map
	policyFiles := []*files.File{}

	for filename, markers := range markersByFile {
		file, err := generator.GenerateFile(filename, markers)
		if err != nil {
			return nil, fmt.Errorf("unable to generate file from markers for path [%s] - %w", filename, err)
		}

		policyFiles = append(policyFiles, file)
	}

	return policyFiles, nil
}
