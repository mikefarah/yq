package yqlib

type ShellVariablesPreferences struct {
	KeySeparator string
}

func NewDefaultShellVariablesPreferences() ShellVariablesPreferences {
	return ShellVariablesPreferences{
		KeySeparator: "_",
	}
}

var ConfiguredShellVariablesPreferences = NewDefaultShellVariablesPreferences()

