package aws

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/nukleros/markers"
	"github.com/nukleros/markers/parser"

	"github.com/scottd018/policy-gen/internal/pkg/constants"
	"github.com/scottd018/policy-gen/internal/pkg/input"
)

var (
	ErrMarkerMissingName        = errors.New("marker missing name field")
	ErrMarkerMissingAction      = errors.New("marker missing action field")
	ErrMarkerInvalidEffect      = errors.New("invalid marker effect")
	ErrMarkerInvalidStatementId = errors.New("invalid statement id - must contain a-z, A-Z, 0-9 and limited to 64 characters")
)

const (
	awsMarkerDefinition = "aws:iam:policy"

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

// MarkerDefinition returns the marker definition for an AWS IAM policy marker.
func MarkerDefinition() string {
	return fmt.Sprintf("+%s:%s", constants.MarkerPrefix, awsMarkerDefinition)
}

// MarkerResults return a set of marker results from a given path.
func MarkerResults(path string) ([]*parser.Result, error) {
	policyMarker := Marker{}

	// create a registry for our field markers
	registry := markers.NewRegistry()

	// define our marker
	definition, err := markers.Define(MarkerDefinition(), policyMarker)
	if err != nil {
		return nil, fmt.Errorf("unable to create policy definition for marker [%s] - %w", MarkerDefinition(), err)
	}

	// add the marker to the registry
	registry.Add(definition)

	// collect the data from the given path
	data, err := input.Collect(path)
	if err != nil {
		return nil, fmt.Errorf("error collecting file data for marker [%s] - %w", err)
	}

	// run the parser
	results := markers.NewParser(string(data), registry).Parse()
	if len(results) == 0 {
		fmt.Printf("no results found for marker [%s] at path [%s]\n", MarkerDefinition(), path)

		return []*parser.Result{}, nil
	}

	return results, nil
}

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

	// ensure the sid is valid if specified
	if marker.Id != nil {
		regex := regexp.MustCompile(statementIdRegex)

		if !regex.MatchString(*marker.Id) {
			return fmt.Errorf("%w - [%s]", ErrMarkerInvalidStatementId, *marker.Id)
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
		marker.Id = &defaultStatementId
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

	return policyFiles
}
