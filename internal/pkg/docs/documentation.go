package docs

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/olekukonko/tablewriter"

	"github.com/scottd018/policy-gen/internal/pkg/files"
)

var (
	ErrMissingFilename = errors.New("missing documentation file name")
)

const (
	HeaderEffect     = "effect"
	HeaderPermission = "permission"
	HeaderResource   = "resource"
	HeaderReason     = "reason"

	DocumentationMarkdownHeader  = "# Policy Justification\n\nThis file contains justification for access policies needed by this project.\n\n"
	DocumentationFilePermissions = 0600
)

type Documentation struct {
	File *files.Markdown
	Data []byte
}

func NewDocumentation(file *files.Markdown) *Documentation {
	return &Documentation{
		File: file,
		Data: []byte(DocumentationMarkdownHeader),
	}
}

// Write writes documentation for a Row definition.
func (docs *Documentation) Write(force bool, rows ...Row) error {
	// ensure our documentation has a file name set
	if docs.File.Path() == "" {
		return ErrMissingFilename
	}

	// create the table
	tableBytes := &bytes.Buffer{}

	table := tablewriter.NewWriter(tableBytes)
	table.SetHeader(header())
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	// append the data for each row to the table
	for _, row := range rows {
		table.Append(
			[]string{
				row.EffectColumn(),
				row.PermissionColumn(),
				row.ResourceColumn(),
				row.ReasonColumn(),
			},
		)
	}

	// write the data to the bytes buffer and return
	table.Render()

	// append the rendered data to the existing data
	docs.Data = append(docs.Data, tableBytes.Bytes()...)

	// set the file options
	var options []files.Option
	if force {
		options = []files.Option{files.WithOverwrite}
	}

	// write the file
	if err := docs.File.Write(docs.Data, DocumentationFilePermissions, options...); err != nil {
		return fmt.Errorf("error writing file object [%s] - %w", docs.File.Path(), err)
	}

	return nil
}

// header defines the table header for our documentation page.  This is ordered, so be
// aware that changing the order will affect the display.
func header() []string {
	return []string{
		HeaderEffect,
		HeaderPermission,
		HeaderResource,
		HeaderReason,
	}
}
