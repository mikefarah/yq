package treeops

import "gopkg.in/yaml.v3"

func BuildCandidateNodeFrom(p *PathElement) *CandidateNode {
	var node yaml.Node = yaml.Node{Kind: yaml.ScalarNode}
	node.Value = p.StringValue

	switch p.Value.(type) {
	case float32, float64:
		node.Tag = "!!float"
	case int, int64, int32:
		node.Tag = "!!int"
	case bool:
		node.Tag = "!!bool"
	case string:
		node.Tag = "!!str"
	}
	return &CandidateNode{Node: &node}
}
