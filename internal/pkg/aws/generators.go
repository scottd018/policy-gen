package aws

import (
	"errors"
	"fmt"

	"github.com/scottd018/policy-gen/internal/pkg/files"
	"github.com/scottd018/policy-gen/internal/pkg/policy"
)

var (
	ErrMarkerConvert      = errors.New("unable to convert policy.Marker interface to aws.Marker object")
	ErrMarkerNameMismatch = errors.New("found mismatching marker names in same file")
)

type PolicyDocumentGenerator struct {
	Directory *files.Directory
}

// ToPolicyMarkerMap generates a map of filenames with their given set of markers.
func (generator *PolicyDocumentGenerator) ToPolicyMarkerMap(markers []policy.Marker) (policy.MarkerMap, error) {
	// markerMap collects all of the markers that belong to a particular file.
	markerMap := policy.MarkerMap{}

	// collect all of the markers that belong to a particular file and then store
	// them in the markersByFile map.
	for _, marker := range markers {
		// ensure we are working with an aws marker
		awsMarker, ok := marker.(*Marker)
		if !ok {
			return nil, ErrMarkerConvert
		}

		// ensure default values for the marker
		awsMarker.WithDefault()

		// generate a full file path path as the unique key for our markersByFile map
		path := files.PolicyFilePath(generator.Directory, *awsMarker.Name)

		// if the map is nil, add the marker to the array
		if markerMap[path] == nil {
			markerMap[path] = []policy.Marker{awsMarker}

			continue
		}

		// if the array is flat, this marker as the first in the array
		if len(markerMap[path]) == 0 {
			markerMap[path] = []policy.Marker{awsMarker}

			continue
		}

		// append the marker to the current list of markers
		markerMap[path] = append(markerMap[path], marker)
	}

	return markerMap, nil
}

// ToDocument generates a document from a given set of markers.  The file includes the content
// based on the policy.
func (generator *PolicyDocumentGenerator) ToDocument(markers []policy.Marker) (policy.Document, error) {
	awsMarkers := make([]Marker, len(markers))

	var name string

	// validate the markers and convert them to the proper type
	for i := range markers {
		marker, ok := markers[i].(*Marker)
		if !ok {
			return nil, ErrMarkerConvert
		}

		if name != "" {
			if name != *marker.Name {
				return nil, fmt.Errorf("[%s/%s] - %w", name, *marker.Name, ErrMarkerNameMismatch)
			}
		} else {
			name = *marker.Name
		}

		awsMarkers[i] = *marker
	}

	return NewPolicyDocument(awsMarkers...), nil
}
