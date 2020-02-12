package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

var parser = NewPathParser()

var parsePathsTests = []struct {
	path          string
	expectedPaths []interface{}
}{
	{"a.b", append(make([]interface{}, 0), "a", "b")},
	{"a.b.**", append(make([]interface{}, 0), "a", "b", "**")},
	{"a.b.*", append(make([]interface{}, 0), "a", "b", "*")},
	{"a.b[0]", append(make([]interface{}, 0), "a", "b", int64(0))},
	{"a.b.0", append(make([]interface{}, 0), "a", "b", "0")},
	{"a.b.d[+]", append(make([]interface{}, 0), "a", "b", "d", "+")},
	{"a", append(make([]interface{}, 0), "a")},
	{"a.b.c", append(make([]interface{}, 0), "a", "b", "c")},
	{"\"a.b\".c", append(make([]interface{}, 0), "a.b", "c")},
	{"a.\"b.c\".d", append(make([]interface{}, 0), "a", "b.c", "d")},
	{"[1].a.d", append(make([]interface{}, 0), int64(1), "a", "d")},
	{"a[0].c", append(make([]interface{}, 0), "a", int64(0), "c")},
	{"[0]", append(make([]interface{}, 0), int64(0))},
}

func TestPathParserParsePath(t *testing.T) {
	for _, tt := range parsePathsTests {
		test.AssertResultComplex(t, tt.expectedPaths, parser.ParsePath(tt.path))
	}
}

func TestPathParserMatchesNextPathElementSplat(t *testing.T) {
	var node = NodeContext{Head: "*"}
	test.AssertResult(t, true, parser.MatchesNextPathElement(node, ""))
}

func TestPathParserMatchesNextPathElementDeepSplat(t *testing.T) {
	var node = NodeContext{Head: "**"}
	test.AssertResult(t, true, parser.MatchesNextPathElement(node, ""))
}

func TestPathParserMatchesNextPathElementAppendArrayValid(t *testing.T) {
	var node = NodeContext{Head: "+"}
	test.AssertResult(t, true, parser.MatchesNextPathElement(node, "3"))
}

func TestPathParserMatchesNextPathElementAppendArrayInvalid(t *testing.T) {
	var node = NodeContext{Head: "+"}
	test.AssertResult(t, false, parser.MatchesNextPathElement(node, "cat"))
}

func TestPathParserMatchesNextPathElementPrefixMatchesWhole(t *testing.T) {
	var node = NodeContext{Head: "cat*"}
	test.AssertResult(t, true, parser.MatchesNextPathElement(node, "cat"))
}

func TestPathParserMatchesNextPathElementPrefixMatchesStart(t *testing.T) {
	var node = NodeContext{Head: "cat*"}
	test.AssertResult(t, true, parser.MatchesNextPathElement(node, "caterpillar"))
}

func TestPathParserMatchesNextPathElementPrefixMismatch(t *testing.T) {
	var node = NodeContext{Head: "cat*"}
	test.AssertResult(t, false, parser.MatchesNextPathElement(node, "dog"))
}

func TestPathParserMatchesNextPathElementExactMatch(t *testing.T) {
	var node = NodeContext{Head: "farahtek"}
	test.AssertResult(t, true, parser.MatchesNextPathElement(node, "farahtek"))
}

func TestPathParserMatchesNextPathElementExactMismatch(t *testing.T) {
	var node = NodeContext{Head: "farahtek"}
	test.AssertResult(t, false, parser.MatchesNextPathElement(node, "othertek"))
}
