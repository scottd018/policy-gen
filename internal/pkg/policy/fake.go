package policy

const (
	FakeString           = "fake"
	FakeDefinition       = FakeString
	FakeName             = FakeString
	FakeEffectColumn     = "Allow"
	FakePermissionColumn = "*"
	FakeReasonColumn     = FakeString
	FakeResourceColumn   = "*"
)

// fake is a struct to fulfill the policymarkers.Marker interface, but is used
// only for testing.
type fake struct{}

// NewFakeMarker generates a new instances of a fake marker.
//
//nolint:revive
func NewFakeMarker() *fake {
	return &fake{}
}

// fake methods for policies.
func (f *fake) Definition() string { return FakeDefinition }
func (f *fake) Validate() error    { return nil }
func (f *fake) GetName() string    { return FakeName }
func (f *fake) WithDefault()       {}

// fake methods for documentation.
func (f *fake) EffectColumn() string     { return FakeEffectColumn }
func (f *fake) PermissionColumn() string { return FakePermissionColumn }
func (f *fake) ReasonColumn() string     { return FakeReasonColumn }
func (f *fake) ResourceColumn() string   { return FakeResourceColumn }
