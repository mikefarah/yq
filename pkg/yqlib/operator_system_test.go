package yqlib

import (
	"testing"
)

var systemOperatorDisabledScenarios = []expressionScenario{
	{
		description:    "system operator returns null when disabled",
		subdescription: "Use `--enable-system-operator` to enable the system operator.",
		document:       "country: Australia",
		expression:     `.country = system("/usr/bin/echo"; "test")`,
		expected: []string{
			"D0, P[], (!!map)::country: null\n",
		},
	},
}

var systemOperatorEnabledScenarios = []expressionScenario{
	{
		description:    "Run a command with an argument",
		subdescription: "Use `--enable-system-operator` to enable the system operator.",
		document:       "country: Australia",
		expression:     `.country = system("/usr/bin/echo"; "test")`,
		expected: []string{
			"D0, P[], (!!map)::country: test\n",
		},
	},
	{
		description:    "Run a command without arguments",
		subdescription: "Omit the semicolon and args to run the command with no extra arguments.",
		document:       "a: hello",
		expression:     `.a = system("/bin/echo")`,
		expected: []string{
			"D0, P[], (!!map)::a: \"\"\n",
		},
	},
	{
		description:    "Run a command with multiple arguments",
		subdescription: "Pass an array of arguments.",
		skipDoc:        true,
		document:       "a: hello",
		expression:     `.a = system("/bin/echo"; ["foo", "bar"])`,
		expected: []string{
			"D0, P[], (!!map)::a: foo bar\n",
		},
	},
	{
		description:   "Command failure returns error",
		skipDoc:       true,
		document:      "a: hello",
		expression:    `.a = system("/bin/false")`,
		expectedError: "system command '/bin/false' failed: exit status 1",
	},
}

func TestSystemOperatorDisabledScenarios(t *testing.T) {
	// ensure system operator is disabled
	originalEnableSystemOps := ConfiguredSecurityPreferences.EnableSystemOps
	defer func() {
		ConfiguredSecurityPreferences.EnableSystemOps = originalEnableSystemOps
	}()

	ConfiguredSecurityPreferences.EnableSystemOps = false

	for _, tt := range systemOperatorDisabledScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "system-operators", systemOperatorDisabledScenarios)
}

func TestSystemOperatorEnabledScenarios(t *testing.T) {
	originalEnableSystemOps := ConfiguredSecurityPreferences.EnableSystemOps
	defer func() {
		ConfiguredSecurityPreferences.EnableSystemOps = originalEnableSystemOps
	}()

	ConfiguredSecurityPreferences.EnableSystemOps = true

	for _, tt := range systemOperatorEnabledScenarios {
		testScenario(t, &tt)
	}
	appendOperatorDocumentScenario(t, "system-operators", systemOperatorEnabledScenarios)
}
