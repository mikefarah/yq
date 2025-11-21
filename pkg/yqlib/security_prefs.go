package yqlib

type SecurityPreferences struct {
	DisableEnvOps  bool
	DisableFileOps bool
}

var ConfiguredSecurityPreferences = SecurityPreferences{
	DisableEnvOps:  false,
	DisableFileOps: false,
}
