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
	return []*CandidateNode{&CandidateNode{Node: dataBucket.Content[0], Document: 0}}
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

func TestDataTreeNavigatorDeleteSimple(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple
  c: camel`)

	path, errPath := treeCreator.ParsePath("a .- b")
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
  c: camel
`
	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorDeleteTwice(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple
  c: camel
  d: dingo`)

	path, errPath := treeCreator.ParsePath("a .- b OR a .- c")
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
  d: dingo
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorDeleteWithUnion(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple
  c: camel
  d: dingo`)

	path, errPath := treeCreator.ParsePath("a .- (b OR c)")
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
  d: dingo
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorDeleteByIndex(t *testing.T) {

	nodes := readDoc(t, `a: 
  - b: apple
  - b: sdfsd
  - b: apple`)

	path, errPath := treeCreator.ParsePath("(a .- (0 or 1))")
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
  Tag: !!seq, Kind: SequenceNode, Anchor: 
  - b: apple
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorDeleteByFind(t *testing.T) {

	nodes := readDoc(t, `a: 
  - b: apple
  - b: sdfsd
  - b: apple`)

	path, errPath := treeCreator.ParsePath("(a .- (* == apple))")
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
  Tag: !!seq, Kind: SequenceNode, Anchor: 
  - b: sdfsd
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorDeleteArrayByFind(t *testing.T) {

	nodes := readDoc(t, `a: 
  - apple
  - sdfsd
  - apple`)

	path, errPath := treeCreator.ParsePath("(a .- (. == apple))")
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
  Tag: !!seq, Kind: SequenceNode, Anchor: 
  - sdfsd
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorDeleteViaSelf(t *testing.T) {

	nodes := readDoc(t, `- apple
- sdfsd
- apple`)

	path, errPath := treeCreator.ParsePath(". .- (. == apple)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: []
  Tag: !!seq, Kind: SequenceNode, Anchor: 
  - sdfsd
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorCountWithFilter(t *testing.T) {

	nodes := readDoc(t, `f:
  a: frog
  b: dally
  c: log`)

	path, errPath := treeCreator.ParsePath("f(count(. == *og))")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [f]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  2
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorCountWithFilter2(t *testing.T) {

	nodes := readDoc(t, `f:
  a: frog
  b: dally
  c: log`)

	path, errPath := treeCreator.ParsePath("count(f(. == *og))")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: []
  Tag: !!int, Kind: ScalarNode, Anchor: 
  2
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorCountMultipleMatchesInside(t *testing.T) {

	nodes := readDoc(t, `f:
  a: [1,2]
  b: dally
  c: [3,4,5]`)

	path, errPath := treeCreator.ParsePath("f(count(a or c))")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [f]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  2
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorCountMultipleMatchesInsideSplat(t *testing.T) {

	nodes := readDoc(t, `f:
  a: [1,2,3]
  b: [1,2,3,4]
  c: [1,2,3,4,5]`)

	path, errPath := treeCreator.ParsePath("f(count( (a or c)*))")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [f]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  8
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorCountMultipleMatchesOutside(t *testing.T) {

	nodes := readDoc(t, `f:
  a: [1,2,3]
  b: [1,2,3,4]
  c: [1,2,3,4,5]`)

	path, errPath := treeCreator.ParsePath("f(a or c)(count(*))")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [f a]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  3
-- Node --
  Document 0, path: [f c]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  5
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorCountOfResults(t *testing.T) {

	nodes := readDoc(t, `- apple
- sdfsd
- apple`)

	path, errPath := treeCreator.ParsePath("count(*)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: []
  Tag: !!int, Kind: ScalarNode, Anchor: 
  3
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorCountNoMatches(t *testing.T) {

	nodes := readDoc(t, `- apple
- sdfsd
- apple`)

	path, errPath := treeCreator.ParsePath("count(5)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: []
  Tag: !!int, Kind: ScalarNode, Anchor: 
  0
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorDeleteAndWrite(t *testing.T) {

	nodes := readDoc(t, `a: 
  - b: apple
  - b: sdfsd
  - { b: apple, c: cat }`)

	path, errPath := treeCreator.ParsePath("(a .- (0 or 1)) or (a[0].b := frog)")
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
  Tag: !!seq, Kind: SequenceNode, Anchor: 
  - {b: frog, c: cat}

-- Node --
  Document 0, path: [a 0 b]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  frog
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorDeleteArray(t *testing.T) {

	nodes := readDoc(t, `a: 
  - b: apple
  - b: sdfsd
  - b: apple`)

	path, errPath := treeCreator.ParsePath("a .- (b == a*)")
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
  Tag: !!seq, Kind: SequenceNode, Anchor: 
  - b: sdfsd
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArraySimple(t *testing.T) {

	nodes := readDoc(t, `- b: apple`)

	path, errPath := treeCreator.ParsePath("[0]")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [0]
  Tag: !!map, Kind: MappingNode, Anchor: 
  b: apple
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSimpleAssign(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple`)

	path, errPath := treeCreator.ParsePath("a.b := frog")
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
  frog
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSimpleAssignSelf(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple`)

	path, errPath := treeCreator.ParsePath("a(. == apple)(. := frog)")
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
  frog
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSimpleAssignByFind(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple`)

	path, errPath := treeCreator.ParsePath("a(. == apple) := frog")
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
  frog
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArraySplat(t *testing.T) {

	nodes := readDoc(t, `- b: apple
- c: banana`)

	path, errPath := treeCreator.ParsePath("[*]")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [0]
  Tag: !!map, Kind: MappingNode, Anchor: 
  b: apple

-- Node --
  Document 0, path: [1]
  Tag: !!map, Kind: MappingNode, Anchor: 
  c: banana
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

func TestDataTreeNavigatorOrDeDupes(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: apple
  mad: things`)

	path, errPath := treeCreator.ParsePath("a.(cat or cat)")
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

func TestDataTreeNavigatorAnd(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: apple
  pat: apple
  cow: apple
  mad: things`)

	path, errPath := treeCreator.ParsePath("a.(*t and c*)")
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

func TestDataTreeNavigatorEqualsSimple(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: {b: apple, c: yes}
  pat: {b: banana}
`)

	path, errPath := treeCreator.ParsePath("a.(b == apple)")
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
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: apple, c: yes}
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorEqualsSelf(t *testing.T) {

	nodes := readDoc(t, `a: frog
b: cat
c: frog`)

	path, errPath := treeCreator.ParsePath("(a or b).(. == frog)")
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
  Tag: !!str, Kind: ScalarNode, Anchor: 
  frog
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorEqualsNested(t *testing.T) {

	nodes := readDoc(t, `a: {t: frog}
b: {t: cat}
c: {t: frog}`)

	path, errPath := treeCreator.ParsePath("(t == frog)")
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
  {t: frog}

-- Node --
  Document 0, path: [c]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {t: frog}
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArrayEqualsSelf(t *testing.T) {

	nodes := readDoc(t, `- cat
- dog
- frog`)

	path, errPath := treeCreator.ParsePath("(. == *og)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [1]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  dog

-- Node --
  Document 0, path: [2]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  frog
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArrayEqualsSelfSplatFirst(t *testing.T) {

	nodes := readDoc(t, `- cat
- dog
- frog`)

	path, errPath := treeCreator.ParsePath("*(. == *og)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [1]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  dog

-- Node --
  Document 0, path: [2]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  frog
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArrayEquals(t *testing.T) {

	nodes := readDoc(t, `- { b: apple, animal: rabbit }
- { b: banana, animal: cat }
- { b: corn, animal: dog }`)

	path, errPath := treeCreator.ParsePath("(b == apple or animal == dog)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [0]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: apple, animal: rabbit}

-- Node --
  Document 0, path: [2]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: corn, animal: dog}
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArrayEqualsDeep(t *testing.T) {

	nodes := readDoc(t, `apples:
  - { b: apple, animal: {legs: 2} }
  - { b: banana, animal: {legs: 4} }
  - { b: corn, animal: {legs: 6} }
`)

	path, errPath := treeCreator.ParsePath("apples(animal.legs == 4)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [apples 1]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: banana, animal: {legs: 4}}
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorEqualsTrickey(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: {b: apso, c: {d : yes}}
  pat: {b: apple, c: {d : no}}
  sat: {b: apsy, c: {d : yes}}
  fat: {b: apple}
`)

	path, errPath := treeCreator.ParsePath("a.(b == ap* and c.d == yes)")
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
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: apso, c: {d: yes}}

-- Node --
  Document 0, path: [a sat]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: apsy, c: {d: yes}}
`

	test.AssertResult(t, expected, resultsToString(results))
}
