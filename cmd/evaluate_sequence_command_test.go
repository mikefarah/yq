package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateEvaluateSequenceCommand(t *testing.T) {
	cmd := createEvaluateSequenceCommand()

	if cmd == nil {
		t.Fatal("createEvaluateSequenceCommand returned nil")
	}

	// Test basic command properties
	if cmd.Use != "eval [expression] [yaml_file1]..." {
		t.Errorf("Expected Use to be 'eval [expression] [yaml_file1]...', got %q", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if cmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}

	// Test aliases
	expectedAliases := []string{"e"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	for i, expected := range expectedAliases {
		if i >= len(cmd.Aliases) || cmd.Aliases[i] != expected {
			t.Errorf("Expected alias %d to be %q, got %q", i, expected, cmd.Aliases[i])
		}
	}
}

func TestProcessExpression(t *testing.T) {
	// Reset global variables
	originalPrettyPrint := prettyPrint
	defer func() { prettyPrint = originalPrettyPrint }()

	tests := []struct {
		name        string
		prettyPrint bool
		expression  string
		expected    string
	}{
		{
			name:        "empty expression without pretty print",
			prettyPrint: false,
			expression:  "",
			expected:    "",
		},
		{
			name:        "empty expression with pretty print",
			prettyPrint: true,
			expression:  "",
			expected:    `(... | (select(tag != "!!str"), select(tag == "!!str") | select(test("(?i)^(y|yes|n|no|on|off)$") | not))  ) style=""`,
		},
		{
			name:        "simple expression without pretty print",
			prettyPrint: false,
			expression:  ".a.b",
			expected:    ".a.b",
		},
		{
			name:        "simple expression with pretty print",
			prettyPrint: true,
			expression:  ".a.b",
			expected:    `.a.b | (... | (select(tag != "!!str"), select(tag == "!!str") | select(test("(?i)^(y|yes|n|no|on|off)$") | not))  ) style=""`,
		},
		{
			name:        "complex expression with pretty print",
			prettyPrint: true,
			expression:  ".items[] | select(.active == true)",
			expected:    `.items[] | select(.active == true) | (... | (select(tag != "!!str"), select(tag == "!!str") | select(test("(?i)^(y|yes|n|no|on|off)$") | not))  ) style=""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prettyPrint = tt.prettyPrint
			result := processExpression(tt.expression)
			if result != tt.expected {
				t.Errorf("processExpression(%q) = %q, want %q", tt.expression, result, tt.expected)
			}
		})
	}
}

func TestEvaluateSequence_NoArgs(t *testing.T) {
	// Create a temporary command
	cmd := createEvaluateSequenceCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with no arguments and no null input
	nullInput = false
	defer func() { nullInput = false }()

	err := evaluateSequence(cmd, []string{})

	// Should not error, but should print usage
	if err != nil {
		t.Errorf("evaluateSequence with no args should not error, got: %v", err)
	}

	// Should have printed usage information
	if output.Len() == 0 {
		t.Error("Expected usage information to be printed")
	}
}

func TestEvaluateSequence_NullInput(t *testing.T) {
	// Create a temporary command
	cmd := createEvaluateSequenceCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with null input
	nullInput = true
	defer func() { nullInput = false }()

	err := evaluateSequence(cmd, []string{})

	// Should not error when using null input
	if err != nil {
		t.Errorf("evaluateSequence with null input should not error, got: %v", err)
	}
}

func TestEvaluateSequence_WithFile(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("name: test\nage: 25\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateSequenceCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with a file
	err = evaluateSequence(cmd, []string{yamlFile})

	// Should not error
	if err != nil {
		t.Errorf("evaluateSequence with file should not error, got: %v", err)
	}

	// Should have some output
	if output.Len() == 0 {
		t.Error("Expected output from evaluateSequence with file")
	}
}

func TestEvaluateSequence_WithExpressionAndFile(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("name: test\nage: 25\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateSequenceCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with expression and file
	err = evaluateSequence(cmd, []string{".name", yamlFile})

	// Should not error
	if err != nil {
		t.Errorf("evaluateSequence with expression and file should not error, got: %v", err)
	}

	// Should have output
	if output.Len() == 0 {
		t.Error("Expected output from evaluateSequence with expression and file")
	}
}

func TestEvaluateSequence_WriteInPlace(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("name: test\nage: 25\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateSequenceCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Enable write in place
	originalWriteInplace := writeInplace
	writeInplace = true
	defer func() { writeInplace = originalWriteInplace }()

	// Test with write in place
	err = evaluateSequence(cmd, []string{".name = \"updated\"", yamlFile})

	// Should not error
	if err != nil {
		t.Errorf("evaluateSequence with write in place should not error, got: %v", err)
	}

	// Verify the file was updated
	updatedContent, err := os.ReadFile(yamlFile)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	// Should contain the updated content
	if !strings.Contains(string(updatedContent), "updated") {
		t.Errorf("Expected file to contain 'updated', got: %s", string(updatedContent))
	}
}

func TestEvaluateSequence_ExitStatus(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("name: test\nage: 25\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateSequenceCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Enable exit status
	originalExitStatus := exitStatus
	exitStatus = true
	defer func() { exitStatus = originalExitStatus }()

	// Test with expression that should find no matches
	err = evaluateSequence(cmd, []string{".nonexistent", yamlFile})

	// Should error when no matches found and exit status is enabled
	if err == nil {
		t.Error("Expected error when no matches found and exit status is enabled")
	}
}
