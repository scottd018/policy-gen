package aws

import (
	"encoding/json"
	"fmt"

	"github.com/scottd018/policy-gen/internal/pkg/files"
)

const (
	defaultVersion = "2012-10-17"
)

// PolicyDocument represents an individual AWS IAM policy document.
type PolicyDocument struct {
	Version    string     `json:"Version"`
	Statements Statements `json:"Statement"`
}

// NewPolicyDocument creates a new policy document from a set of markers.
func NewPolicyDocument(markers ...Marker) *PolicyDocument {
	document := &PolicyDocument{Version: defaultVersion}

	for i := range markers {
		document.AddStatementFor(markers[i])
	}

	return document
}

// AddStatementFor takes in a marker input, converts it to a statement, and adds it
// to an existing policy document.
func (document *PolicyDocument) AddStatementFor(marker Marker) {
	// if we do not have any statements set, add this statement as the first
	if len(document.Statements) == 0 {
		document.Statements = []Statement{marker.ToStatement()}

		return
	}

	// find the statement with the id
	statement := document.Statements.Find(*marker.Id)

	// if we do not have a matching statement with an id in our document, create a
	// new statement in the list of existing statements.
	if statement == nil {
		document.Statements = append(document.Statements, marker.ToStatement())

		return
	} else if !statement.HasResource(*marker.Resource) || !statement.HasEffect(*marker.Effect) || !statement.HasCondition(marker.Condition()) {
		marker.AdjustID()
		document.AddStatementFor(marker)

		return
	}

	// append the marker data to the existing statement
	statement.AppendFor(marker)
}

// ToFile converts a policy document to a files.File object reference.
func (document *PolicyDocument) ToFile(path string) (*files.File, error) {
	// we do not need to pass the pre-existing directory option here because it
	// was validated on input
	file, err := files.NewFile(path)
	if err != nil {
		return nil, fmt.Errorf("error converting document to file: [%s] - %w", path, err)
	}

	// convert object to json
	data, err := json.MarshalIndent(document, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json for file: [%s] - %w", path, err)
	}

	// add content to the object
	file.Content = data

	return file, nil
}
