package treeops

import (
	"bytes"
	"fmt"

	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("yq-treeops")

type CandidateNode struct {
	Node     *yaml.Node    // the actual node
	Path     []interface{} /// the path we took to get to this node
	Document uint          // the document index of this node

	// middle nodes are nodes that match along the original path, but not a
	// target match of the path. This is only relevant when ShouldOnlyDeeplyVisitLeaves is false.
	IsMiddleNode bool
}

func (n *CandidateNode) getKey() string {
	return fmt.Sprintf("%v - %v", n.Document, n.Path)
}

type YqTreeLib interface {
	Get(rootNode *yaml.Node, path string) ([]*CandidateNode, error)
	// GetForMerge(rootNode *yaml.Node, path string, arrayMergeStrategy ArrayMergeStrategy) ([]*NodeContext, error)
	// Update(rootNode *yaml.Node, updateCommand UpdateCommand, autoCreate bool) error
	// New(path string) yaml.Node

	// PathStackToString(pathStack []interface{}) string
	// MergePathStackToString(pathStack []interface{}, arrayMergeStrategy ArrayMergeStrategy) string
}

type lib struct {
	treeCreator PathTreeCreator
}

func NodeToString(node *CandidateNode) string {
	if !log.IsEnabledFor(logging.DEBUG) {
		return ""
	}
	value := node.Node
	if value == nil {
		return "-- node is nil --"
	}
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	errorEncoding := encoder.Encode(value)
	if errorEncoding != nil {
		log.Error("Error debugging node, %v", errorEncoding.Error())
	}
	encoder.Close()
	return fmt.Sprintf(`-- Node --
  Document %v, path: %v
  Tag: %v, Kind: %v, Anchor: %v
  %v`, node.Document, node.Path, value.Tag, KindString(value.Kind), value.Anchor, buf.String())
}

func KindString(kind yaml.Kind) string {
	switch kind {
	case yaml.ScalarNode:
		return "ScalarNode"
	case yaml.SequenceNode:
		return "SequenceNode"
	case yaml.MappingNode:
		return "MappingNode"
	case yaml.DocumentNode:
		return "DocumentNode"
	case yaml.AliasNode:
		return "AliasNode"
	default:
		return "unknown!"
	}
}
