package files

type Option int

const (
	WithPreExistingDirectory Option = iota
	WithOverwrite
)

// hasOption checks to see if an option exists in a set of given options.
func hasOption(option Option, options ...Option) bool {
	if len(options) == 0 {
		return false
	}

	for i := range options {
		if options[i] == option {
			return true
		}
	}

	return false
}
