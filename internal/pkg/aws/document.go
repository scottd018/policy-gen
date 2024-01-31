package aws

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/scottd018/go-utils/pkg/pointers"
)

const (
	defaultVersion = "2012-10-17"
)

// PolicyDocument represents an individual AWS IAM policy document.
type PolicyDocument struct {
	Version    string     `json:"Version"`
	Statements Statements `json:"Statements"`
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
	} else if !statement.HasEffect(*marker.Effect) {
		// adjust the id if we have effects that are mismatched (e.g. Allow/Deny)
		marker.Id = pointers.String(fmt.Sprintf("%s%s", *marker.Id, *marker.Effect))

		document.AddStatementFor(marker)
	}

	// append the marker data to the existing statement
	statement.AppendFor(marker)
}

// Write writes a policy document to a file.
func (document *PolicyDocument) Write(file string, force bool) error {
	// convert struct to json
	jsonData, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal json for file [%s] - %w", file, err)
	}

	// check if the file already exists
	if _, err = os.Stat(file); os.IsNotExist(err) {
		// write the file
		err = os.WriteFile(file, jsonData, policyFilePermissions)
		if err != nil {
			return fmt.Errorf("unable to write file [%s] - %w", file, err)
		}
	}

	// write the file only if force is requested
	if !force {
		return fmt.Errorf("unable to write file [%s]; use --force if you wish to overwrite", file)
	}

	if err := os.WriteFile(file, jsonData, policyFilePermissions); err != nil {
		return fmt.Errorf("unable to write file [%s] - %w", file, err)
	}

	return nil
}
