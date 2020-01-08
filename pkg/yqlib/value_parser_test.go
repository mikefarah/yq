package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
	yaml "gopkg.in/yaml.v3"
)

var parseValueTests = []struct {
	argument        string
	customTag       string
	expectedTag     string
	testDescription string
}{
	{"true", "", "!!bool", "boolean"},
	{"true", "!!string", "!!string", "boolean forced as string"},
	{"3.4", "", "!!float", "float"},
	{"1212121", "", "!!int", "big number"},
	{"1212121.1", "", "!!float", "big float number"},
	{"3", "", "!!int", "int"},
	{"null", "", "!!null", "null"},
}

func TestValueParserParse(t *testing.T) {
	for _, tt := range parseValueTests {
		actual := NewValueParser().Parse(tt.argument, tt.customTag)
		test.AssertResultWithContext(t, tt.argument, actual.Value, tt.testDescription)
		test.AssertResultWithContext(t, tt.expectedTag, actual.Tag, tt.testDescription)
		test.AssertResult(t, yaml.ScalarNode, actual.Kind)
	}
}

func TestValueParserParseEmptyArray(t *testing.T) {
	actual := NewValueParser().Parse("[]", "")
	test.AssertResult(t, "!!seq", actual.Tag)
	test.AssertResult(t, yaml.SequenceNode, actual.Kind)
}
