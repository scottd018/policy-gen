package docs

type Row interface {
	EffectColumn() string
	PermissionColumn() string
	ResourceColumn() string
	ReasonColumn() string
	ConditionColumn() string
}
