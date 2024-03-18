package policy

const (
	MarkerPrefixStart  = "+"
	MarkerPrefixString = "policy-gen"
)

// Marker is a generic interface which represents a marker within a file.
type Marker interface {
	// for policies
	Definition() string
	Validate() error
	GetName() string
	WithDefault()

	// for documentation
	EffectColumn() string
	PermissionColumn() string
	ReasonColumn() string
	ResourceColumn() string
}

// MarkerMap is a map of a string to a set of markers.  In this case the string represents
// a file name where the markers will be used to generate content in a file.
type MarkerMap map[string][]Marker
