package yqlib

type ShellVariablesPreferences struct {
	KeySeparator string
	UnwrapScalar bool
	// SecurityForceNeverEscape restores the legacy behaviour of emitting unwrapped (-r) values unquoted. Unsafe.
	SecurityForceNeverEscape bool
}

func NewDefaultShellVariablesPreferences() ShellVariablesPreferences {
	return ShellVariablesPreferences{
		KeySeparator: "_",
		UnwrapScalar: false,
	}
}

var ConfiguredShellVariablesPreferences = NewDefaultShellVariablesPreferences()
