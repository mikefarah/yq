package yqlib

type ShellVariablesPreferences struct {
	KeySeparator string
	UnwrapScalar bool
}

func NewDefaultShellVariablesPreferences() ShellVariablesPreferences {
	return ShellVariablesPreferences{
		KeySeparator: "_",
		UnwrapScalar: false,
	}
}

var ConfiguredShellVariablesPreferences = NewDefaultShellVariablesPreferences()
