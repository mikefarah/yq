package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

type expressionScenario struct {
	description string
	document    string
	expression  string
	expected    []string
	skipDoc     bool
}

func testScenario(t *testing.T, s *expressionScenario) {
	var results *list.List
	var err error
	if s.document != "" {
		results, err = EvaluateStream("sample.yaml", strings.NewReader(s.document), s.expression)
	} else {
		results, err = EvaluateExpression(s.expression)
	}

	if err != nil {
		t.Error(err)
		return
	}
	test.AssertResultComplexWithContext(t, s.expected, resultsToString(results), fmt.Sprintf("exp: %v\ndoc: %v", s.expression, s.document))
}

func documentScenarios(t *testing.T, title string, scenarios []expressionScenario) {
	f, err := os.Create(fmt.Sprintf("doc/%v.md", title))

	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	w.WriteString(fmt.Sprintf("# %v\n", title))
	w.WriteString(fmt.Sprintf("## Examples\n"))

	printer := NewPrinter(false, true, false, 2, true)

	for index, s := range scenarios {
		if !s.skipDoc {

			if s.description != "" {
				w.WriteString(fmt.Sprintf("### %v\n", s.description))
			} else {
				w.WriteString(fmt.Sprintf("### Example %v\n", index))
			}
			if s.document != "" {
				w.WriteString(fmt.Sprintf("sample.yml:\n"))
				w.WriteString(fmt.Sprintf("```yaml\n%v\n```\n", s.document))
			}
			if s.expression != "" {
				w.WriteString(fmt.Sprintf("Expression\n"))
				w.WriteString(fmt.Sprintf("```bash\nyq '%v' < sample.yml\n```\n", s.expression))
			}

			w.WriteString(fmt.Sprintf("Result\n"))

			var output bytes.Buffer
			var results *list.List
			var err error
			if s.document != "" {
				results, err = EvaluateStream("sample.yaml", strings.NewReader(s.document), s.expression)
			} else {
				results, err = EvaluateExpression(s.expression)
			}

			printer.PrintResults(results, bufio.NewWriter(&output))

			w.WriteString(fmt.Sprintf("```yaml\n%v```\n", output.String()))

			if err != nil {
				panic(err)
			}
		}

	}
	w.Flush()
}
