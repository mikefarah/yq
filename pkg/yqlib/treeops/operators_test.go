package treeops

import (
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

type expressionScenario struct {
	document   string
	expression string
	expected   []string
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
