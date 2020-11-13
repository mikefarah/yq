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

func writeOrPanic(w *bufio.Writer, text string) {
	_, err := w.WriteString(text)
	if err != nil {
		panic(err)
	}
}

func documentScenarios(t *testing.T, title string, scenarios []expressionScenario) {
	f, err := os.Create(fmt.Sprintf("doc/%v.md", title))

	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	writeOrPanic(w, fmt.Sprintf("# %v\n", title))
	writeOrPanic(w, "## Examples\n")

	for index, s := range scenarios {
		if !s.skipDoc {

			if s.description != "" {
				writeOrPanic(w, fmt.Sprintf("### %v\n", s.description))
			} else {
				writeOrPanic(w, fmt.Sprintf("### Example %v\n", index))
			}
			if s.document != "" {
				writeOrPanic(w, "sample.yml:\n")
				writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.document))
			}
			if s.expression != "" {
				writeOrPanic(w, "Expression\n")
				writeOrPanic(w, fmt.Sprintf("```bash\nyq '%v' < sample.yml\n```\n", s.expression))
			}

			writeOrPanic(w, "Result\n")

			var output bytes.Buffer
			var err error
			printer := NewPrinter(bufio.NewWriter(&output), false, true, false, 2, true)

			if s.document != "" {
				node, err := treeCreator.ParsePath(s.expression)
				if err != nil {
					t.Error(err)
				}
				err = EvaluateStream("sample.yaml", strings.NewReader(s.document), node, printer)
				if err != nil {
					t.Error(err)
				}
			} else {
				err = EvaluateAllFileStreams(s.expression, []string{}, printer)
				if err != nil {
					t.Error(err)
				}
			}

			writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n", output.String()))

		}

	}
	w.Flush()
}
