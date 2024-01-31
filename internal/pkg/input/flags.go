package input

const (
	FlagInputPath  = "input-path"
	FlagOutputPath = "output-path"
	FlagForce      = "force"
	FlagDebug      = "debug"

	FlagInputPathShort  = "i"
	FlagOutputPathShort = "o"
	FlagForceShort      = "f"

	FlagInputPathDefault  = "./"
	FlagOutputPathDefault = "./"
	FlagForceDefault      = false
	FlagDebugDefault      = false

	FlagInputPathDescription  = "Input path to recursively begin parsing markers"
	FlagOutputPathDescription = "Output path to output generated policies"
	FlagForceDescription      = "Forcefully overwrite files with matching names"
	FlagDebugDescription      = "Enable debug logging"
)
