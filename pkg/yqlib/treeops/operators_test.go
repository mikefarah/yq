package treeops

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
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
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}
	test.AssertResultComplexWithContext(t, s.expected, resultsToString(results), s.expression)
}
