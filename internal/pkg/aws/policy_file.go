package aws

import (
	"fmt"
	"strings"
)

const (
	policyFilePermissions = 0600
)

type PolicyFiles map[string]*PolicyDocument

// Filename constructs a name for a policy file to write.
func Filename(path, name string) string {
	if strings.HasSuffix(path, "/") {
		return fmt.Sprintf("%s%s.json", path, name)
	}

	return fmt.Sprintf("%s/%s.json", path, name)
}
