package yqlib

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var evaluateNodesScenario = []expressionScenario{
	{
		document:   `a: hello`,
		expression: `.a`,
		expected: []string{
			"D0, P[a], (!!str)::hello\n",
		},
	},
	{
		document:   `a: hello`,
		expression: `.`,
		expected: []string{
			"D0, P[], (!!map)::a: hello\n",
		},
	},
	{
		document:   `- a: "yes"`,
		expression: `.[] | has("a")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
		},
	},
}

func TestAllAtOnceEvaluateNodes(t *testing.T) {
	var evaluator = NewAllAtOnceEvaluator()
	// logging.SetLevel(logging.DEBUG, "")
	for _, tt := range evaluateNodesScenario {
		decoder := NewYamlDecoder(ConfiguredYamlPreferences)
		reader := bufio.NewReader(strings.NewReader(tt.document))
		err := decoder.Init(reader)
		if err != nil {
			t.Error(err)
			return
		}
		candidateNode, errorReading := decoder.Decode()

		if errorReading != nil {
			t.Error(errorReading)
			return
		}

		list, _ := evaluator.EvaluateNodes(tt.expression, candidateNode)
		test.AssertResultComplex(t, tt.expected, resultsToString(t, list))
	}
}

func TestTomlDecoderCanBeReinitializedAcrossDocuments(t *testing.T) {
	decoder := NewTomlDecoder()

	firstDocuments, err := ReadDocuments(strings.NewReader("id = \"Foobar\"\n"), decoder)
	if err != nil {
		t.Fatalf("failed to read first TOML document: %v", err)
	}
	if firstDocuments.Len() != 1 {
		t.Fatalf("expected first document count to be 1, got %d", firstDocuments.Len())
	}
	test.AssertResult(t, "Foobar", firstDocuments.Front().Value.(*CandidateNode).Content[1].Value)

	secondDocuments, err := ReadDocuments(strings.NewReader("id = \"Banana\"\n"), decoder)
	if err != nil {
		t.Fatalf("failed to read second TOML document: %v", err)
	}
	if secondDocuments.Len() != 1 {
		t.Fatalf("expected second document count to be 1, got %d", secondDocuments.Len())
	}
	test.AssertResult(t, "Banana", secondDocuments.Front().Value.(*CandidateNode).Content[1].Value)
}

func TestReadDocumentsLeavesReaderOpen(t *testing.T) {
	filename := t.TempDir() + "/input.yml"
	if err := os.WriteFile(filename, []byte("a: apple\n"), 0600); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open input file: %v", err)
	}
	defer file.Close()

	if _, err := ReadDocuments(file, NewYamlDecoder(ConfiguredYamlPreferences)); err != nil {
		t.Fatalf("failed to read documents: %v", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("ReadDocuments closed the caller-owned reader: %v", err)
	}
}
