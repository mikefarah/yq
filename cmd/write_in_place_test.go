package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
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

func TestEvaluateWriteInPlaceSymlinkToReadOnlyTarget(t *testing.T) {
	originalWriteInplace := writeInplace
	originalExitStatus := exitStatus
	originalCompletedSuccessfully := completedSuccessfully
	t.Cleanup(func() {
		writeInplace = originalWriteInplace
		exitStatus = originalExitStatus
		completedSuccessfully = originalCompletedSuccessfully
	})

	for _, command := range []string{"eval", "eval-all"} {
		t.Run(command, func(t *testing.T) {
			dir := t.TempDir()
			target := filepath.Join(dir, "target.yaml")
			link := filepath.Join(dir, "link.yaml")
			if err := os.WriteFile(target, []byte("name: original\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(target, 0444); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("target.yaml", link); err != nil {
				t.Fatal(err)
			}

			cmd := New()
			cmd.SetArgs([]string{command, "-i", `.name = "updated"`, link})
			if err := cmd.Execute(); err != nil {
				t.Fatalf("in-place update failed: %v", err)
			}

			info, err := os.Lstat(link)
			if err != nil {
				t.Fatal(err)
			}
			if !info.Mode().IsRegular() || info.Mode().Perm() != 0444 {
				t.Errorf("replacement mode = %v, want regular file with mode 0444", info.Mode())
			}
			content, err := os.ReadFile(link)
			if err != nil {
				t.Fatal(err)
			}
			if string(content) != "name: updated\n" {
				t.Errorf("replacement content = %q", content)
			}
			content, err = os.ReadFile(target)
			if err != nil {
				t.Fatal(err)
			}
			if string(content) != "name: original\n" {
				t.Errorf("target content changed: %q", content)
			}
			info, err = os.Stat(target)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != 0444 {
				t.Errorf("target permissions changed: %v", info.Mode().Perm())
			}
		})
	}
}

// TestEvaluateWriteInPlacePreservesFormatting is a regression test for
// https://github.com/mikefarah/yq/issues/465 - writing a file in place must
// preserve the comments, blank lines, indentation and comment spacing of the
// untouched parts of the document.
func TestEvaluateWriteInPlacePreservesFormatting(t *testing.T) {
	originalWriteInplace := writeInplace
	originalPrefs := yqlib.ConfiguredYamlPreferences
	t.Cleanup(func() {
		writeInplace = originalWriteInplace
		yqlib.ConfiguredYamlPreferences = originalPrefs
	})

	tests := []struct {
		name            string
		expression      string
		originalContent string
		expectedContent string
	}{
		{
			// the exact example from issue #465
			name:       "comments, blank lines and inline comment spacing preserved",
			expression: `.logLevel = 4`,
			originalContent: `# Log level for operator
logLevel: 3

# Annotations applied to operator deployment
annotations: {}

# Annotations passed to operator pod(s).
podAnnotations: {}

# Node labels for pod assignment
nodeSelector:  # +doc-gen:break
  beta.kubernetes.io/os: linux
  beta.kubernetes.io/arch: amd64
`,
			expectedContent: `# Log level for operator
logLevel: 4

# Annotations applied to operator deployment
annotations: {}

# Annotations passed to operator pod(s).
podAnnotations: {}

# Node labels for pod assignment
nodeSelector:  # +doc-gen:break
  beta.kubernetes.io/os: linux
  beta.kubernetes.io/arch: amd64
`,
		},
		{
			name:       "custom indentation and quotes preserved",
			expression: `.parent.child1 = "changed"`,
			originalContent: `# top comment

parent:
    # nested comment
    child1: 'single quoted'

    child2: "double quoted"

    list:
        - one

        - two
`,
			expectedContent: `# top comment

parent:
    # nested comment
    child1: 'changed'

    child2: "double quoted"

    list:
        - one

        - two
`,
		},
		{
			name:       "blank lines around literal blocks and foot comments preserved",
			expression: `.after = false`,
			originalContent: `literal: |
  line1
  line2

after: true

foot: 1
# foot comment of foot

last: 2
`,
			expectedContent: `literal: |
  line1
  line2

after: false

foot: 1
# foot comment of foot

last: 2
`,
		},
		{
			name:       "blank lines between sequence entries preserved",
			expression: `.seq[1] = "B"`,
			originalContent: `seq:
  - a

  - b
  - c
last: 1
`,
			expectedContent: `seq:
  - a

  - B
  - c
last: 1
`,
		},
		{
			name:       "deleting a node removes its comment but keeps surrounding blank lines",
			expression: `del(.gone)`,
			originalContent: `keep: 1

# comment on gone
gone: 2

after: 3
`,
			expectedContent: `keep: 1


after: 3
`,
		},
		{
			name:            "crlf line endings preserved",
			expression:      `.a = 10`,
			originalContent: "a: 1\r\n\r\n# comment\r\nb: 2\r\n",
			expectedContent: "a: 10\r\n\r\n# comment\r\nb: 2\r\n",
		},
	}

	for _, command := range []string{"eval", "eval-all"} {
		t.Run(command, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					inputDir := t.TempDir()
					tempDir := t.TempDir()
					filename := filepath.Join(inputDir, "input.yaml")
					if err := os.WriteFile(filename, []byte(tt.originalContent), 0600); err != nil {
						t.Fatal(err)
					}
					for _, key := range []string{"TMPDIR", "TMP", "TEMP"} {
						t.Setenv(key, tempDir)
					}

					cmd := New()
					var output bytes.Buffer
					cmd.SetOut(&output)
					cmd.SetErr(&output)
					cmd.SetArgs([]string{command, "-i", tt.expression, filename})
					if err := cmd.Execute(); err != nil {
						t.Fatalf("unexpected error: %v (output: %s)", err, output.String())
					}

					content, err := os.ReadFile(filename)
					if err != nil {
						t.Fatal(err)
					}
					if string(content) != tt.expectedContent {
						t.Errorf("Expected file contents:\n%q\ngot:\n%q", tt.expectedContent, content)
					}
				})
			}
		})
	}
}

// TestEvaluateWriteInPlaceExplicitIndentWins checks that an explicit --indent
// flag is respected and disables the automatic indentation detection that
// in-place formatting preservation otherwise performs.
func TestEvaluateWriteInPlaceExplicitIndentWins(t *testing.T) {
	originalWriteInplace := writeInplace
	originalPrefs := yqlib.ConfiguredYamlPreferences
	t.Cleanup(func() {
		writeInplace = originalWriteInplace
		yqlib.ConfiguredYamlPreferences = originalPrefs
	})

	dir := t.TempDir()
	filename := filepath.Join(dir, "input.yaml")
	original := "parent:\n    child: 1\n"
	if err := os.WriteFile(filename, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := New()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"eval", "-i", "-I", "8", `.parent.child = 2`, filename})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v (output: %s)", err, output.String())
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expected := "parent:\n        child: 2\n"
	if string(content) != expected {
		t.Errorf("Expected file contents %q, got %q", expected, content)
	}
}
