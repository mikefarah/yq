package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

var parsePathsTests = []struct {
	path          string
	expectedPaths []string
}{
	{"a.b", []string{"a", "b"}},
	{"a.b[0]", []string{"a", "b", "0"}},
	{"a.b.d[+]", []string{"a", "b", "d", "+"}},
	{"a", []string{"a"}},
	{"a.b.c", []string{"a", "b", "c"}},
	{"\"a.b\".c", []string{"a.b", "c"}},
	{"a.\"b.c\".d", []string{"a", "b.c", "d"}},
	{"[1].a.d", []string{"1", "a", "d"}},
	{"a[0].c", []string{"a", "0", "c"}},
	{"[0]", []string{"0"}},
}

func TestParsePath(t *testing.T) {
	for _, tt := range parsePathsTests {
		test.AssertResultComplex(t, tt.expectedPaths, NewPathParser().ParsePath(tt.path))
	}
}
