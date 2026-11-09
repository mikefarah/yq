package yqlib

import (
	"bufio"
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupSplitWriterTest(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	workDir := filepath.Join(tempDir, "work")
	if err := os.Mkdir(workDir, 0750); err != nil {
		t.Fatal(err)
	}
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(workDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Error(err)
		}
	})

	originalAllowWriteAnyPath := ConfiguredSecurityPreferences.AllowWriteAnyPath
	ConfiguredSecurityPreferences.AllowWriteAnyPath = false
	t.Cleanup(func() {
		ConfiguredSecurityPreferences.AllowWriteAnyPath = originalAllowWriteAnyPath
	})

	return tempDir
}

func getTestSplitWriter(t *testing.T, name string) (*bufio.Writer, error) {
	t.Helper()
	InitExpressionParser()
	expression, err := ExpressionParser.ParseExpression(".")
	if err != nil {
		t.Fatal(err)
	}
	printerWriter := NewMultiPrinterWriter(expression, &Format{})
	return printerWriter.GetWriter(&CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: name})
}

func TestMultiPrinterWriterRejectsPathsOutsideWorkingDirectory(t *testing.T) {
	tempDir := setupSplitWriterTest(t)

	logger := GetLogger()
	originalLogger := logger.slogger
	var logged bytes.Buffer
	logger.SetSlogger(slog.New(slog.NewTextHandler(&logged, nil)))
	t.Cleanup(func() {
		logger.SetSlogger(originalLogger)
	})

	outsidePath := filepath.Join(tempDir, "absolute.txt")
	for _, path := range []string{"../parent.txt", outsidePath} {
		_, err := getTestSplitWriter(t, path)
		if err == nil {
			t.Fatalf("expected path %q to be rejected", path)
		}
		if !strings.Contains(err.Error(), "--security-allow-write-any-path") {
			t.Errorf("error does not explain the override flag: %v", err)
		}
	}

	if !strings.Contains(logged.String(), "--security-allow-write-any-path") {
		t.Errorf("warning does not explain the override flag: %q", logged.String())
	}
	if _, err := os.Stat(outsidePath); !os.IsNotExist(err) {
		t.Errorf("rejected destination was created: %v", err)
	}
}

func TestMultiPrinterWriterAllowsRelativePaths(t *testing.T) {
	setupSplitWriterTest(t)

	writer, err := getTestSplitWriter(t, filepath.Join("nested", "output.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("output")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(filepath.Join("nested", "output.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "output" {
		t.Errorf("unexpected output: %q", output)
	}
}

func TestMultiPrinterWriterAllowsAbsolutePathsWithinWorkingDirectory(t *testing.T) {
	setupSplitWriterTest(t)
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(workDir, "absolute.txt")
	writer, err := getTestSplitWriter(t, path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("output")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "output" {
		t.Errorf("unexpected output: %q", output)
	}
}

func TestMultiPrinterWriterRejectsSymlinkDirectories(t *testing.T) {
	tempDir := setupSplitWriterTest(t)
	targetDir := filepath.Join(tempDir, "target")
	if err := os.Mkdir(targetDir, 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(targetDir, "linked"); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	_, err := getTestSplitWriter(t, filepath.Join("linked", "output.txt"))
	if err == nil || !strings.Contains(err.Error(), "--security-allow-write-any-path") {
		t.Fatalf("expected symlink directory to be rejected with override guidance, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(targetDir, "output.txt")); !os.IsNotExist(err) {
		t.Errorf("rejected destination was created: %v", err)
	}
}

func TestMultiPrinterWriterReplacesFinalSymlink(t *testing.T) {
	tempDir := setupSplitWriterTest(t)
	targetPath := filepath.Join(tempDir, "target.txt")
	if err := os.WriteFile(targetPath, []byte("target"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(targetPath, "output.txt"); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	writer, err := getTestSplitWriter(t, "output.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("output")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}

	info, err := os.Lstat("output.txt")
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("output path is still a symlink")
	}
	target, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(target) != "target" {
		t.Errorf("symlink target was modified: %q", target)
	}
}

func TestMultiPrinterWriterDoesNotTruncateHardLinkTarget(t *testing.T) {
	tempDir := setupSplitWriterTest(t)
	targetPath := filepath.Join(tempDir, "target.txt")
	if err := os.WriteFile(targetPath, []byte("target"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(targetPath, "output.txt"); err != nil {
		t.Skipf("cannot create hard link: %v", err)
	}

	writer, err := getTestSplitWriter(t, "output.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("output")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}

	target, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(target) != "target" {
		t.Errorf("hard-link target was modified: %q", target)
	}
	output, err := os.ReadFile("output.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "output" {
		t.Errorf("unexpected output: %q", output)
	}
}

func TestMultiPrinterWriterAllowsAnyPathWithOverride(t *testing.T) {
	tempDir := setupSplitWriterTest(t)
	ConfiguredSecurityPreferences.AllowWriteAnyPath = true
	path := filepath.Join(tempDir, "outside.txt")

	writer, err := getTestSplitWriter(t, path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("output")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "output" {
		t.Errorf("unexpected output: %q", output)
	}
}
