package aws

const (
	defaultVersion = "2012=10-17"
)

// PolicyDocument represents an individual AWS IAM policy document.
type PolicyDocument struct {
	Version    string
	Statements Statements
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
	}

	// append the marker data to the existing statement
	statement.AppendFor(marker)
}
