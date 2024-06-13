package docs

const (
	HeaderEffect     = "effect"
	HeaderPermission = "permission"
	HeaderResource   = "resource"
	HeaderReason     = "reason"
	HeaderCondition  = "condition"
)

// Header defines the table Header for our documentation page.  This is ordered, so be
// aware that changing the order will affect the display.
func Header() []string {
	return []string{
		HeaderEffect,
		HeaderPermission,
		HeaderResource,
		HeaderReason,
		HeaderCondition,
	}
}
