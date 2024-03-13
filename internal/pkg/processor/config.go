package processor

import (
	"github.com/scottd018/policy-gen/internal/pkg/files"
)

// Config represents the configuration for a processor.
type Config struct {
	InputDirectory    *files.Directory
	OutputDirectory   *files.Directory
	DocumentationFile *files.File
	Recursive         bool
	Force             bool
	Debug             bool
}
