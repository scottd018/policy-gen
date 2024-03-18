package policy

import "github.com/scottd018/policy-gen/internal/pkg/files"

// DocumentGenerator is a generic interface which represents an object that generates
// files from a set of policy markers.
type DocumentGenerator interface {
	ToDocument([]Marker) (Document, error)
	ToPolicyMarkerMap(markers []Marker) (MarkerMap, error)
}

// Document is a generic interface which represents a policy document.
type Document interface {
	ToFile(string) (*files.File, error)
}
