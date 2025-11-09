package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateEvaluateAllCommand(t *testing.T) {
	cmd := createEvaluateAllCommand()

	if cmd == nil {
		t.Fatal("createEvaluateAllCommand returned nil")
	}

	// Test basic command properties
	if cmd.Use != "eval-all [expression] [yaml_file1]..." {
		t.Errorf("Expected Use to be 'eval-all [expression] [yaml_file1]...', got %q", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if cmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}

	// Test aliases
	expectedAliases := []string{"ea"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	for i, expected := range expectedAliases {
		if i >= len(cmd.Aliases) || cmd.Aliases[i] != expected {
			t.Errorf("Expected alias %d to be %q, got %q", i, expected, cmd.Aliases[i])
		}
	}
}

func TestEvaluateAll_NoArgs(t *testing.T) {
	// Create a temporary command
	cmd := createEvaluateAllCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with no arguments and no null input
	nullInput = false
	defer func() { nullInput = false }()

	err := evaluateAll(cmd, []string{})

	// Should not error, but should print usage
	if err != nil {
		t.Errorf("evaluateAll with no args should not error, got: %v", err)
	}

	// Should have printed usage information
	if output.Len() == 0 {
		t.Error("Expected usage information to be printed")
	}
}

func TestEvaluateAll_NullInput(t *testing.T) {
	// Create a temporary command
	cmd := createEvaluateAllCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with null input
	nullInput = true
	defer func() { nullInput = false }()

	err := evaluateAll(cmd, []string{})

	// Should not error when using null input
	if err != nil {
		t.Errorf("evaluateAll with null input should not error, got: %v", err)
	}
}

func TestEvaluateAll_WithSingleFile(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("name: test\nage: 25\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateAllCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with a single file
	err = evaluateAll(cmd, []string{yamlFile})

	// Should not error
	if err != nil {
		t.Errorf("evaluateAll with single file should not error, got: %v", err)
	}

	// Should have some output
	if output.Len() == 0 {
		t.Error("Expected output from evaluateAll with single file")
	}
}

func TestEvaluateAll_WithMultipleFiles(t *testing.T) {
	// Create temporary YAML files
	tempDir := t.TempDir()

	yamlFile1 := filepath.Join(tempDir, "test1.yaml")
	yamlContent1 := []byte("name: test1\nage: 25\n")
	err := os.WriteFile(yamlFile1, yamlContent1, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file 1: %v", err)
	}

	yamlFile2 := filepath.Join(tempDir, "test2.yaml")
	yamlContent2 := []byte("name: test2\nage: 30\n")
	err = os.WriteFile(yamlFile2, yamlContent2, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file 2: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateAllCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with multiple files
	err = evaluateAll(cmd, []string{yamlFile1, yamlFile2})

	// Should not error
	if err != nil {
		t.Errorf("evaluateAll with multiple files should not error, got: %v", err)
	}

	// Should have output
	if output.Len() == 0 {
		t.Error("Expected output from evaluateAll with multiple files")
	}
}

func TestEvaluateAll_WithExpression(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("name: test\nage: 25\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateAllCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with expression
	err = evaluateAll(cmd, []string{".name", yamlFile})

	// Should not error
	if err != nil {
		t.Errorf("evaluateAll with expression should not error, got: %v", err)
	}

	// Should have output
	if output.Len() == 0 {
		t.Error("Expected output from evaluateAll with expression")
	}
}

func TestEvaluateAll_WriteInPlace(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("name: test\nage: 25\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateAllCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Enable write in place
	originalWriteInplace := writeInplace
	writeInplace = true
	defer func() { writeInplace = originalWriteInplace }()

	// Test with write in place
	err = evaluateAll(cmd, []string{".name = \"updated\"", yamlFile})

	// Should not error
	if err != nil {
		t.Errorf("evaluateAll with write in place should not error, got: %v", err)
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

func TestEvaluateAll_ExitStatus(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("name: test\nage: 25\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateAllCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Enable exit status
	originalExitStatus := exitStatus
	exitStatus = true
	defer func() { exitStatus = originalExitStatus }()

	// Test with expression that should find no matches
	err = evaluateAll(cmd, []string{".nonexistent", yamlFile})

	// Should error when no matches found and exit status is enabled
	if err == nil {
		t.Error("Expected error when no matches found and exit status is enabled")
	}
}

func TestEvaluateAll_WithMultipleDocuments(t *testing.T) {
	// Create a temporary YAML file with multiple documents
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("---\nname: doc1\nage: 25\n---\nname: doc2\nage: 30\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateAllCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Test with multiple documents
	err = evaluateAll(cmd, []string{".", yamlFile})

	// Should not error
	if err != nil {
		t.Errorf("evaluateAll with multiple documents should not error, got: %v", err)
	}

	// Should have output
	if output.Len() == 0 {
		t.Error("Expected output from evaluateAll with multiple documents")
	}
}

func TestEvaluateAll_NulSepOutput(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := []byte("name: test\nage: 25\n")
	err := os.WriteFile(yamlFile, yamlContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create a temporary command
	cmd := createEvaluateAllCommand()

	// Set up command to capture output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Enable nul separator output
	originalNulSepOutput := nulSepOutput
	nulSepOutput = true
	defer func() { nulSepOutput = originalNulSepOutput }()

	// Test with nul separator output
	err = evaluateAll(cmd, []string{".name", yamlFile})

	// Should not error
	if err != nil {
		t.Errorf("evaluateAll with nul separator output should not error, got: %v", err)
	}

	// Should have output
	if output.Len() == 0 {
		t.Error("Expected output from evaluateAll with nul separator output")
	}
}
