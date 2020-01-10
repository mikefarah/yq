package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v2/test"
)

var parseValueTests = []struct {
	argument        string
	expectedResult  interface{}
	testDescription string
}{
	{"true", true, "boolean"},
	{"\"true\"", "true", "boolean as string"},
	{"3.4", 3.4, "number"},
	{"\"3.4\"", "3.4", "number as string"},
	{"", "", "empty string"},
	{"1212121", int64(1212121), "big number"},
}

func TestParseValue(t *testing.T) {
	for _, tt := range parseValueTests {
		test.AssertResultWithContext(t, tt.expectedResult, NewValueParser().ParseValue(tt.argument), tt.testDescription)
	}
}
