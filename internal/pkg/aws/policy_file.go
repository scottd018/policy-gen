package aws

import (
	"errors"
	"fmt"

	"github.com/scottd018/policy-gen/internal/pkg/files"
	"github.com/scottd018/policy-gen/internal/pkg/policymarkers"
)

var (
	ErrMarkerConvert      = errors.New("unable to convert policymarker.Marker interface to aws marker")
	ErrMarkerNameMismatch = errors.New("found mismatching marker names in same file")
)

const (
	policyFilePermissions = 0600
)

type PolicyFileGenerator struct {
	Directory *files.Directory
}

// GenerateFile generates a file for a given path and set of markers.  The file includes the content
// based on the policy.
func (generator *PolicyFileGenerator) GenerateFile(path string, markers []policymarkers.Marker) (*files.File, error) {
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

	return NewPolicyDocument(awsMarkers...).ToFile(path)
}

// GetDirectory prints the directory path.  It is use to satisfy the policymarkers.FileGenerator interface.
func (generator *PolicyFileGenerator) GetDirectory() *files.Directory {
	return generator.Directory
}
