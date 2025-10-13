package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var shellVariablesScenarios = []formatScenario{
	{
		description:    "Encode shell variables",
		subdescription: "Note that comments are dropped and values will be enclosed in single quotes as needed.",
		input: "" +
			"# comment" + "\n" +
			"name: Mike Wazowski" + "\n" +
			"eyes:" + "\n" +
			"  color: turquoise" + "\n" +
			"  number: 1" + "\n" +
			"friends:" + "\n" +
			"  - James P. Sullivan" + "\n" +
			"  - Celia Mae",
		expected: "" +
			"name='Mike Wazowski'" + "\n" +
			"eyes_color=turquoise" + "\n" +
			"eyes_number=1" + "\n" +
			"friends_0='James P. Sullivan'" + "\n" +
			"friends_1='Celia Mae'" + "\n",
	},
	{
		description:    "Encode shell variables: illegal variable names as key.",
		subdescription: "Keys that would be illegal as variable keys are adapted.",
		input: "" +
			"ascii_=_symbols: replaced with _" + "\n" +
			"\"ascii_\t_controls\": dropped (this example uses \\t)" + "\n" +
			"nonascii_\u05d0_characters: dropped" + "\n" +
			"effort_expe\u00f1ded_t\u00f2_preserve_accented_latin_letters: moderate (via unicode NFKD)" + "\n",
		expected: "" +
			"ascii___symbols='replaced with _'" + "\n" +
			"ascii__controls='dropped (this example uses \\t)'" + "\n" +
			"nonascii__characters=dropped" + "\n" +
			"effort_expended_to_preserve_accented_latin_letters='moderate (via unicode NFKD)'" + "\n",
	},
	{
		description:    "Encode shell variables: empty values, arrays and maps",
		subdescription: "Empty values are encoded to empty variables, but empty arrays and maps are skipped.",
		input:          "empty:\n  value:\n  array: []\n  map:   {}",
		expected:       "empty_value=" + "\n",
	},
	{
		description:    "Encode shell variables: single quotes in values",
		subdescription: "Single quotes in values are encoded as '\"'\"' (close single quote, double-quoted single quote, open single quote).",
		input:          "name: Miles O'Brien",
		expected:       `name='Miles O'"'"'Brien'` + "\n",
	},
	{
		description:    "Encode shell variables: custom separator",
		subdescription: "Use --shell-key-separator to specify a custom separator between keys. This is useful when the original keys contain underscores.",
		input: "" +
			"my_app:" + "\n" +
			"  db_config:" + "\n" +
			"    host: localhost" + "\n" +
			"    port: 5432",
		expected: "" +
			"my_app__db_config__host=localhost" + "\n" +
			"my_app__db_config__port=5432" + "\n",
		scenarioType: "shell-separator",
	},
}

func TestShellVariableScenarios(t *testing.T) {
	for _, s := range shellVariablesScenarios {
		//fmt.Printf("\t<%s> <%s>\n", s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewShellVariablesEncoder()))
		if s.scenarioType == "shell-separator" {
			// Save and restore the original separator
			originalSeparator := ConfiguredShellVariablesPreferences.KeySeparator
			ConfiguredShellVariablesPreferences.KeySeparator = "__"
			test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewShellVariablesEncoder()), s.description)
			ConfiguredShellVariablesPreferences.KeySeparator = originalSeparator
		} else {
			test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewShellVariablesEncoder()), s.description)
		}
	}
	genericScenarios := make([]interface{}, len(shellVariablesScenarios))
	for i, s := range shellVariablesScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "shellvariables", genericScenarios, documentShellVariableScenario)
}

func documentShellVariableScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)
	if s.skipDoc {
		return
	}
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression

	if s.scenarioType == "shell-separator" {
		writeOrPanic(w, "```bash\nyq -o=shell --shell-key-separator=\"__\" sample.yml\n```\n")
	} else if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=shell '%v' sample.yml\n```\n", expression))
	} else {
		writeOrPanic(w, "```bash\nyq -o=shell sample.yml\n```\n")
	}
	writeOrPanic(w, "will output\n")

	if s.scenarioType == "shell-separator" {
		// Save and restore the original separator
		originalSeparator := ConfiguredShellVariablesPreferences.KeySeparator
		ConfiguredShellVariablesPreferences.KeySeparator = "__"
		writeOrPanic(w, fmt.Sprintf("```sh\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewShellVariablesEncoder())))
		ConfiguredShellVariablesPreferences.KeySeparator = originalSeparator
	} else {
		writeOrPanic(w, fmt.Sprintf("```sh\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewShellVariablesEncoder())))
	}
}
