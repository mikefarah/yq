package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var formattingExpressionScenarios = []formatScenario{
	{
		description: "Using expression files and comments",
		input:       "a:\n  b: old",
		expression:  "\n# This is a yq expression that updates the map\n# for several great reasons outlined here.\n\n.a.b = \"new\" # line comment here\n| .a.c = \"frog\"\n\n# Now good things will happen.\n",
		expected:    "a:\n  b: new\n  c: frog\n",
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
	writeOrPanic(w, fmt.Sprintf("```bash%v```\n", s.expression))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq --from-file update.yq sample.yml\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewYamlEncoder(2, false, ConfiguredYamlPreferences))))
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
