package input

const (
	// input flags
	FlagInputPath     = "input-path"
	FlagOutputPath    = "output-path"
	FlagDocumentation = "documentation"
	FlagRecursive     = "recursive"
	FlagForce         = "force"
	FlagDebug         = "debug"

	// input flag short values
	FlagInputPathShort     = "i"
	FlagOutputPathShort    = "o"
	FlagDocumentationShort = "d"
	FlagRecursiveShort     = "r"
	FlagForceShort         = "f"

	// input flag default values
	FlagInputPathDefault     = "./"
	FlagOutputPathDefault    = "./"
	FlagDocumentationDefault = ""
	FlagRecursiveDefault     = false
	FlagForceDefault         = false
	FlagDebugDefault         = false

	// input flag descriptions
	FlagInputPathDescription     = "Input path to recursively begin parsing markers"
	FlagOutputPathDescription    = "Output path to output generated policies"
	FlagDocumentationDescription = "Documentation file to write"
	FlagRecursiveDescription     = "Recursively find markers from the input-path input"
	FlagForceDescription         = "Forcefully overwrite files with matching names"
	FlagDebugDescription         = "Enable debug logging"
)
