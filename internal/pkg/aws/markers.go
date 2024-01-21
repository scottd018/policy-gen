package aws

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/nukleros/markers/parser"
)

var (
	ErrMarkerMissingName        = errors.New("marker missing name field")
	ErrMarkerMissingAction      = errors.New("marker missing action field")
	ErrMarkerInvalidEffect      = errors.New("invalid marker effect")
	ErrMarkerInvalidStatementId = errors.New("invalid statement id (SID)")
)

const (
	ValidEffectAllow = "Allow"
	ValidEffectDeny  = "Deny"
)

type Marker struct {
	Name     *string
	Id       *string
	Action   *string
	Effect   *string
	Resource *string
}

type Markers []Marker

// FindMarkers finds all the markers in a given set of parsed results.
func FindMarkers(results []*parser.Result) (Markers, error) {
	markers := make(Markers, len(results))

	for i := range results {
		// ensure the marker we found is the appropriate object type
		marker, ok := results[i].Object.(Marker)
		if !ok {
			return nil, fmt.Errorf(
				"found invalid marker with text [%s] at position [%d]",
				results[i].MarkerText,
				i,
			)
		}

		// ensure the marker we found is valid
		if err := marker.Validate(); err != nil {
			return nil, fmt.Errorf(
				"found invalid marker with text [%s] at position [%d] - %w",
				results[i].MarkerText,
				i,
				err,
			)
		}

		// add the markers to the slice
		markers[i] = marker
	}

	return markers, nil
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

	// ensure effect is valid
	if marker.Effect == nil {
		return nil
	}

	if *marker.Effect != ValidEffectAllow && *marker.Effect != ValidEffectDeny {
		return fmt.Errorf("[%s] - %w", *marker.Effect, ErrMarkerInvalidEffect)
	}

	// ensure the sid is valid if specified
	if marker.Id != nil {
		regex := regexp.MustCompile(statementIdRegex)

		if !regex.MatchString(*marker.Id) {
			return fmt.Errorf("invalid statement id [%s] - %w", *marker.Id, ErrMarkerInvalidStatementId)
		}
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
		marker.Id = &defaultStatementId
	}
}

// ToStatement converts a marker to an AWS IAM policy statement.
func (marker Marker) ToStatement() Statement {
	// create a statement with the base fields that contains all
	// required marker fields
	statement := Statement{
		Action: []string{*marker.Action},
	}

	// add the effect if we specified one otherwise default to allow
	if marker.Effect == nil {
		statement.Effect = *marker.Effect
	} else {
		statement.Effect = defaultStatementEffect
	}

	// add the resource if we specified one otherwise default to all
	if marker.Resource == nil {
		statement.Resources = []string{defaultStatementResource}
	} else {
		statement.Resources = []string{*marker.Resource}
	}

	// add the id if we specified one otherwise use the default statement id
	if marker.Id == nil {
		statement.SID = defaultStatementId
	}

	return statement
}

// Process processes a set of markers into their output policy files.
func (markers Markers) Process() PolicyFiles {
	// markersByFile collects all of the markers that belong to a particular file.
	markersByFile := map[string][]Marker{}

	// collect all of the markers that belong to a particular file and then store
	// them in the markersByFile map.
	for _, marker := range markers {
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

	return nil
}
