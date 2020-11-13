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

	node, err := treeCreator.ParsePath(s.expression)
	if err != nil {
		t.Error(err)
		return
	}
	inputs := list.New()

	if s.document != "" {
		inputs, err = readDocuments(strings.NewReader(s.document), "sample.yml")
		if err != nil {
			t.Error(err)
			return
		}
	}

	results, err = treeNavigator.GetMatchingNodes(inputs, node)

	if err != nil {
		t.Error(err)
		return
	}
	test.AssertResultComplexWithContext(t, s.expected, resultsToString(results), fmt.Sprintf("exp: %v\ndoc: %v", s.expression, s.document))
}

func documentScenarios(t *testing.T, title string, scenarios []expressionScenario) {
	f, err := os.Create(fmt.Sprintf("doc/%v.md", title))

	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	w.WriteString(fmt.Sprintf("# %v\n", title))
	w.WriteString(fmt.Sprintf("## Examples\n"))

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
			var err error
			printer := NewPrinter(bufio.NewWriter(&output), false, true, false, 2, true)

			if s.document != "" {
				node, err := treeCreator.ParsePath(s.expression)
				if err != nil {
					t.Error(err)
				}
				err = EvaluateStream("sample.yaml", strings.NewReader(s.document), node, printer)
			} else {
				err = EvaluateAllFileStreams(s.expression, []string{}, printer)
			}

			w.WriteString(fmt.Sprintf("```yaml\n%v```\n", output.String()))

			if err != nil {
				t.Error(err)
			}
		}

	}
	w.Flush()
}
