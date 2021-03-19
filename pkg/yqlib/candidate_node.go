package yqlib

import (
	"fmt"

	"github.com/jinzhu/copier"
	yaml "gopkg.in/yaml.v3"
)

type CandidateNode struct {
	Node      *yaml.Node    // the actual node
	Path      []interface{} /// the path we took to get to this node
	Document  uint          // the document index of this node
	Filename  string
	FileIndex int
	// when performing op against all nodes given, this will treat all the nodes as one
	// (e.g. top level cross document merge). This property does not propegate to child nodes.
	EvaluateTogether bool
	IsMapKey         bool
}

func (n *CandidateNode) GetKey() string {
	keyPrefix := ""
	if n.IsMapKey {
		keyPrefix = "key-"
	}
	return fmt.Sprintf("%v%v - %v", keyPrefix, n.Document, n.Path)
}

func (n *CandidateNode) CreateChild(path interface{}, node *yaml.Node) *CandidateNode {
	return &CandidateNode{
		Node:      node,
		Path:      n.createChildPath(path),
		Document:  n.Document,
		Filename:  n.Filename,
		FileIndex: n.FileIndex,
	}
}

func (n *CandidateNode) createChildPath(path interface{}) []interface{} {
	if path == nil {
		newPath := make([]interface{}, len(n.Path))
		copy(newPath, n.Path)
		return newPath
	}

	//don't use append as they may actually modify the path of the orignal node!
	newPath := make([]interface{}, len(n.Path)+1)
	copy(newPath, n.Path)
	newPath[len(n.Path)] = path
	return newPath
}

func (n *CandidateNode) Copy() (*CandidateNode, error) {
	clone := &CandidateNode{}
	err := copier.Copy(clone, n)
	if err != nil {
		return nil, err
	}
	return clone, nil
}

// updates this candidate from the given candidate node
func (n *CandidateNode) UpdateFrom(other *CandidateNode) {

	n.UpdateAttributesFrom(other)
	n.Node.Content = other.Node.Content
	n.Node.Value = other.Node.Value
}

func (n *CandidateNode) UpdateAttributesFrom(other *CandidateNode) {
	log.Debug("UpdateAttributesFrom: n: %v other: %v", n.GetKey(), other.GetKey())
	if n.Node.Kind != other.Node.Kind {
		// clear out the contents when switching to a different type
		// e.g. map to array
		n.Node.Content = make([]*yaml.Node, 0)
		n.Node.Value = ""
	}
	n.Node.Kind = other.Node.Kind
	n.Node.Tag = other.Node.Tag
	n.Node.Alias = other.Node.Alias
	n.Node.Anchor = other.Node.Anchor

	// merge will pickup the style of the new thing
	// when autocreating nodes
	if n.Node.Style == 0 {
		n.Node.Style = other.Node.Style
	}

	if other.Node.FootComment != "" {
		n.Node.FootComment = other.Node.FootComment
	}
	if other.Node.HeadComment != "" {
		n.Node.HeadComment = other.Node.HeadComment
	}
	if other.Node.LineComment != "" {
		n.Node.LineComment = other.Node.LineComment
	}
}
