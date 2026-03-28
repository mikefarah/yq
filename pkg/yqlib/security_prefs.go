package yqlib

type SecurityPreferences struct {
	DisableEnvOps   bool
	DisableFileOps  bool
	EnableSystemOps bool
}

var ConfiguredSecurityPreferences = SecurityPreferences{
	DisableEnvOps:   false,
	DisableFileOps:  false,
	EnableSystemOps: false,
}
