package yqlib

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

type expressionScenario struct {
	description string
	document    string
	expression  string
	expected    []string
}

func testScenario(t *testing.T, s *expressionScenario) {

	nodes := readDoc(t, s.document)
	path, errPath := treeCreator.ParsePath(s.expression)
	if errPath != nil {
		t.Error(errPath)
		return
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
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

		nodes := readDoc(t, s.document)
		path, errPath := treeCreator.ParsePath(s.expression)
		if errPath != nil {
			t.Error(errPath)
			return
		}
		var output bytes.Buffer
		results, err := treeNavigator.GetMatchingNodes(nodes, path)
		printer.PrintResults(results, bufio.NewWriter(&output))

		w.WriteString(fmt.Sprintf("```yaml\n%v```\n", output.String()))

		if err != nil {
			panic(err)
		}

	}
	w.Flush()
}
