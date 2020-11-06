package yqlib

import (
	"container/list"
)

func resultsToString(results *list.List) []string {
	var pretty []string = make([]string, 0)
	for el := results.Front(); el != nil; el = el.Next() {
		n := el.Value.(*CandidateNode)
		pretty = append(pretty, NodeToString(n))
	}
	return pretty
}

// func TestDataTreeNavigatorDeleteSimple(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   b: apple
//   c: camel`)

// 	path, errPath := treeCreator.ParsePath("a .- b")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a]
//   Tag: !!map, Kind: MappingNode, Anchor:
//   c: camel
// `
// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorDeleteTwice(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   b: apple
//   c: camel
//   d: dingo`)

// 	path, errPath := treeCreator.ParsePath("a .- b OR a .- c")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a]
//   Tag: !!map, Kind: MappingNode, Anchor:
//   d: dingo
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorDeleteWithUnion(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   b: apple
//   c: camel
//   d: dingo`)

// 	path, errPath := treeCreator.ParsePath("a .- (b OR c)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a]
//   Tag: !!map, Kind: MappingNode, Anchor:
//   d: dingo
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorDeleteByIndex(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   - b: apple
//   - b: sdfsd
//   - b: apple`)

// 	path, errPath := treeCreator.ParsePath("(a .- (0 or 1))")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a]
//   Tag: !!seq, Kind: SequenceNode, Anchor:
//   - b: apple
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorDeleteByFind(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   - b: apple
//   - b: sdfsd
//   - b: apple`)

// 	path, errPath := treeCreator.ParsePath("(a .- (* == apple))")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a]
//   Tag: !!seq, Kind: SequenceNode, Anchor:
//   - b: sdfsd
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorDeleteArrayByFind(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   - apple
//   - sdfsd
//   - apple`)

// 	path, errPath := treeCreator.ParsePath("(a .- (. == apple))")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a]
//   Tag: !!seq, Kind: SequenceNode, Anchor:
//   - sdfsd
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorDeleteViaSelf(t *testing.T) {

// 	nodes := readDoc(t, `- apple
// - sdfsd
// - apple`)

// 	path, errPath := treeCreator.ParsePath(". .- (. == apple)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: []
//   Tag: !!seq, Kind: SequenceNode, Anchor:
//   - sdfsd
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorFilterWithSplat(t *testing.T) {

// 	nodes := readDoc(t, `f:
//   a: frog
//   b: dally
//   c: log`)

// 	path, errPath := treeCreator.ParsePath(".f | .[] == \"frog\"")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [f]
//   Tag: !!int, Kind: ScalarNode, Anchor:
//   2
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCountAndCollectWithFilterCmd(t *testing.T) {

// 	nodes := readDoc(t, `f:
//   a: frog
//   b: dally
//   c: log`)

// 	path, errPath := treeCreator.ParsePath(".f | .[] == *og ")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [f]
//   Tag: !!int, Kind: ScalarNode, Anchor:
//   2
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCollectWithFilter(t *testing.T) {

// 	nodes := readDoc(t, `f:
//   a: frog
//   b: dally
//   c: log`)

// 	path, errPath := treeCreator.ParsePath("f(collect(. == *og))")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [f]
//   Tag: , Kind: SequenceNode, Anchor:
//   - frog
// - log
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCountWithFilter2(t *testing.T) {

// 	nodes := readDoc(t, `f:
//   a: frog
//   b: dally
//   c: log`)

// 	path, errPath := treeCreator.ParsePath("count(f(. == *og))")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: []
//   Tag: !!int, Kind: ScalarNode, Anchor:
//   2
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCollectWithFilter2(t *testing.T) {

// 	nodes := readDoc(t, `f:
//   a: frog
//   b: dally
//   c: log`)

// 	path, errPath := treeCreator.ParsePath("collect(f(. == *og))")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: []
//   Tag: , Kind: SequenceNode, Anchor:
//   - frog
// - log
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCountMultipleMatchesInside(t *testing.T) {

// 	nodes := readDoc(t, `f:
//   a: [1,2]
//   b: dally
//   c: [3,4,5]`)

// 	path, errPath := treeCreator.ParsePath("f | count(a or c)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [f]
//   Tag: !!int, Kind: ScalarNode, Anchor:
//   2
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCollectMultipleMatchesInside(t *testing.T) {

// 	nodes := readDoc(t, `f:
//   a: [1,2]
//   b: dally
//   c: [3,4,5]`)

// 	path, errPath := treeCreator.ParsePath("f | collect(a or c)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [f]
//   Tag: , Kind: SequenceNode, Anchor:
//   - [1, 2]
// - [3, 4, 5]
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCountMultipleMatchesInsideSplat(t *testing.T) {

// 	nodes := readDoc(t, `f:
//   a: [1,2,3]
//   b: [1,2,3,4]
//   c: [1,2,3,4,5]`)

// 	path, errPath := treeCreator.ParsePath("f(count( (a or c)*))")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [f]
//   Tag: !!int, Kind: ScalarNode, Anchor:
//   8
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCountMultipleMatchesOutside(t *testing.T) {

// 	nodes := readDoc(t, `f:
//   a: [1,2,3]
//   b: [1,2,3,4]
//   c: [1,2,3,4,5]`)

// 	path, errPath := treeCreator.ParsePath("f(a or c)(count(*))")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [f a]
//   Tag: !!int, Kind: ScalarNode, Anchor:
//   3
// -- Node --
//   Document 0, path: [f c]
//   Tag: !!int, Kind: ScalarNode, Anchor:
//   5
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCountOfResults(t *testing.T) {

// 	nodes := readDoc(t, `- apple
// - sdfsd
// - apple`)

// 	path, errPath := treeCreator.ParsePath("count(*)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: []
//   Tag: !!int, Kind: ScalarNode, Anchor:
//   3
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorCountNoMatches(t *testing.T) {

// 	nodes := readDoc(t, `- apple
// - sdfsd
// - apple`)

// 	path, errPath := treeCreator.ParsePath("count(5)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: []
//   Tag: !!int, Kind: ScalarNode, Anchor:
//   0
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorDeleteAndWrite(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   - b: apple
//   - b: sdfsd
//   - { b: apple, c: cat }`)

// 	path, errPath := treeCreator.ParsePath("(a .- (0 or 1)) or (a[0].b := frog)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a]
//   Tag: !!seq, Kind: SequenceNode, Anchor:
//   - {b: frog, c: cat}

// -- Node --
//   Document 0, path: [a 0 b]
//   Tag: !!str, Kind: ScalarNode, Anchor:
//   frog
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorDeleteArray(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   - b: apple
//   - b: sdfsd
//   - b: apple`)

// 	path, errPath := treeCreator.ParsePath("a .- (b == a*)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a]
//   Tag: !!seq, Kind: SequenceNode, Anchor:
//   - b: sdfsd
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorArraySimple(t *testing.T) {

// 	nodes := readDoc(t, `- b: apple`)

// 	path, errPath := treeCreator.ParsePath("[0]")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [0]
//   Tag: !!map, Kind: MappingNode, Anchor:
//   b: apple
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorSimpleAssignByFind(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   b: apple`)

// 	path, errPath := treeCreator.ParsePath("a(. == apple) := frog")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a b]
//   Tag: !!str, Kind: ScalarNode, Anchor:
//   frog
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorOrDeDupes(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   cat: apple
//   mad: things`)

// 	path, errPath := treeCreator.ParsePath("a.(cat or cat)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a cat]
//   Tag: !!str, Kind: ScalarNode, Anchor:
//   apple
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }

// func TestDataTreeNavigatorAnd(t *testing.T) {

// 	nodes := readDoc(t, `a:
//   cat: apple
//   pat: apple
//   cow: apple
//   mad: things`)

// 	path, errPath := treeCreator.ParsePath("a.(*t and c*)")
// 	if errPath != nil {
// 		t.Error(errPath)
// 	}
// 	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

// 	if errNav != nil {
// 		t.Error(errNav)
// 	}

// 	expected := `
// -- Node --
//   Document 0, path: [a cat]
//   Tag: !!str, Kind: ScalarNode, Anchor:
//   apple
// `

// 	test.AssertResult(t, expected, resultsToString(results))
// }
