package yqlib

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var testYamlDocs = []string{
	`---
name: John
age: 30
address:
  street: 123 Main St
  city: Anytown
  state: ST
  zip: 12345
hobbies:
  - reading
  - swimming
  - cooking
`,
	`---
# This is a comment
version: "1.0"
database:
  host: localhost
  port: 5432
  name: myapp
  credentials:
    username: admin
    password: secret123
servers:
  - name: web1
    ip: 192.168.1.10
  - name: web2  
    ip: 192.168.1.11
`,
	`---
# Complex YAML with various types
config:
  enabled: true
  timeout: 30
  ratio: 3.14159
  options:
    debug: false
    verbose: true
  items: []
  metadata: null
  multiline: |
    This is a multiline
    string with line breaks
    and proper formatting
`,
}

func TestYamlParsersStructuralEquivalence(t *testing.T) {
	// Save original preferences
	originalPrefs := ConfiguredYamlPreferences.Copy()
	defer func() {
		ConfiguredYamlPreferences = originalPrefs
	}()

	for i, yamlDoc := range testYamlDocs {
		t.Run(fmt.Sprintf("Document_%d", i), func(t *testing.T) {
			// Test with yaml.v3 parser
			ConfiguredYamlPreferences.UseGoccyParser = false
			v3Node := parseYamlDocument(t, yamlDoc)

			// Test with goccy parser
			ConfiguredYamlPreferences.UseGoccyParser = true
			goccyNode := parseYamlDocument(t, yamlDoc)

			// Compare the parsed structures (not text output)
			if !recursiveNodeEqual(v3Node, goccyNode) {
				t.Errorf("Parsed structures differ between yaml.v3 and goccy for document %d", i)
				t.Logf("yaml.v3 result: %s", NodeToString(v3Node))
				t.Logf("goccy result: %s", NodeToString(goccyNode))
			}
		})
	}
}

func parseYamlDocument(t *testing.T, yamlDoc string) *CandidateNode {
	decoder := YamlFormat.DecoderFactory()

	inputs, err := readDocuments(strings.NewReader(yamlDoc), "test.yml", 0, decoder)
	if err != nil {
		t.Fatalf("Failed to decode YAML: %v", err)
	}

	if inputs.Len() == 0 {
		t.Fatal("No documents found")
	}

	return inputs.Front().Value.(*CandidateNode)
}

func TestYamlParsersEquivalence(t *testing.T) {
	// Save original preferences
	originalPrefs := ConfiguredYamlPreferences.Copy()
	defer func() {
		ConfiguredYamlPreferences = originalPrefs
	}()

	// Note: This test documents the differences in output formatting between parsers
	// These differences are acceptable as they don't affect semantic meaning
	t.Log("Note: yaml.v3 and goccy may produce different key ordering, which is semantically equivalent")

	for i, yamlDoc := range testYamlDocs {
		t.Run(fmt.Sprintf("Document_%d", i), func(t *testing.T) {
			// Test with yaml.v3 parser
			ConfiguredYamlPreferences.UseGoccyParser = false
			v3Result := processYamlDocument(t, yamlDoc)

			// Test with goccy parser
			ConfiguredYamlPreferences.UseGoccyParser = true
			goccyResult := processYamlDocument(t, yamlDoc)

			// Log the differences for documentation purposes
			if v3Result != goccyResult {
				t.Logf("Output format differs (this is expected):")
				t.Logf("yaml.v3 output:\n%s", v3Result)
				t.Logf("goccy output:\n%s", goccyResult)

				// But both should be valid YAML that parses to the same structure
				v3Reparse := parseYamlDocument(t, v3Result)
				goccyReparse := parseYamlDocument(t, goccyResult)

				if !recursiveNodeEqual(v3Reparse, goccyReparse) {
					t.Error("Reparsed structures differ - this indicates a semantic issue")
				}
			}
		})
	}
}

func processYamlDocument(t *testing.T, yamlDoc string) string {
	decoder := YamlFormat.DecoderFactory()
	encoder := YamlFormat.EncoderFactory()

	inputs, err := readDocuments(strings.NewReader(yamlDoc), "test.yml", 0, decoder)
	if err != nil {
		t.Fatalf("Failed to decode YAML: %v", err)
	}

	if inputs.Len() == 0 {
		t.Fatal("No documents found")
	}

	node := inputs.Front().Value.(*CandidateNode)

	var output bytes.Buffer
	err = encoder.Encode(&output, node)
	if err != nil {
		t.Fatalf("Failed to encode YAML: %v", err)
	}

	return output.String()
}

func TestGoccyParserBasicFunctionality(t *testing.T) {
	// Save original preferences
	originalPrefs := ConfiguredYamlPreferences.Copy()
	defer func() {
		ConfiguredYamlPreferences = originalPrefs
	}()

	// Enable goccy parser
	ConfiguredYamlPreferences.UseGoccyParser = true

	yamlDoc := `---
test:
  field1: value1
  field2: 42
  field3: true
`

	decoder := YamlFormat.DecoderFactory()
	encoder := YamlFormat.EncoderFactory()

	inputs, err := readDocuments(strings.NewReader(yamlDoc), "test.yml", 0, decoder)
	test.AssertResult(t, nil, err)
	test.AssertResult(t, 1, inputs.Len())

	node := inputs.Front().Value.(*CandidateNode)

	var output bytes.Buffer
	err = encoder.Encode(&output, node)
	test.AssertResult(t, nil, err)

	result := output.String()

	// Basic checks - the exact format might differ but structure should be preserved
	if !strings.Contains(result, "test:") {
		t.Error("Expected 'test:' in output")
	}
	if !strings.Contains(result, "field1: value1") {
		t.Error("Expected 'field1: value1' in output")
	}
	if !strings.Contains(result, "field2: 42") {
		t.Error("Expected 'field2: 42' in output")
	}
	if !strings.Contains(result, "field3: true") {
		t.Error("Expected 'field3: true' in output")
	}
}
