package yqlib

import (
	"os/exec"
	"testing"
)

func findExec(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("skipping: %v not found: %v", name, err)
	}
	return path
}

var systemOperatorDisabledScenarios = []expressionScenario{
	{
		description:    "system operator returns error when disabled",
		subdescription: "Use `--security-enable-system-operator` to enable the system operator.",
		document:       "country: Australia",
		expression:     `.country = system("/usr/bin/echo"; "test")`,
		expectedError:  "system operations are disabled, use --security-enable-system-operator to enable",
	},
}

func TestSystemOperatorDisabledScenarios(t *testing.T) {
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
	echoPath := findExec(t, "echo")
	falsePath := findExec(t, "false")

	originalEnableSystemOps := ConfiguredSecurityPreferences.EnableSystemOps
	defer func() {
		ConfiguredSecurityPreferences.EnableSystemOps = originalEnableSystemOps
	}()

	ConfiguredSecurityPreferences.EnableSystemOps = true

	scenarios := []expressionScenario{
		{
			description:    "Run a command with an argument",
			subdescription: "Use `--security-enable-system-operator` to enable the system operator.",
			yqFlags:        "--security-enable-system-operator",
			document:       "country: Australia",
			expression:     `.country = system("` + echoPath + `"; "test")`,
			expected: []string{
				"D0, P[], (!!map)::country: test\n",
			},
		},
		{
			description:    "Run a command without arguments",
			subdescription: "Omit the semicolon and args to run the command with no extra arguments.",
			yqFlags:        "--security-enable-system-operator",
			document:       "a: hello",
			expression:     `.a = system("` + echoPath + `")`,
			expected: []string{
				"D0, P[], (!!map)::a: \"\"\n",
			},
		},
		{
			description:    "Run a command with multiple arguments",
			subdescription: "Pass an array of arguments.",
			skipDoc:        true,
			document:       "a: hello",
			expression:     `.a = system("` + echoPath + `"; ["foo", "bar"])`,
			expected: []string{
				"D0, P[], (!!map)::a: foo bar\n",
			},
		},
		{
			description: "Command and args are evaluated per matched node",
			skipDoc:     true,
			document:    "cmd: " + echoPath + "\narg: hello",
			expression:  `.result = system(.cmd; .arg)`,
			expected: []string{
				"D0, P[], (!!map)::cmd: " + echoPath + "\narg: hello\nresult: hello\n",
			},
		},
		{
			description:   "Command failure returns error",
			skipDoc:       true,
			document:      "a: hello",
			expression:    `.a = system("` + falsePath + `")`,
			expectedError: "system command '" + falsePath + "' failed: exit status 1",
		},
		{
			description:   "Null command returns error",
			skipDoc:       true,
			document:      "a: hello",
			expression:    `.a = system(null)`,
			expectedError: "system operator: command must be a string scalar",
		},
		{
			description: "System operator processes multiple matched nodes",
			skipDoc:     true,
			document:    "a: first",
			document2:   "a: second",
			expression:  `.a = system("` + echoPath + `"; "replaced")`,
			expected: []string{
				"D0, P[], (!!map)::a: replaced\n",
				"D0, P[], (!!map)::a: replaced\n",
			},
		},
	}

	for _, tt := range scenarios {
		testScenario(t, &tt)
	}
	appendOperatorDocumentScenario(t, "system-operators", scenarios)
}
