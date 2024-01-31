package input

import "github.com/scottd018/policy-gen/internal/pkg/files"

// Processor represents the processed flag input that has been pre-validated and is
// ready to be passed into a marker processor.
type Processor struct {
	InputDirectory    *files.Directory
	OutputDirectory   *files.Directory
	DocumentationFile *files.Markdown
	Force             bool
	Debug             bool
}
