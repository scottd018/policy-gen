package aws

import (
	"fmt"
)

const (
	policyFilePermissions = 0600
)

type PolicyFiles map[string]*PolicyDocument

// PolicyFilename constructs a name for a policy file to write.  This
func PolicyFilename(path, name string) string {
	return fmt.Sprintf("%s/%s", path, name)
}
