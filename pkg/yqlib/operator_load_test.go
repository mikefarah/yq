package yqlib

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
)

var loadScenarios = []expressionScenario{
	{
		skipDoc:     true,
		description: "Load empty file with a comment",
		expression:  `load("../../examples/empty.yaml")`,
		expected: []string{
			"D0, P[], (!!null)::# comment\n",
		},
	},
	{
		skipDoc:     true,
		description: "Load and splat",
		expression:  `load("../../examples/small.yaml")[]`,
		expected: []string{
			"D0, P[a], (!!str)::cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "Load and traverse",
		expression:  `load("../../examples/small.yaml").a`,
		expected: []string{
			"D0, P[a], (!!str)::cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "Load file with a header comment into an array",
		document:    `- "../../examples/small.yaml"`,
		expression:  `.[] |= load(.)`,
		expected: []string{
			"D0, P[], (!!seq)::- # comment\n  # about things\n  a: cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "Load empty file with no comment",
		expression:  `load("../../examples/empty-no-comment.yaml")`,
		expected: []string{
			"D0, P[], (!!null)::\n",
		},
	},
	{
		skipDoc:     true,
		description: "Load multiple documents",
		expression:  `load("../../examples/multiple_docs_small.yaml")`,
		expected: []string{
			"D0, P[], ()::- a: Easy! as one two three\n- another:\n    document: here\n- - 1\n  - 2\n",
		},
	},
	{
		description: "Simple example",
		document:    `{myFile: "../../examples/thing.yml"}`,
		expression:  `load(.myFile)`,
		expected: []string{
			"D0, P[], (!!map)::a: apple is included\nb: cool.\n",
		},
	},
	{
		description:    "Replace node with referenced file",
		subdescription: "Note that you can modify the filename in the load operator if needed.",
		document:       `{something: {file: "thing.yml"}}`,
		expression:     `.something |= load("../../examples/" + .file)`,
		expected: []string{
			"D0, P[], (!!map)::{something: {a: apple is included, b: cool.}}\n",
		},
	},
	{
		description:    "Replace _all_ nodes with referenced file",
		subdescription: "Recursively match all the nodes (`..`) and then filter the ones that have a 'file' attribute. ",
		document:       `{something: {file: "thing.yml"}, over: {here: [{file: "thing.yml"}]}}`,
		expression:     `(.. | select(has("file"))) |= load("../../examples/" + .file)`,
		expected: []string{
			"D0, P[], (!!map)::{something: {a: apple is included, b: cool.}, over: {here: [{a: apple is included, b: cool.}]}}\n",
		},
	},
	{
		description:    "Replace node with referenced file as string",
		subdescription: "This will work for any text based file",
		document:       `{something: {file: "thing.yml"}}`,
		expression:     `.something |= load_str("../../examples/" + .file)`,
		expected: []string{
			"D0, P[], (!!map)::{something: \"a: apple is included\\nb: cool.\"}\n",
		},
	},
	{
		requiresFormat: "xml",
		description:    "Load from XML",
		document:       "cool: things",
		expression:     `.more_stuff = load_xml("../../examples/small.xml")`,
		expected: []string{
			"D0, P[], (!!map)::cool: things\nmore_stuff:\n    this: is some xml\n",
		},
	},
	{
		description: "Load from Properties",
		document:    "cool: things",
		expression:  `.more_stuff = load_props("../../examples/small.properties")`,
		expected: []string{
			"D0, P[], (!!map)::cool: things\nmore_stuff:\n    this:\n        is: a properties file\n",
		},
	},
	{
		description:    "Merge from properties",
		subdescription: "This can be used as a convenient way to update a yaml document",
		document:       "this:\n  is: from yaml\n  cool: ay\n",
		expression:     `. *= load_props("../../examples/small.properties")`,
		expected: []string{
			"D0, P[], (!!map)::this:\n    is: a properties file\n    cool: ay\n",
		},
	},
	{
		description: "Load from base64 encoded file",
		document:    "cool: things",
		expression:  `.more_stuff = load_base64("../../examples/base64.txt")`,
		expected: []string{
			"D0, P[], (!!map)::cool: things\nmore_stuff: my secret chilli recipe is....\n",
		},
	},
}

func TestLoadScenarios(t *testing.T) {
	for _, tt := range loadScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "load", loadScenarios)
}

var loadOperatorSecurityDisabledScenarios = []expressionScenario{
	{
		description:    "load() operation fails when security is enabled",
		subdescription: "Use `--security-disable-file-ops` to disable file operations for security.",
		expression:     `load("../../examples/thing.yml")`,
		expectedError:  "file operations have been disabled",
	},
	{
		description:    "load_str() operation fails when security is enabled",
		subdescription: "Use `--security-disable-file-ops` to disable file operations for security.",
		expression:     `load_str("../../examples/thing.yml")`,
		expectedError:  "file operations have been disabled",
	},
	{
		description:    "load_xml() operation fails when security is enabled",
		subdescription: "Use `--security-disable-file-ops` to disable file operations for security.",
		expression:     `load_xml("../../examples/small.xml")`,
		expectedError:  "file operations have been disabled",
	},
	{
		description:    "load_props() operation fails when security is enabled",
		subdescription: "Use `--security-disable-file-ops` to disable file operations for security.",
		expression:     `load_props("../../examples/small.properties")`,
		expectedError:  "file operations have been disabled",
	},
	{
		description:    "load_base64() operation fails when security is enabled",
		subdescription: "Use `--security-disable-file-ops` to disable file operations for security.",
		expression:     `load_base64("../../examples/base64.txt")`,
		expectedError:  "file operations have been disabled",
	},
}

func TestLoadOperatorSecurityDisabledScenarios(t *testing.T) {
	// Save original security preferences
	originalDisableFileOps := ConfiguredSecurityPreferences.DisableFileOps
	defer func() {
		ConfiguredSecurityPreferences.DisableFileOps = originalDisableFileOps
	}()

	// Test that load operations fail when DisableFileOps is true
	ConfiguredSecurityPreferences.DisableFileOps = true

	for _, tt := range loadOperatorSecurityDisabledScenarios {
		testScenario(t, &tt)
	}
	appendOperatorDocumentScenario(t, "load", loadOperatorSecurityDisabledScenarios)
}

func TestLoadWithDecoderClosesFiles(t *testing.T) {
	InitExpressionParser()

	filenames := make([]string, 40)
	dir := t.TempDir()
	for i := range filenames {
		filename := filepath.Join(dir, fmt.Sprintf("input-%d.yml", i))
		if err := os.WriteFile(filename, []byte("value: test\n"), 0600); err != nil {
			t.Fatal(err)
		}
		filenames[i] = filename
	}

	previousGCPercent := debug.SetGCPercent(-1)
	t.Cleanup(func() { debug.SetGCPercent(previousGCPercent) })
	before := countOpenFileDescriptors()
	for _, filename := range filenames {
		node, err := loadWithDecoder(filename, NewYamlDecoder(ConfiguredYamlPreferences))
		if err != nil {
			t.Fatal(err)
		}
		if node.Kind != MappingNode {
			t.Fatalf("got node kind %v, want mapping", node.Kind)
		}
	}
	after := countOpenFileDescriptors()
	if after > before {
		t.Fatalf("open file descriptors grew from %d to %d", before, after)
	}
}

type failingLoadDecoder struct {
	err error
}

func (decoder failingLoadDecoder) Init(io.Reader) error {
	return nil
}

func (decoder failingLoadDecoder) Decode() (*CandidateNode, error) {
	return nil, decoder.err
}

func TestLoadWithDecoderClosesFilesOnDecodeError(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "input.yml")
	if err := os.WriteFile(filename, []byte("value: test\n"), 0600); err != nil {
		t.Fatal(err)
	}

	previousGCPercent := debug.SetGCPercent(-1)
	t.Cleanup(func() { debug.SetGCPercent(previousGCPercent) })
	decodeErr := errors.New("decode failed")
	before := countOpenFileDescriptors()
	for range 40 {
		_, err := loadWithDecoder(filename, failingLoadDecoder{err: decodeErr})
		if !errors.Is(err, decodeErr) {
			t.Fatalf("loadWithDecoder() error = %v, want %v", err, decodeErr)
		}
	}
	after := countOpenFileDescriptors()
	if after > before {
		t.Fatalf("open file descriptors grew from %d to %d", before, after)
	}
}
