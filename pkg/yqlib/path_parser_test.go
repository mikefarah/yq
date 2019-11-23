package yqlib

import (
	"testing"
	"github.com/mikefarah/yq/test"
)

var parsePathsTests = []struct {
	path          string
	expectedPaths []string
}{
	{"a.b", []string{"a", "b"}},
	{"a.b[0]", []string{"a", "b", "0"}},
	{"a.b.d[+]", []string{"a", "b", "d", "+"}},
}

func TestParsePath(t *testing.T) {
	for _, tt := range parsePathsTests {
		test.AssertResultComplex(t, tt.expectedPaths, ParsePath(tt.path))
	}
}

var nextYamlPathTests = []struct {
	path              string
	expectedElement   string
	expectedRemaining string
}{
	{"a.b", "a", "b"},
	{"a", "a", ""},
	{"a.b.c", "a", "b.c"},
	{"\"a.b\".c", "a.b", "c"},
	{"a.\"b.c\".d", "a", "\"b.c\".d"},
	{"[1].a.d", "1", "a.d"},
	{"a[0].c", "a", "[0].c"},
	{"[0]", "0", ""},
}

func TestNextYamlPath(t *testing.T) {
	for _, tt := range nextYamlPathTests {
		var element, remaining = nextYamlPath(tt.path)
		test.AssertResultWithContext(t, tt.expectedElement, element, tt)
		test.AssertResultWithContext(t, tt.expectedRemaining, remaining, tt)
	}
}
