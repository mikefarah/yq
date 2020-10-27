package treeops

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
)

type CandidateNode struct {
	Node     *yaml.Node    // the actual node
	Path     []interface{} /// the path we took to get to this node
	Document uint          // the document index of this node
	Filename string
}

func (n *CandidateNode) GetKey() string {
	return fmt.Sprintf("%v - %v - %v", n.Document, n.Path, n.Node.Value)
}

func (n *CandidateNode) Copy() *CandidateNode {
	clone := &CandidateNode{}
	copier.Copy(clone, n)
	return clone
}

// updates this candidate from the given candidate node
func (n *CandidateNode) UpdateFrom(other *CandidateNode) {
	n.UpdateAttributesFrom(other)
	n.Node.Content = other.Node.Content
	n.Node.Value = other.Node.Value
}

func (n *CandidateNode) UpdateAttributesFrom(other *CandidateNode) {
	if n.Node.Kind != other.Node.Kind {
		// clear out the contents when switching to a different type
		// e.g. map to array
		n.Node.Content = make([]*yaml.Node, 0)
		n.Node.Value = ""
	}
	n.Node.Kind = other.Node.Kind
	n.Node.Tag = other.Node.Tag
	// not sure if this ever should happen here...
	// if other.Node.Style != 0 {
	// 	n.Node.Style = other.Node.Style
	// }
	n.Node.FootComment = other.Node.FootComment
	n.Node.HeadComment = other.Node.HeadComment
	n.Node.LineComment = other.Node.LineComment
}

func (n *CandidateNode) PathStackToString() string {
	return mergePathStackToString(n.Path)
}

func mergePathStackToString(pathStack []interface{}) string {
	var sb strings.Builder
	for index, path := range pathStack {
		switch path.(type) {
		case int, int64:
			// if arrayMergeStrategy == AppendArrayMergeStrategy {
			// sb.WriteString("[+]")
			// } else {
			sb.WriteString(fmt.Sprintf("[%v]", path))
			// }

		default:
			s := fmt.Sprintf("%v", path)
			var _, errParsingInt = strconv.ParseInt(s, 10, 64) // nolint

			hasSpecial := strings.Contains(s, ".") || strings.Contains(s, "[") || strings.Contains(s, "]") || strings.Contains(s, "\"")
			hasDoubleQuotes := strings.Contains(s, "\"")
			wrappingCharacterStart := "\""
			wrappingCharacterEnd := "\""
			if hasDoubleQuotes {
				wrappingCharacterStart = "("
				wrappingCharacterEnd = ")"
			}
			if hasSpecial || errParsingInt == nil {
				sb.WriteString(wrappingCharacterStart)
			}
			sb.WriteString(s)
			if hasSpecial || errParsingInt == nil {
				sb.WriteString(wrappingCharacterEnd)
			}
		}

		if index < len(pathStack)-1 {
			sb.WriteString(".")
		}
	}
	return sb.String()
}
