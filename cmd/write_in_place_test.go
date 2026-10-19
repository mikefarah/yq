package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluateWriteInPlace(t *testing.T) {
	originalWriteInplace := writeInplace
	originalExitStatus := exitStatus
	originalCompletedSuccessfully := completedSuccessfully
	t.Cleanup(func() {
		writeInplace = originalWriteInplace
		exitStatus = originalExitStatus
		completedSuccessfully = originalCompletedSuccessfully
	})

	const originalContent = "name: original\n"
	tests := []struct {
		name            string
		expression      string
		exitStatus      bool
		missingInput    bool
		expectedError   string
		expectedContent string
	}{
		{
			name:            "no matches",
			expression:      "select(false)",
			exitStatus:      true,
			expectedError:   "no matches found",
			expectedContent: originalContent,
		},
		{
			name:            "null result",
			expression:      "null",
			exitStatus:      true,
			expectedError:   "no matches found",
			expectedContent: originalContent,
		},
		{
			name:            "false result",
			expression:      "false",
			exitStatus:      true,
			expectedError:   "no matches found",
			expectedContent: originalContent,
		},
		{
			name:            "evaluation error",
			expression:      `error("evaluation failed")`,
			expectedError:   "evaluation failed",
			expectedContent: originalContent,
		},
		{
			name:            "later input error",
			expression:      `.name = "updated"`,
			missingInput:    true,
			expectedContent: originalContent,
		},
		{
			name:            "successful update",
			expression:      `.name = "updated"`,
			expectedContent: "name: updated\n",
		},
		{
			name:            "successful update with exit status",
			expression:      `.name = "updated"`,
			exitStatus:      true,
			expectedContent: "name: updated\n",
		},
		{
			name:       "empty result without exit status",
			expression: "select(false)",
		},
	}

	for _, command := range []string{"eval", "eval-all"} {
		t.Run(command, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					inputDir := t.TempDir()
					tempDir := t.TempDir()
					filename := filepath.Join(inputDir, "input.yaml")
					if err := os.WriteFile(filename, []byte(originalContent), 0600); err != nil {
						t.Fatal(err)
					}
					for _, key := range []string{"TMPDIR", "TMP", "TEMP"} {
						t.Setenv(key, tempDir)
					}
					if filepath.Clean(os.TempDir()) != filepath.Clean(tempDir) {
						t.Fatalf("Expected temporary files in %q, got %q", tempDir, os.TempDir())
					}

					args := []string{command, "-i"}
					if tt.exitStatus {
						args = append(args, "-e")
					}
					args = append(args, tt.expression, filename)
					missingFilename := filepath.Join(inputDir, "missing.yaml")
					if tt.missingInput {
						args = append(args, missingFilename)
					}
					cmd := New()
					var output bytes.Buffer
					cmd.SetOut(&output)
					cmd.SetErr(&output)
					cmd.SetArgs(args)
					err := cmd.Execute()
					if tt.missingInput {
						var pathError *os.PathError
						if !errors.As(err, &pathError) || pathError.Path != missingFilename || !errors.Is(err, os.ErrNotExist) {
							t.Errorf("Expected missing input error for %q, got %v", missingFilename, err)
						}
					} else if tt.expectedError != "" {
						if err == nil || err.Error() != tt.expectedError {
							t.Errorf("Expected error %q, got %v", tt.expectedError, err)
						}
					} else if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}

					content, err := os.ReadFile(filename)
					if err != nil {
						t.Fatal(err)
					}
					if string(content) != tt.expectedContent {
						t.Errorf("Expected file contents %q, got %q", tt.expectedContent, content)
					}
					entries, err := os.ReadDir(tempDir)
					if err != nil {
						t.Fatal(err)
					}
					if len(entries) != 0 {
						t.Errorf("Expected no temporary files, got %v", entries)
					}
				})
			}
		})
	}
}
