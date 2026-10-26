package yqlib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMultiPrinterWriter_ClobbersByDefault(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "out.yml")

	if err := os.WriteFile(name, []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to seed existing file: %v", err)
	}

	InitExpressionParser()
	exp, err := ExpressionParser.ParseExpression(`"` + name + `"`)
	if err != nil {
		t.Fatalf("failed to parse expression: %v", err)
	}

	writer := NewMultiPrinterWriter(exp, YamlFormat, false)
	node := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "hi"}

	bufWriter, err := writer.GetWriter(node)
	if err != nil {
		t.Fatalf("expected no error when overwriting is allowed, got: %v", err)
	}
	if _, err := bufWriter.WriteString("new content"); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if err := bufWriter.Flush(); err != nil {
		t.Fatalf("unexpected flush error: %v", err)
	}

	content, readErr := os.ReadFile(name)
	if readErr != nil {
		t.Fatalf("failed to read back file: %v", readErr)
	}
	if string(content) != "new content" {
		t.Fatalf("expected existing file to be overwritten, got: %q", string(content))
	}
}

func TestMultiPrinterWriter_NoClobberErrorsOnExistingFile(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "out.yml")

	if err := os.WriteFile(name, []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to seed existing file: %v", err)
	}

	InitExpressionParser()
	exp, err := ExpressionParser.ParseExpression(`"` + name + `"`)
	if err != nil {
		t.Fatalf("failed to parse expression: %v", err)
	}

	writer := NewMultiPrinterWriter(exp, YamlFormat, true)
	node := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "hi"}

	_, err = writer.GetWriter(node)
	if err == nil {
		t.Fatalf("expected an error when the target file already exists with no-clobber enabled")
	}

	content, readErr := os.ReadFile(name)
	if readErr != nil {
		t.Fatalf("failed to read back file: %v", readErr)
	}
	if string(content) != "existing" {
		t.Fatalf("expected existing file contents to be preserved, got: %q", string(content))
	}
}

func TestMultiPrinterWriter_NoClobberAllowsNewFile(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "brandnew.yml")

	InitExpressionParser()
	exp, err := ExpressionParser.ParseExpression(`"` + name + `"`)
	if err != nil {
		t.Fatalf("failed to parse expression: %v", err)
	}

	writer := NewMultiPrinterWriter(exp, YamlFormat, true)
	node := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "hi"}

	bufWriter, err := writer.GetWriter(node)
	if err != nil {
		t.Fatalf("expected no error creating a brand new file, got: %v", err)
	}
	if err := bufWriter.Flush(); err != nil {
		t.Fatalf("unexpected flush error: %v", err)
	}

	if _, statErr := os.Stat(name); statErr != nil {
		t.Fatalf("expected file to be created: %v", statErr)
	}
}

func TestMultiPrinterWriter_BackwardsCompatibleTwoArgConstructor(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "out.yml")

	InitExpressionParser()
	exp, err := ExpressionParser.ParseExpression(`"` + name + `"`)
	if err != nil {
		t.Fatalf("failed to parse expression: %v", err)
	}

	// callers compiled against the pre-existing two-argument signature must still build and run
	writer := NewMultiPrinterWriter(exp, YamlFormat)
	node := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "hi"}

	bufWriter, err := writer.GetWriter(node)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if err := bufWriter.Flush(); err != nil {
		t.Fatalf("unexpected flush error: %v", err)
	}
}
