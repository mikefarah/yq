package treeops

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type CandidateNode struct {
	Node     *yaml.Node    // the actual node
	Path     []interface{} /// the path we took to get to this node
	Document uint          // the document index of this node
}

func (n *CandidateNode) GetKey() string {
	return fmt.Sprintf("%v - %v", n.Document, n.Path)
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
