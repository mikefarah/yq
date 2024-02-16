package yqlib

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var formattingExpressionScenarios = []formatScenario{
	{
		description: "Using expression files and comments",
		skipDoc:     true,
		input:       "a:\n  b: old",
		expression:  "#! yq\n\n# This is a yq expression that updates the map\n# for several great reasons outlined here.\n\n.a.b = \"new\" # line comment here\n| .a.c = \"frog\"\n\n# Now good things will happen.\n",
		expected:    "a:\n  b: new\n  c: frog\n",
	},
	{
		description:    "Using expression files and comments",
		subdescription: "Note that you can execute the file directly - but make sure you make the expression file executable.",
		input:          "a:\n  b: old",
		expression:     "#! yq\n\n# This is a yq expression that updates the map\n# for several great reasons outlined here.\n\n.a.b = \"new\" # line comment here\n| .a.c = \"frog\"\n\n# Now good things will happen.\n",
		expected:       "a:\n  b: new\n  c: frog\n",
		scenarioType:   "shebang",
	},
	{
		description:    "Flags in expression files",
		subdescription: "You can specify flags on the shebang line, this only works when executing the file directly.",
		input:          "a:\n  b: old",
		expression:     "#! yq -oj\n\n# This is a yq expression that updates the map\n# for several great reasons outlined here.\n\n.a.b = \"new\" # line comment here\n| .a.c = \"frog\"\n\n# Now good things will happen.\n",
		expected:       "a:\n  b: new\n  c: frog\n",
		scenarioType:   "shebang-json",
	},
	{
		description:    "Commenting out yq expressions",
		subdescription: "Note that `c` is no longer set to 'frog'. In this example we're calling yq directly and passing the expression file into `--from-file`, this is no different from executing the expression file directly.",
		input:          "a:\n  b: old",
		expression:     "#! yq\n# This is a yq expression that updates the map\n# for several great reasons outlined here.\n\n.a.b = \"new\" # line comment here\n# | .a.c = \"frog\"\n\n# Now good things will happen.\n",
		expected:       "a:\n  b: new\n",
	},
}

func documentExpressionScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.yaml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "And an 'update.yq' expression file of:\n")
	writeOrPanic(w, fmt.Sprintf("```bash\n%v```\n", s.expression))

	writeOrPanic(w, "then\n")
	if strings.HasPrefix(s.scenarioType, "shebang") {
		writeOrPanic(w, "```bash\n./update.yq sample.yaml\n```\n")
	} else {
		writeOrPanic(w, "```bash\nyq --from-file update.yq sample.yml\n```\n")
	}
	writeOrPanic(w, "will output\n")
	encoder := NewYamlEncoder(2, false, ConfiguredYamlPreferences)

	if s.scenarioType == "shebang-json" {
		encoder = NewJSONEncoder(2, false, false)
	}

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), encoder)))
}

func TestExpressionCommentScenarios(t *testing.T) {
	for _, tt := range formattingExpressionScenarios {
		test.AssertResultComplexWithContext(t, tt.expected,
			mustProcessFormatScenario(tt, NewYamlDecoder(ConfiguredYamlPreferences), NewYamlEncoder(2, false, ConfiguredYamlPreferences)),
			tt.description)
	}
	genericScenarios := make([]interface{}, len(formattingExpressionScenarios))
	for i, s := range formattingExpressionScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "formatting-expressions", genericScenarios, documentExpressionScenario)
}
