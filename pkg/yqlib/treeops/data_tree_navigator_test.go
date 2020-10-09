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
  cat: apple
  mad: things`)

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

-- Node --
  Document 0, path: [a mad]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  things
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorWildDeepish(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: {b: 3}
  mad: {b: 4}
  fad: {c: t}`)

	path, errPath := treeCreator.ParsePath("a.*a*.b")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat b]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  3

-- Node --
  Document 0, path: [a mad b]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  4
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorOrSimple(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: apple
  mad: things`)

	path, errPath := treeCreator.ParsePath("a.(cat or mad)")
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

-- Node --
  Document 0, path: [a mad]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  things
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorOrSimpleWithDepth(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: {b: 3}
  mad: {b: 4}
  fad: {c: t}`)

	path, errPath := treeCreator.ParsePath("a.(cat.* or fad.*)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat b]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  3

-- Node --
  Document 0, path: [a fad c]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  t
`
	test.AssertResult(t, expected, resultsToString(results))
}
