package treeops

import (
	"bytes"
	"fmt"

	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

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

var log = logging.MustGetLogger("yq-treeops")

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
