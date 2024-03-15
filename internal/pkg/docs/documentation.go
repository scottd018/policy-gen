package docs

import (
	"bytes"
	"errors"

	"github.com/olekukonko/tablewriter"

	"github.com/scottd018/policy-gen/internal/pkg/files"
)

var (
	ErrMissingFilename = errors.New("missing documentation file name")
)

const (
	DocumentationMarkdownHeader  = "# Policy Justification\n\nThis file contains justification for access policies needed by this project.\n\n"
	DocumentationFilePermissions = 0600
)

type Documentation struct {
	File *files.File
}

func NewDocumentation(file *files.File) *Documentation {
	file.Content = append(file.Content, []byte(DocumentationMarkdownHeader)...)

	return &Documentation{File: file}
}

// Generate generates a document from a set of rows.
func (docs *Documentation) Generate(rows ...Row) {
	// create the table
	tableBytes := &bytes.Buffer{}

	table := tablewriter.NewWriter(tableBytes)
	table.SetHeader(Header())
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
	docs.File.Content = append(docs.File.Content, tableBytes.Bytes()...)
}
