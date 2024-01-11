//go:build !yq_nolua

package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var luaScenarios = []formatScenario{
	{
		description: "Basic input example",
		input: `return {
	["country"] = "Australia"; -- this place
	["cities"] = {
		"Sydney",
		"Melbourne",
		"Brisbane",
		"Perth",
	};
};
`,
		expected: `country: Australia
cities:
    - Sydney
    - Melbourne
    - Brisbane
    - Perth
`,
	},
	{
		skipDoc:     true,
		description: "path",
		expression:  ".cities[2] | path",
		input: `return {
	["country"] = "Australia"; -- this place
	["cities"] = {
		"Sydney",
		"Melbourne",
		"Brisbane",
		"Perth",
	};
};
`,
		expected: "- cities\n- 2\n",
	},
	{
		skipDoc:     true,
		description: "path",
		expression:  ".cities[2] | key",
		input: `return {
	["country"] = "Australia"; -- this place
	["cities"] = {
		"Sydney",
		"Melbourne",
		"Brisbane",
		"Perth",
	};
};
`,
		expected: "2\n",
	},
	{
		description:  "Basic output example",
		scenarioType: "encode",
		input: `---
country: Australia # this place
cities:
- Sydney
- Melbourne
- Brisbane
- Perth`,
		expected: `return {
	["country"] = "Australia"; -- this place
	["cities"] = {
		"Sydney",
		"Melbourne",
		"Brisbane",
		"Perth",
	};
};
`,
	},
	{
		description:  "Basic roundtrip",
		skipDoc:      true,
		scenarioType: "roundtrip",
		expression:   `.cities[0] = "Adelaide"`,
		input: `return {
	["country"] = "Australia"; -- this place
	["cities"] = {
		"Sydney",
		"Melbourne",
		"Brisbane",
		"Perth",
	};
};
`,
		expected: `return {
	["country"] = "Australia";
	["cities"] = {
		"Adelaide",
		"Melbourne",
		"Brisbane",
		"Perth",
	};
};
`,
	},
	{
		description:    "Unquoted keys",
		subdescription: "Uses the `--lua-unquoted` option to produce a nicer-looking output.",
		scenarioType:   "unquoted-encode",
		input: `---
country: Australia # this place
cities:
- Sydney
- Melbourne
- Brisbane
- Perth`,
		expected: `return {
	country = "Australia"; -- this place
	cities = {
		"Sydney",
		"Melbourne",
		"Brisbane",
		"Perth",
	};
};
`,
	},
	{
		description:    "Globals",
		subdescription: "Uses the `--lua-globals` option to export the values into the global scope.",
		scenarioType:   "globals-encode",
		input: `---
country: Australia # this place
cities:
- Sydney
- Melbourne
- Brisbane
- Perth`,
		expected: `country = "Australia"; -- this place
cities = {
	"Sydney",
	"Melbourne",
	"Brisbane",
	"Perth",
};
`,
	},
	{
		description: "Elaborate example",
		input: `---
hello: world
tables:
  like: this
  keys: values
  ? look: non-string keys
  : True
numbers:
  - decimal: 12345
  - hex: 0x7fabc123
  - octal: 0o30
  - float: 123.45
  - infinity: .inf
    plus_infinity: +.inf
    minus_infinity: -.inf
  - not: .nan
`,
		expected: `return {
	["hello"] = "world";
	["tables"] = {
		["like"] = "this";
		["keys"] = "values";
		[{
			["look"] = "non-string keys";
		}] = true;
	};
	["numbers"] = {
		{
			["decimal"] = 12345;
		},
		{
			["hex"] = 0x7fabc123;
		},
		{
			["octal"] = 24;
		},
		{
			["float"] = 123.45;
		},
		{
			["infinity"] = (1/0);
			["plus_infinity"] = (1/0);
			["minus_infinity"] = (-1/0);
		},
		{
			["not"] = (0/0);
		},
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
		input:        "str: |\n  foo\n  bar\nanother: 'single'\nand: \"double\"",
		expected:     "return {\n\t[\"str\"] = [[\nfoo\nbar\n]];\n\t[\"another\"] = 'single';\n\t[\"and\"] = \"double\";\n};\n",
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
	case "", "decode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewLuaDecoder(ConfiguredLuaPreferences), NewYamlEncoder(4, false, ConfiguredYamlPreferences)), s.description)
	case "encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewLuaEncoder(ConfiguredLuaPreferences)), s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewLuaDecoder(ConfiguredLuaPreferences), NewLuaEncoder(ConfiguredLuaPreferences)), s.description)
	case "unquoted-encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewLuaEncoder(LuaPreferences{
			DocPrefix:    "return ",
			DocSuffix:    ";\n",
			UnquotedKeys: true,
			Globals:      false,
		})), s.description)
	case "globals-encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewLuaEncoder(LuaPreferences{
			DocPrefix:    "return ",
			DocSuffix:    ";\n",
			UnquotedKeys: false,
			Globals:      true,
		})), s.description)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentLuaScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "", "decode":
		documentLuaDecodeScenario(w, s)
	case "encode", "unquoted-encode", "globals-encode":
		documentLuaEncodeScenario(w, s)
	case "roundtrip":
		documentLuaRoundTripScenario(w, s)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentLuaDecodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.lua file of:\n")
	writeOrPanic(w, fmt.Sprintf("```lua\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	expression := s.expression
	if expression == "" {
		expression = "."
	}
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -oy '%v' sample.lua\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewLuaDecoder(ConfiguredLuaPreferences), NewYamlEncoder(2, false, ConfiguredYamlPreferences))))
}

func documentLuaEncodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	prefs := ConfiguredLuaPreferences
	switch s.scenarioType {
	case "unquoted-encode":
		prefs = LuaPreferences{
			DocPrefix:    "return ",
			DocSuffix:    ";\n",
			UnquotedKeys: true,
			Globals:      false,
		}
	case "globals-encode":
		prefs = LuaPreferences{
			DocPrefix:    "return ",
			DocSuffix:    ";\n",
			UnquotedKeys: false,
			Globals:      true,
		}
	}
	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	switch s.scenarioType {
	case "unquoted-encode":
		writeOrPanic(w, "```bash\nyq -o=lua --lua-unquoted '.' sample.yml\n```\n")
	case "globals-encode":
		writeOrPanic(w, "```bash\nyq -o=lua --lua-globals '.' sample.yml\n```\n")
	default:
		writeOrPanic(w, "```bash\nyq -o=lua '.' sample.yml\n```\n")
	}
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```lua\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewLuaEncoder(prefs))))
}

func documentLuaRoundTripScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.lua file of:\n")
	writeOrPanic(w, fmt.Sprintf("```lua\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq '.' sample.lua\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```lua\n%v```\n\n", mustProcessFormatScenario(s, NewLuaDecoder(ConfiguredLuaPreferences), NewLuaEncoder(ConfiguredLuaPreferences))))
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
