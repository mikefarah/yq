//go:build windows

package yqlib

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Regression test for the ACL-loss bug: replacing an existing file via renameFile
// must preserve explicit (non-inherited) ACEs on the replaced file, the way an
// atomic in-place edit is expected to. Before the fix, this went through os.Rename
// (MoveFileEx), which does not guarantee that.
func TestRenameFilePreservesACL(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.yaml")
	replacement := filepath.Join(dir, "replacement.yaml")

	if err := os.WriteFile(target, []byte("a: 1\n"), 0644); err != nil {
		t.Fatalf("failed to create target file: %v", err)
	}
	if err := os.WriteFile(replacement, []byte("a: 2\n"), 0644); err != nil {
		t.Fatalf("failed to create replacement file: %v", err)
	}

	// Grant an explicit (non-inherited) ACE that would not otherwise exist on a
	// fresh temp file, so its survival is a meaningful signal.
	if out, err := exec.Command("icacls", target, "/grant", "Users:(F)").CombinedOutput(); err != nil {
		t.Fatalf("icacls grant failed: %v\n%s", err, out)
	}

	before, err := exec.Command("icacls", target).CombinedOutput()
	if err != nil {
		t.Fatalf("icacls (before) failed: %v\n%s", err, before)
	}
	if !strings.Contains(string(before), `BUILTIN\Users:(F)`) {
		t.Fatalf("test setup did not actually grant the expected ACE, got:\n%s", before)
	}

	if err := renameFile(replacement, target); err != nil {
		t.Fatalf("renameFile failed: %v", err)
	}

	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read target after renameFile: %v", err)
	}
	if string(content) != "a: 2\n" {
		t.Fatalf("target file content wrong after renameFile: got %q", content)
	}

	after, err := exec.Command("icacls", target).CombinedOutput()
	if err != nil {
		t.Fatalf("icacls (after) failed: %v\n%s", err, after)
	}
	if !strings.Contains(string(after), `BUILTIN\Users:(F)`) {
		t.Fatalf("renameFile did not preserve the explicit ACE on the replaced file.\nbefore:\n%s\nafter:\n%s", before, after)
	}
}
