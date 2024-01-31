package aws

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/scottd018/policy-gen/internal/pkg/constants"
)

var (
	ErrMarkerMissingName        = errors.New("marker missing name field")
	ErrMarkerMissingAction      = errors.New("marker missing action field")
	ErrMarkerInvalidEffect      = errors.New("invalid marker effect")
	ErrMarkerInvalidStatementID = errors.New("invalid statement id - must contain a-z, A-Z, 0-9 and limited to 64 characters")
)

const (
	awsMarkerDefinition = "aws:iam:policy"

	ValidEffectAllow = "Allow"
	ValidEffectDeny  = "Deny"
)

// we must not lint Id for ID here as the markers package incorrectly parses a
// capitalized ID.
//
//nolint:revive,stylecheck
type Marker struct {
	Name     *string
	Id       *string
	Action   *string
	Effect   *string
	Resource *string
}

type Markers []Marker

// MarkerDefinition returns the marker definition for an AWS IAM policy marker.
func MarkerDefinition() string {
	return fmt.Sprintf("+%s:%s", constants.MarkerPrefix, awsMarkerDefinition)
}

// Validate validates that a marker is valid.
func (marker Marker) Validate() error {
	// ensure required markers are set
	for err, markerField := range map[error]*string{
		ErrMarkerMissingName:   marker.Name,
		ErrMarkerMissingAction: marker.Action,
	} {
		if markerField == nil {
			return err
		}

		if *markerField == "" {
			return err
		}
	}

	// ensure the sid is valid if specified
	if marker.Id != nil {
		regex := regexp.MustCompile(statementIDRegex)

		if !regex.MatchString(*marker.Id) {
			return fmt.Errorf("%w - [%s]", ErrMarkerInvalidStatementID, *marker.Id)
		}
	}

	// ensure effect is valid
	if marker.Effect == nil {
		return nil
	}

	if *marker.Effect != ValidEffectAllow && *marker.Effect != ValidEffectDeny {
		return fmt.Errorf("%w [%s]", ErrMarkerInvalidEffect, *marker.Effect)
	}

	return nil
}

// WithDefault sets a marker with its default values.
func (marker *Marker) WithDefault() {
	// add the effect if we specified one otherwise default to allow
	if marker.Effect == nil {
		marker.Effect = &defaultStatementEffect
	}

	// add the resource if we specified one otherwise default to all
	if marker.Resource == nil {
		marker.Resource = &defaultStatementResource
	}

	// add the id if we specified one otherwise use the default statement id
	if marker.Id == nil {
		marker.Id = &defaultStatementID
	}
}

// ToStatement converts a marker to an AWS IAM policy statement.
func (marker Marker) ToStatement() Statement {
	return Statement{
		Action:    []string{*marker.Action},
		Effect:    *marker.Effect,
		Resources: []string{*marker.Resource},
		SID:       *marker.Id,
	}
}

// PolicyFiles processes a set of markers into their output policy files.
func (m Markers) PolicyFiles() PolicyFiles {
	// markersByFile collects all of the markers that belong to a particular file.
	markersByFile := map[string][]Marker{}

	// collect all of the markers that belong to a particular file and then store
	// them in the markersByFile map.
	for _, marker := range m {
		// ensure default values for the marker
		marker.WithDefault()

		// if the map is nil, add the marker to the array
		if markersByFile[*marker.Name] == nil {
			markersByFile[*marker.Name] = []Marker{marker}

			continue
		}

		// if the array is flat, this marker as the first in the array
		if len(markersByFile[*marker.Name]) == 0 {
			markersByFile[*marker.Name] = []Marker{marker}

			continue
		}

		// append the marker to the current list of markers
		markersByFile[*marker.Name] = append(markersByFile[*marker.Name], marker)
	}

	// create a new policy file for each unique key in the markersByFile map
	policyFiles := PolicyFiles{}

	for filename, markers := range markersByFile {
		policyFiles[filename] = NewPolicyDocument(markers...)
	}

	return policyFiles
}
