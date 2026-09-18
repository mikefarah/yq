package yqlib

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// helper to build an ExpressionNode that just yields a fixed string for the file name
func parseFilenameExp(t *testing.T, exp string) *ExpressionNode {
	t.Helper()
	InitExpressionParser()
	node, err := ExpressionParser.ParseExpression(exp)
	if err != nil {
		t.Fatalf("failed to parse split-exp test expression %q: %v", exp, err)
	}
	return node
}

func TestMultiPrinterWriterOverwriteDefault(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "out.yml")
	if err := os.WriteFile(target, []byte("pre-existing\n"), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	exp := parseFilenameExp(t, `"`+target+`"`)
	pw := NewMultiPrinterWriter(exp, YamlFormat)

	node := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "hello"}
	w, err := pw.GetWriter(node)
	if err != nil {
		t.Fatalf("default behaviour should silently overwrite, got error: %v", err)
	}
	if w == nil {
		t.Fatalf("expected a writer, got nil")
	}
	// confirm the file was truncated/recreated by os.Create
	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat target: %v", err)
	}
	if info.Size() != 0 {
		t.Fatalf("expected file to be truncated (size 0) before writes, got %d bytes", info.Size())
	}
}

func TestMultiPrinterWriterNoOverwriteRefusesExisting(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "out.yml")
	if err := os.WriteFile(target, []byte("pre-existing\n"), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	exp := parseFilenameExp(t, `"`+target+`"`)
	pw := NewMultiPrinterWriterWithOptions(exp, YamlFormat, true)

	node := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "hello"}
	_, err := pw.GetWriter(node)
	if err == nil {
		t.Fatalf("expected error when --no-overwrite is set and target exists, got nil")
	}
	if !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("expected refusing-to-overwrite error message, got: %v", err)
	}

	// file must be untouched
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(data) != "pre-existing\n" {
		t.Fatalf("file should be untouched, contents = %q", string(data))
	}
}

func TestMultiPrinterWriterNoOverwriteCreatesNew(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "new.yml")

	exp := parseFilenameExp(t, `"`+target+`"`)
	pw := NewMultiPrinterWriterWithOptions(exp, YamlFormat, true)

	node := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "hello"}
	w, err := pw.GetWriter(node)
	if err != nil {
		t.Fatalf("no-overwrite should still create new files, got: %v", err)
	}
	if w == nil {
		t.Fatalf("expected a writer, got nil")
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("expected new file to exist, stat err: %v", err)
	}
}
