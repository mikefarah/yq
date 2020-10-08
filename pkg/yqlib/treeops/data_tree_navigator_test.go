package treeops

import (
	"strings"
	"testing"

	"github.com/mikefarah/yq/v3/test"
	yaml "gopkg.in/yaml.v3"
)

var treeNavigator = NewDataTreeNavigator(NavigationPrefs{})
var treeCreator = NewPathTreeCreator()

func readDoc(t *testing.T, content string) []*CandidateNode {
	decoder := yaml.NewDecoder(strings.NewReader(content))
	var dataBucket yaml.Node
	err := decoder.Decode(&dataBucket)
	if err != nil {
		t.Error(err)
	}
	return []*CandidateNode{&CandidateNode{Node: &dataBucket, Document: 0}}
}

func resultsToString(results []*CandidateNode) string {
	var pretty string = ""
	for _, n := range results {
		pretty = pretty + "\n" + NodeToString(n)
	}
	return pretty
}

func TestDataTreeNavigatorSimple(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple`)

	path, errPath := treeCreator.ParsePath("a")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a]
  Tag: !!map, Kind: MappingNode, Anchor: 
  b: apple
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSimpleDeep(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple`)

	path, errPath := treeCreator.ParsePath("a.b")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a b]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  apple
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSimpleMismatch(t *testing.T) {

	nodes := readDoc(t, `a: 
  c: apple`)

	path, errPath := treeCreator.ParsePath("a.b")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := ``

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorWild(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: apple`)

	path, errPath := treeCreator.ParsePath("a.*a*")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  apple
`

	test.AssertResult(t, expected, resultsToString(results))
}
