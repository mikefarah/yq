package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var luaScenarios = []formatScenario{
	{
		description: "Basic example",
		input: `---
hello: world
? look: non-string keys
: True
numbers: [123,456]
`,
		expected: `return {
	["hello"] = "world";
	[{
		["look"] = "non-string keys";
	}] = true;
	["numbers"] = {
		123,
		456,
	};
};
`,
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Sequence",
		input:        "- a\n- b\n- c\n",
		expected:     "return {\n\t\"a\",\n\t\"b\",\n\t\"c\",\n};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Mapping",
		input:        "a: b\nc:\n  d: e\nf: 0\n",
		expected:     "return {\n\t[\"a\"] = \"b\";\n\t[\"c\"] = {\n\t\t[\"d\"] = \"e\";\n\t};\n\t[\"f\"] = 0;\n};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Scalar str",
		input:        "str: |\n  foo\n  bar\n",
		expected:     "return {\n\t[\"str\"] = \"foo\\nbar\\n\";\n};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Scalar null",
		input:        "x: null\n",
		expected:     "return {\n\t[\"x\"] = nil;\n};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Scalar int",
		input:        "- 1\n- 2\n- 0x10\n- 0o30\n- -999\n",
		expected:     "return {\n\t1,\n\t2,\n\t0x10,\n\t24,\n\t-999,\n};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Scalar float",
		input:        "- 1.0\n- 3.14\n- 1e100\n- .Inf\n- .NAN\n",
		expected:     "return {\n\t1.0,\n\t3.14,\n\t1e100,\n\t(1/0),\n\t(0/0),\n};\n",
		scenarioType: "encode",
	},
}

func testLuaScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewLuaEncoder(ConfiguredLuaPreferences)), s.description)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentLuaScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "encode":
		documentLuaEncodeScenario(w, s)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentLuaEncodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq -o=lua '.' sample.yml\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```lua\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewLuaEncoder(ConfiguredLuaPreferences))))
}

func TestLuaScenarios(t *testing.T) {
	for _, tt := range luaScenarios {
		testLuaScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(luaScenarios))
	for i, s := range luaScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "lua", genericScenarios, documentLuaScenario)
}
