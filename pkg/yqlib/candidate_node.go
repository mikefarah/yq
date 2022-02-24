package yqlib

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/jinzhu/copier"
	yaml "gopkg.in/yaml.v3"
)

type CandidateNode struct {
	Node   *yaml.Node     // the actual node
	Parent *CandidateNode // parent node
	Key    *yaml.Node     // node key, if this is a value from a map (or index in an array)

	LeadingContent string

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

func (n *CandidateNode) GetNiceTag() string {
	return unwrapDoc(n.Node).Tag
}

func (n *CandidateNode) GetNicePath() string {
	if n.Path != nil && len(n.Path) >= 0 {
		pathStr := make([]string, len(n.Path))
		for i, v := range n.Path {
			pathStr[i] = fmt.Sprintf("%v", v)
		}
		return strings.Join(pathStr, ".")
	}
	return ""
}

func (n *CandidateNode) AsList() *list.List {
	elMap := list.New()
	elMap.PushBack(n)
	return elMap
}

func (n *CandidateNode) CreateChildInMap(key *yaml.Node, node *yaml.Node) *CandidateNode {
	var value interface{}
	if key != nil {
		value = key.Value
	}
	return &CandidateNode{
		Node:      node,
		Path:      n.createChildPath(value),
		Parent:    n,
		Key:       key,
		Document:  n.Document,
		Filename:  n.Filename,
		FileIndex: n.FileIndex,
	}
}

func (n *CandidateNode) CreateChildInArray(index int, node *yaml.Node) *CandidateNode {
	return &CandidateNode{
		Node:      node,
		Path:      n.createChildPath(index),
		Parent:    n,
		Key:       &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", index), Tag: "!!int"},
		Document:  n.Document,
		Filename:  n.Filename,
		FileIndex: n.FileIndex,
	}
}

func (n *CandidateNode) CreateReplacement(node *yaml.Node) *CandidateNode {
	return &CandidateNode{
		Node:      node,
		Path:      n.createChildPath(nil),
		Parent:    n.Parent,
		Key:       n.Key,
		IsMapKey:  n.IsMapKey,
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
	clone.Node = deepClone(n.Node)
	return clone, nil
}

// updates this candidate from the given candidate node
func (n *CandidateNode) UpdateFrom(other *CandidateNode, prefs assignPreferences) {

	// if this is an empty map or empty array, use the style of other node.
	if (n.Node.Kind != yaml.ScalarNode && len(n.Node.Content) == 0) ||
		// if the tag has changed (e.g. from str to bool)
		(guessTagFromCustomType(n.Node) != guessTagFromCustomType(other.Node)) {
		n.Node.Style = other.Node.Style
	}

	n.Node.Content = deepCloneContent(other.Node.Content)
	n.Node.Kind = other.Node.Kind
	n.Node.Value = other.Node.Value

	n.UpdateAttributesFrom(other, prefs)

}

func (n *CandidateNode) UpdateAttributesFrom(other *CandidateNode, prefs assignPreferences) {
	log.Debug("UpdateAttributesFrom: n: %v other: %v", n.GetKey(), other.GetKey())
	if n.Node.Kind != other.Node.Kind {
		// clear out the contents when switching to a different type
		// e.g. map to array
		n.Node.Content = make([]*yaml.Node, 0)
		n.Node.Value = ""
	}
	n.Node.Kind = other.Node.Kind

	// don't clobber custom tags...
	if strings.HasPrefix(n.Node.Tag, "!!") || n.Node.Tag == "" {
		n.Node.Tag = other.Node.Tag
	}

	n.Node.Alias = other.Node.Alias

	if !prefs.DontOverWriteAnchor {
		n.Node.Anchor = other.Node.Anchor
	}

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
