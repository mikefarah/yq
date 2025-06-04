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
	`---
# Anchors and Aliases
defaults: &defaults
  adapter: postgres
  host: localhost

development:
  <<: *defaults
  database: myapp_dev

test:
  <<: *defaults
  database: myapp_test
`,
	`---
# Multiple documents
doc: 1
key: value1
---
doc: 2
key: value2
`,
	`---
# Flow style and tags
flow_map: {item1: val1, item2: val2}
flow_seq: [alpha, bravo, charlie]
custom_tag: !myTag data
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
			v3Nodes := parseYamlDocument(t, yamlDoc)

			// Test with goccy parser
			ConfiguredYamlPreferences.UseGoccyParser = true
			goccyNodes := parseYamlDocument(t, yamlDoc)

			if len(v3Nodes) != len(goccyNodes) {
				t.Errorf("Number of parsed documents differ between yaml.v3 (%d) and goccy (%d) for document %d", len(v3Nodes), len(goccyNodes), i)
				// Log nodes for debugging if lengths differ significantly or one is empty
				if len(v3Nodes) < 5 && len(goccyNodes) < 5 { // Avoid excessive logging for large diffs
					for idx, node := range v3Nodes {
						t.Logf("yaml.v3 result[%d]: %s", idx, NodeToString(node))
					}
					for idx, node := range goccyNodes {
						t.Logf("goccy result[%d]: %s", idx, NodeToString(node))
					}
				}
				return // No point comparing nodes if doc counts differ
			}

			// Compare the parsed structures (not text output)
			for j := range v3Nodes {
				v3Node := v3Nodes[j]
				goccyNode := goccyNodes[j]
				if !recursiveNodeEqual(v3Node, goccyNode) {
					t.Errorf("Parsed structures differ between yaml.v3 and goccy for document %d, sub-document %d", i, j)
					t.Logf("yaml.v3 result: %s", NodeToString(v3Node))
					t.Logf("goccy result: %s", NodeToString(goccyNode))
				}
			}
		})
	}
}

func parseYamlDocument(t *testing.T, yamlDoc string) []*CandidateNode {
	decoder := YamlFormat.DecoderFactory()

	inputs, err := readDocuments(strings.NewReader(yamlDoc), "test.yml", 0, decoder)
	if err != nil {
		t.Fatalf("Failed to decode YAML: %v", err)
	}

	if inputs.Len() == 0 {
		t.Fatal("No documents found")
	}

	var nodes []*CandidateNode
	for el := inputs.Front(); el != nil; el = el.Next() {
		nodes = append(nodes, el.Value.(*CandidateNode))
	}
	return nodes
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
				v3ReparsedNodes := parseYamlDocument(t, v3Result)
				goccyReparsedNodes := parseYamlDocument(t, goccyResult)

				if len(v3ReparsedNodes) != len(goccyReparsedNodes) {
					t.Errorf("Number of reparsed documents differ between yaml.v3 (%d) and goccy (%d) for document %d", len(v3ReparsedNodes), len(goccyReparsedNodes), i)
					// Log nodes for debugging if lengths differ significantly
					if len(v3ReparsedNodes) < 5 && len(goccyReparsedNodes) < 5 { // Avoid excessive logging
						for idx, node := range v3ReparsedNodes {
							t.Logf("yaml.v3 reparsed result[%d]: %s", idx, NodeToString(node))
						}
						for idx, node := range goccyReparsedNodes {
							t.Logf("goccy reparsed result[%d]: %s", idx, NodeToString(node))
						}
					}
				} else {
					for j := range v3ReparsedNodes {
						v3ReparseNode := v3ReparsedNodes[j]
						goccyReparseNode := goccyReparsedNodes[j]
						if !recursiveNodeEqualWithMergeHandling(v3ReparseNode, goccyReparseNode) {
							t.Errorf("Reparsed structures differ for document %d, sub-document %d - this indicates a semantic issue", i, j)
							// Force logging of nodes involved in the direct failing comparison
							t.Logf("Failed comparison details:")
							t.Logf("--- yaml.v3 reparsed node ---\n%s\nNode Structure:\n%s", NodeToString(v3ReparseNode), GetNodeStructure(v3ReparseNode, 0))
							t.Logf("--- goccy reparsed node ---\n%s\nNode Structure:\n%s", NodeToString(goccyReparseNode), GetNodeStructure(goccyReparseNode, 0))
						}
					}
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

	var allDocsOutput strings.Builder
	for i, el := 0, inputs.Front(); el != nil; i, el = i+1, el.Next() {
		node := el.Value.(*CandidateNode)
		var singleDocOutput bytes.Buffer
		err = encoder.Encode(&singleDocOutput, node)
		if err != nil {
			t.Fatalf("Failed to encode YAML document %d: %v", i, err)
		}
		if i > 0 {
			allDocsOutput.WriteString("\n---\n") // Add separator for subsequent documents
		}
		allDocsOutput.Write(singleDocOutput.Bytes())
	}

	return allDocsOutput.String()
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

// Add a helper function to get detailed node structure for logging
func GetNodeStructure(node *CandidateNode, depth int) string {
	if node == nil {
		return strings.Repeat("  ", depth) + "<nil>\n"
	}
	indent := strings.Repeat("  ", depth)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%sKind: %s, Tag: '%s', Value: '%s', Path: %s\n", indent, KindString(node.Kind), node.Tag, node.Value, node.GetNicePath()))
	if node.Alias != nil {
		sb.WriteString(fmt.Sprintf("%s  AliasToTag: '%s', AliasToValue: '%s'\n", indent, node.Alias.Tag, node.Alias.Value))
		if node.Alias.Alias != nil {
			sb.WriteString(fmt.Sprintf("%s    PointsToTag: '%s', PointsToValue: '%s'\n", indent, node.Alias.Alias.Tag, node.Alias.Alias.Value))
		}
	}
	if len(node.Content) > 0 {
		sb.WriteString(fmt.Sprintf("%s  Content (%d items):\n", indent, len(node.Content)))
		for i, child := range node.Content {
			sb.WriteString(fmt.Sprintf("%s    [%d]:\n", indent, i))
			sb.WriteString(GetNodeStructure(child, depth+2))
		}
	}
	return sb.String()
}

// recursiveNodeEqualWithMergeHandling compares two nodes but handles the case where
// one parser uses alias nodes for merge tags while the other expands them directly
func recursiveNodeEqualWithMergeHandling(lhs *CandidateNode, rhs *CandidateNode) bool {
	return recursiveNodeEqualWithMergeHandlingHelper(lhs, rhs, lhs, rhs)
}

// recursiveNodeEqualWithMergeHandlingHelper includes root nodes for alias resolution
func recursiveNodeEqualWithMergeHandlingHelper(lhs *CandidateNode, rhs *CandidateNode, lhsRoot *CandidateNode, rhsRoot *CandidateNode) bool {
	if recursiveNodeEqual(lhs, rhs) {
		return true
	}

	// Special handling for merge tag differences between parsers
	if lhs != nil && rhs != nil && lhs.Kind == MappingNode && rhs.Kind == MappingNode {
		return compareMappingNodesWithMergeHandling(lhs, rhs, lhsRoot, rhsRoot)
	}

	// Handle sequence nodes
	if lhs != nil && rhs != nil && lhs.Kind == SequenceNode && rhs.Kind == SequenceNode {
		if len(lhs.Content) != len(rhs.Content) {
			return false
		}
		for i := range lhs.Content {
			if !recursiveNodeEqualWithMergeHandlingHelper(lhs.Content[i], rhs.Content[i], lhsRoot, rhsRoot) {
				return false
			}
		}
		return true
	}

	return false
}

// compareMappingNodesWithMergeHandling compares mapping nodes and handles merge tag differences
func compareMappingNodesWithMergeHandling(lhs *CandidateNode, rhs *CandidateNode, lhsRoot *CandidateNode, rhsRoot *CandidateNode) bool {
	// Create maps of key-value pairs for comparison
	lhsMap := make(map[string]*CandidateNode)
	rhsMap := make(map[string]*CandidateNode)

	// Helper function to find alias target in the root mapping
	findAliasTarget := func(rootNode *CandidateNode, aliasName string) *CandidateNode {
		for i := 0; i < len(rootNode.Content); i += 2 {
			if rootNode.Content[i].Value == aliasName && rootNode.Content[i+1].Kind == MappingNode {
				return rootNode.Content[i+1]
			}
		}
		return nil
	}

	// Parse lhs content
	for i := 0; i < len(lhs.Content); i += 2 {
		key := lhs.Content[i].Value
		value := lhs.Content[i+1]

		if key == "<<" && value.Kind == AliasNode {
			// This is a merge tag with alias - find the target node
			aliasNode := findAliasTarget(lhsRoot, value.Value)
			if aliasNode != nil {
				// Add each key-value pair from the alias
				for k := 0; k < len(aliasNode.Content); k += 2 {
					aliasKey := aliasNode.Content[k].Value
					aliasValue := aliasNode.Content[k+1]
					if _, exists := lhsMap[aliasKey]; !exists {
						lhsMap[aliasKey] = aliasValue
					}
				}
			}
		} else if key == "<<" && value.Kind == MappingNode {
			// This is a merge tag with expanded content
			for j := 0; j < len(value.Content); j += 2 {
				mergeKey := value.Content[j].Value
				mergeValue := value.Content[j+1]
				if _, exists := lhsMap[mergeKey]; !exists {
					lhsMap[mergeKey] = mergeValue
				}
			}
		} else {
			lhsMap[key] = value
		}
	}

	// Parse rhs content
	for i := 0; i < len(rhs.Content); i += 2 {
		key := rhs.Content[i].Value
		value := rhs.Content[i+1]

		if key == "<<" && value.Kind == AliasNode {
			// This is a merge tag with alias - find the target node
			aliasNode := findAliasTarget(rhsRoot, value.Value)
			if aliasNode != nil {
				// Add each key-value pair from the alias
				for k := 0; k < len(aliasNode.Content); k += 2 {
					aliasKey := aliasNode.Content[k].Value
					aliasValue := aliasNode.Content[k+1]
					if _, exists := rhsMap[aliasKey]; !exists {
						rhsMap[aliasKey] = aliasValue
					}
				}
			}
		} else if key == "<<" && value.Kind == MappingNode {
			// This is a merge tag with expanded content
			for j := 0; j < len(value.Content); j += 2 {
				mergeKey := value.Content[j].Value
				mergeValue := value.Content[j+1]
				if _, exists := rhsMap[mergeKey]; !exists {
					rhsMap[mergeKey] = mergeValue
				}
			}
		} else {
			rhsMap[key] = value
		}
	}

	// Compare the flattened maps
	if len(lhsMap) != len(rhsMap) {
		return false
	}

	for key, lhsValue := range lhsMap {
		rhsValue, exists := rhsMap[key]
		if !exists {
			return false
		}
		if !recursiveNodeEqualWithMergeHandlingHelper(lhsValue, rhsValue, lhsRoot, rhsRoot) {
			return false
		}
	}

	return true
}
