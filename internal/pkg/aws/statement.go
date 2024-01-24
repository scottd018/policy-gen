package aws

const (
	statementIdRegex = "^[a-zA-Z0-9]{1,64}$"
)

var (
	defaultStatementEffect   = ValidEffectAllow
	defaultStatementResource = "*"
	defaultStatementId       = "Default"
)

type Statement struct {
	SID       string   `json:"SID"`
	Effect    string   `json:"Effect"`
	Action    []string `json:"Action"`
	Resources []string `json:"Resources"`
}

type Statements []Statement

// Find finds a statement by its SID.  It return a nil value if no statements with
// a matching SID is found.
func (statements Statements) Find(statementId string) *Statement {
	// return immediately if there are no statements.
	if len(statements) == 0 {
		return nil
	}

	// return a statement with the given SID.
	for i := range statements {
		if statements[i].SID == statementId {
			return &statements[i]
		}
	}

	return nil
}

// HasAction determines if a particular statement has an action.
func (statement Statement) HasAction(action string) bool {
	if len(statement.Action) == 0 {
		return false
	}

	for i := range statement.Action {
		if statement.Action[i] == action {
			return true
		}
	}

	return false
}

// HasAction determines if a particular statement has a resource.
func (statement Statement) HasResource(resource string) bool {
	if len(statement.Resources) == 0 {
		return false
	}

	for i := range statement.Resources {
		if statement.Resources[i] == resource {
			return true
		}
	}

	return false
}

// HasEffect determines if a particular statement has an effect.  It is a helper
// method to make the calling code a bit more readable and consistent with the
// other Has* methods.
func (statement Statement) HasEffect(effect string) bool {
	return statement.Effect == effect
}

// AppendAction appends an action to an existing statement.
func (statement *Statement) AppendAction(action string) {
	// if the statement actions are missing add them
	if statement.Action == nil {
		statement.Action = []string{action}

		return
	}

	// if the statement has no actions, simply add them
	if len(statement.Action) == 0 {
		statement.Action = []string{action}

		return
	}

	// if the statement already has the action, simply return
	if statement.HasAction(action) {
		return
	}

	statement.Action = append(statement.Action, action)
}

// AppendResource appends an action to an existing statement.
func (statement *Statement) AppendResource(resource string) {
	// if the statement resources are missing add them
	if statement.Resources == nil {
		statement.Resources = []string{resource}

		return
	}

	// if the statement has no resources, simply add them
	if len(statement.Resources) == 0 {
		statement.Resources = []string{resource}

		return
	}

	// if the statement has exactly one resource equal to the default resource
	// replace it.
	if len(statement.Resources) == 1 && statement.HasResource(defaultStatementResource) {
		statement.Resources = []string{resource}
	}

	// if the statement does not already have the resource, add it
	if !statement.HasResource(resource) {
		statement.Resources = append(statement.Resources, resource)
	}
}

// AppendFor appends a statement to an existing statement given a marker.
func (statement *Statement) AppendFor(marker Marker) {
	// append the action
	statement.AppendAction(*marker.Action)

	// append the resource
	statement.AppendResource(*marker.Resource)
}
