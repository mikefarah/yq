package yqlib

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

var LoadYamlPreferences = YamlPreferences{
	LeadingContentPreProcessing: false,
	PrintDocSeparators:          true,
	UnwrapScalar:                true,
	EvaluateTogether:            false,
}

type loadPrefs struct {
	decoder Decoder
}

func loadString(filename string) (*CandidateNode, error) {
	// ignore CWE-22 gosec issue - that's more targeted for http based apps that run in a public directory,
	// and ensuring that it's not possible to give a path to a file outside that directory.

	filebytes, err := os.ReadFile(filename) // #nosec
	if err != nil {
		return nil, err
	}

	return &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: string(filebytes)}, nil
}

func loadWithDecoder(filename string, decoder Decoder) (*CandidateNode, error) {
	if decoder == nil {
		return nil, fmt.Errorf("could not load %s", filename)
	}

	file, err := os.Open(filename) // #nosec
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)

	documents, err := readDocuments(reader, filename, 0, decoder)
	if err != nil {
		return nil, err
	}

	if documents.Len() == 0 {
		// return null candidate
		return &CandidateNode{Kind: ScalarNode, Tag: "!!null"}, nil
	} else if documents.Len() == 1 {
		candidate := documents.Front().Value.(*CandidateNode)
		return candidate, nil

	}
	sequenceNode := &CandidateNode{Kind: SequenceNode}
	for doc := documents.Front(); doc != nil; doc = doc.Next() {
		sequenceNode.AddChild(doc.Value.(*CandidateNode))
	}
	return sequenceNode, nil
}

func loadStringOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("loadString")
	if ConfiguredSecurityPreferences.DisableFileOps {
		return Context{}, fmt.Errorf("file operations have been disabled")
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}
		if rhs.MatchingNodes.Front() == nil {
			return Context{}, fmt.Errorf("filename expression returned nil")
		}
		nameCandidateNode := rhs.MatchingNodes.Front().Value.(*CandidateNode)

		filename := nameCandidateNode.Value

		contentsCandidate, err := loadString(filename)
		if err != nil {
			return Context{}, fmt.Errorf("failed to load %v: %w", filename, err)
		}

		results.PushBack(contentsCandidate)

	}

	return context.ChildContext(results), nil
}

func loadOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("loadOperator")
	if ConfiguredSecurityPreferences.DisableFileOps {
		return Context{}, fmt.Errorf("file operations have been disabled")
	}

	loadPrefs := expressionNode.Operation.Preferences.(loadPrefs)

	// need to evaluate the 1st parameter against the context
	// and return the data accordingly.

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}
		if rhs.MatchingNodes.Front() == nil {
			return Context{}, fmt.Errorf("filename expression returned nil")
		}
		nameCandidateNode := rhs.MatchingNodes.Front().Value.(*CandidateNode)

		filename := nameCandidateNode.Value

		contentsCandidate, err := loadWithDecoder(filename, loadPrefs.decoder)
		if err != nil {
			return Context{}, fmt.Errorf("failed to load %v: %w", filename, err)
		}

		results.PushBack(contentsCandidate)

	}

	return context.ChildContext(results), nil
}
