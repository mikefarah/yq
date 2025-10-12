package cmd

import (
	"strings"
	"testing"
)

func TestNewRuneVar(t *testing.T) {
	var r rune
	runeVar := newRuneVar(&r)

	if runeVar == nil {
		t.Fatal("newRuneVar returned nil")
	}
}

func TestRuneValue_String(t *testing.T) {
	tests := []struct {
		name     string
		runeVal  rune
		expected string
	}{
		{
			name:     "simple character",
			runeVal:  'a',
			expected: "a",
		},
		{
			name:     "special character",
			runeVal:  '\n',
			expected: "\n",
		},
		{
			name:     "unicode character",
			runeVal:  'ñ',
			expected: "ñ",
		},
		{
			name:     "zero rune",
			runeVal:  0,
			expected: string(rune(0)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runeVal := runeValue(tt.runeVal)
			result := runeVal.String()
			if result != tt.expected {
				t.Errorf("runeValue.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRuneValue_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    rune
		expectError bool
	}{
		{
			name:        "simple character",
			input:       "a",
			expected:    'a',
			expectError: false,
		},
		{
			name:        "newline escape",
			input:       "\\n",
			expected:    '\n',
			expectError: false,
		},
		{
			name:        "tab escape",
			input:       "\\t",
			expected:    '\t',
			expectError: false,
		},
		{
			name:        "carriage return escape",
			input:       "\\r",
			expected:    '\r',
			expectError: false,
		},
		{
			name:        "form feed escape",
			input:       "\\f",
			expected:    '\f',
			expectError: false,
		},
		{
			name:        "vertical tab escape",
			input:       "\\v",
			expected:    '\v',
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    0,
			expectError: true,
		},
		{
			name:        "multiple characters",
			input:       "ab",
			expected:    0,
			expectError: true,
		},
		{
			name:        "special character",
			input:       "ñ",
			expected:    'ñ',
			expectError: true, // This will fail because the Set function checks len(val) != 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rune
			runeVal := newRuneVar(&r)

			err := runeVal.Set(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if r != tt.expected {
					t.Errorf("Expected rune %q (%d), got %q (%d)",
						string(tt.expected), tt.expected, string(r), r)
				}
			}
		})
	}
}

func TestRuneValue_Set_ErrorMessages(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "empty string error",
			input:         "",
			expectedError: "[] is not a valid character. Must be length 1 was 0",
		},
		{
			name:          "multiple characters error",
			input:         "abc",
			expectedError: "[abc] is not a valid character. Must be length 1 was 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rune
			runeVal := newRuneVar(&r)

			err := runeVal.Set(tt.input)

			if err == nil {
				t.Errorf("Expected error for input %q, but got none", tt.input)
				return
			}

			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error message to contain %q, got %q",
					tt.expectedError, err.Error())
			}
		})
	}
}

func TestRuneValue_Type(t *testing.T) {
	var r rune
	runeVal := newRuneVar(&r)

	result := runeVal.Type()
	expected := "char"

	if result != expected {
		t.Errorf("runeValue.Type() = %q, want %q", result, expected)
	}
}

func TestNew(t *testing.T) {
	rootCmd := New()

	if rootCmd == nil {
		t.Fatal("New() returned nil")
	}

	// Test basic command properties
	if rootCmd.Use != "yq" {
		t.Errorf("Expected Use to be 'yq', got %q", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if rootCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}

	// Test that the command has the expected subcommands
	expectedCommands := []string{"eval", "eval-all", "completion"}
	actualCommands := make([]string, 0, len(rootCmd.Commands()))

	for _, cmd := range rootCmd.Commands() {
		actualCommands = append(actualCommands, cmd.Name())
	}

	for _, expected := range expectedCommands {
		found := false
		for _, actual := range actualCommands {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command %q not found in actual commands: %v",
				expected, actualCommands)
		}
	}
}

func TestNew_FlagCompletions(t *testing.T) {
	rootCmd := New()

	// Test that flag completion functions are registered
	// This is a basic smoke test - we can't easily test the actual completion logic
	// without more complex setup
	flags := []string{
		"output-format",
		"input-format",
		"xml-attribute-prefix",
		"xml-content-name",
		"xml-proc-inst-prefix",
		"xml-directive-name",
		"lua-prefix",
		"lua-suffix",
		"properties-separator",
		"indent",
		"front-matter",
		"expression",
		"split-exp",
	}

	for _, flagName := range flags {
		flag := rootCmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag %q to exist", flagName)
		}
	}
}
