package yqlib

type SecurityPreferences struct {
	DisableEnvOps   bool
	DisableFileOps  bool
	EnableSystemOps bool
	// AllowedEnvVars, when non-empty, restricts env()/strenv()/envsubst()
	// to the listed variable names. Empty means no restriction.
	AllowedEnvVars []string
}

// IsEnvAllowed reports whether name may be read by env operations.
// An empty allowlist allows everything (default behavior unchanged).
func (p SecurityPreferences) IsEnvAllowed(name string) bool {
	if len(p.AllowedEnvVars) == 0 {
		return true
	}
	for _, allowed := range p.AllowedEnvVars {
		if allowed == name {
			return true
		}
	}
	return false
}

var ConfiguredSecurityPreferences = SecurityPreferences{
	DisableEnvOps:   false,
	DisableFileOps:  false,
	EnableSystemOps: false,
}
