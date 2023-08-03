package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var luaScenarios = []formatScenario{
	{
		description:  "Basic example",
		input:        "hello: world\n? look: non-string keys\n: True\nnumbers: [123,456]\n",
		expected:     "return {[\"hello\"]=\"world\";[{[\"look\"]=\"non-string keys\";}]=true;[\"numbers\"]={123,456,};};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Sequence",
		input:        "- a\n- b\n- c\n",
		expected:     "return {\"a\",\"b\",\"c\",};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Mapping",
		input:        "a: b\nc:\n  d: e\nf: 0\n",
		expected:     "return {[\"a\"]=\"b\";[\"c\"]={[\"d\"]=\"e\";};[\"f\"]=0;};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Scalar str",
		input:        "str: |\n  foo\n  bar\n",
		expected:     "return {[\"str\"]=\"foo\\nbar\\n\";};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Scalar null",
		input:        "x: null\n",
		expected:     "return {[\"x\"]=nil;};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Scalar int",
		input:        "- 1\n- 2\n- 0x10\n- 0o30\n- -999\n",
		expected:     "return {1,2,0x10,24,-999,};\n",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "Scalar float",
		input:        "- 1.0\n- 3.14\n- 1e100\n- .Inf\n- .NAN\n",
		expected:     "return {1.0,3.14,1e100,(1/0),(0/0),};\n",
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
