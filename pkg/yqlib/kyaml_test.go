package yqlib

import (
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func testKYamlScenario(t *testing.T, s formatScenario) {
	prefs := ConfiguredKYamlPreferences.Copy()
	prefs.Indent = s.indent
	prefs.UnwrapScalar = false

	switch s.scenarioType {
	case "encode":
		test.AssertResultWithContext(
			t,
			s.expected,
			mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewKYamlEncoder(prefs)),
			s.description,
		)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func TestKYamlFormatScenarios(t *testing.T) {
	scenarios := []formatScenario{
		{
			description:  "encode flow mapping and sequence",
			scenarioType: "encode",
			indent:       2,
			input:        "a: b\nc:\n  - d\n",
			expected: "{\n" +
				"  a: \"b\",\n" +
				"  c: [\n" +
				"    \"d\",\n" +
				"  ],\n" +
				"}\n",
		},
		{
			description:  "encode non-string scalars",
			scenarioType: "encode",
			indent:       2,
			input: "a: 12\n" +
				"b: true\n" +
				"c: null\n" +
				"d: \"true\"\n",
			expected: "{\n" +
				"  a: 12,\n" +
				"  b: true,\n" +
				"  c: null,\n" +
				"  d: \"true\",\n" +
				"}\n",
		},
		{
			description:  "quote non-identifier keys",
			scenarioType: "encode",
			indent:       2,
			input:        "\"1a\": b\n\"has space\": c\n",
			expected: "{\n" +
				"  \"1a\": \"b\",\n" +
				"  \"has space\": \"c\",\n" +
				"}\n",
		},
		{
			description:  "escape quoted strings",
			scenarioType: "encode",
			indent:       2,
			input:        "a: \"line1\\nline2\\t\\\"q\\\"\"\n",
			expected: "{\n" +
				"  a: \"line1\\nline2\\t\\\"q\\\"\",\n" +
				"}\n",
		},
		{
			description:  "preserve comments when encoding",
			scenarioType: "encode",
			indent:       2,
			input: "# leading\n" +
				"a: 1 # a line\n" +
				"# head b\n" +
				"b: 2\n" +
				"c:\n" +
				"  # head d\n" +
				"  - d # d line\n" +
				"  - e\n" +
				"# trailing\n",
			expected: "# leading\n" +
				"{\n" +
				"  a: 1, # a line\n" +
				"  # head b\n" +
				"  b: 2,\n" +
				"  c: [\n" +
				"    # head d\n" +
				"    \"d\", # d line\n" +
				"    \"e\",\n" +
				"  ],\n" +
				"  # trailing\n" +
				"}\n",
		},
	}

	for _, s := range scenarios {
		testKYamlScenario(t, s)
	}
}
