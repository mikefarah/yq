package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
	yaml "gopkg.in/yaml.v3"
)

var parseStyleTests = []struct {
	customStyle   string
	expectedStyle yaml.Style
}{
	{"", 0},
	{"tagged", yaml.TaggedStyle},
	{"double", yaml.DoubleQuotedStyle},
	{"single", yaml.SingleQuotedStyle},
	{"folded", yaml.FoldedStyle},
	{"literal", yaml.LiteralStyle},
}

func TestValueParserStyleTag(t *testing.T) {
	for _, tt := range parseStyleTests {
		actual := NewValueParser().Parse("cat", "", tt.customStyle)
		test.AssertResultWithContext(t, tt.expectedStyle, actual.Style, tt.customStyle)
	}
}

var parseValueTests = []struct {
	argument        string
	customTag       string
	expectedTag     string
	testDescription string
}{
	{"true", "!!str", "!!str", "boolean forced as string"},
	{"3", "!!int", "!!int", "int"},
	{"cat", "", "", "default"},
}

func TestValueParserParse(t *testing.T) {
	for _, tt := range parseValueTests {
		actual := NewValueParser().Parse(tt.argument, tt.customTag, "")
		test.AssertResultWithContext(t, tt.argument, actual.Value, tt.testDescription)
		test.AssertResultWithContext(t, tt.expectedTag, actual.Tag, tt.testDescription)
		test.AssertResult(t, yaml.ScalarNode, actual.Kind)
	}
}

func TestValueParserParseEmptyArray(t *testing.T) {
	actual := NewValueParser().Parse("[]", "", "")
	test.AssertResult(t, "!!seq", actual.Tag)
	test.AssertResult(t, yaml.SequenceNode, actual.Kind)
}
