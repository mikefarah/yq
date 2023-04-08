package yqlib

import (
	"container/list"
	"fmt"
	"strings"
)

type Kind uint32

const (
	DocumentNode Kind = 1 << iota
	SequenceNode
	MappingNode
	ScalarNode
	AliasNode
)

type Style uint32

const (
	TaggedStyle Style = 1 << iota
	DoubleQuotedStyle
	SingleQuotedStyle
	LiteralStyle
	FoldedStyle
	FlowStyle
)

func createIntegerScalarNode(num int) *CandidateNode {
	return &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!int",
		Value: fmt.Sprintf("%v", num),
	}
}

func createStringScalarNode(stringValue string) *CandidateNode {
	var node = &CandidateNode{Kind: ScalarNode}
	node.Value = stringValue
	node.Tag = "!!str"
	return node
}

func createScalarNode(value interface{}, stringValue string) *CandidateNode {
	var node = &CandidateNode{Kind: ScalarNode}
	node.Value = stringValue

	switch value.(type) {
	case float32, float64:
		node.Tag = "!!float"
	case int, int64, int32:
		node.Tag = "!!int"
	case bool:
		node.Tag = "!!bool"
	case string:
		node.Tag = "!!str"
	case nil:
		node.Tag = "!!null"
	}
	return node
}

type CandidateNode struct {
	Kind  Kind
	Style Style

	Tag     string
	Value   string
	Anchor  string
	Alias   *CandidateNode
	Content []*CandidateNode

	HeadComment string
	LineComment string
	FootComment string

	Parent *CandidateNode // parent node
	Key    *CandidateNode // node key, if this is a value from a map (or index in an array)

	LeadingContent  string
	TrailingContent string

	Path     []interface{} /// the path we took to get to this node
	Document uint          // the document index of this node
	Filename string

	Line   int
	Column int

	FileIndex int
	// when performing op against all nodes given, this will treat all the nodes as one
	// (e.g. top level cross document merge). This property does not propagate to child nodes.
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

func (n *CandidateNode) unwrapDocument() *CandidateNode {
	if n.Kind == DocumentNode {
		return n.Content[0]
	}
	return n
}

func (n *CandidateNode) GetNiceTag() string {
	return n.unwrapDocument().Tag
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

func (n *CandidateNode) guessTagFromCustomType() string {
	if strings.HasPrefix(n.Tag, "!!") {
		return n.Tag
	} else if n.Value == "" {
		log.Debug("guessTagFromCustomType: node has no value to guess the type with")
		return n.Tag
	}
	dataBucket, errorReading := parseSnippet(n.Value)

	if errorReading != nil {
		log.Debug("guessTagFromCustomType: could not guess underlying tag type %v", errorReading)
		return n.Tag
	}
	guessedTag := dataBucket.unwrapDocument().Tag
	log.Info("im guessing the tag %v is a %v", n.Tag, guessedTag)
	return guessedTag
}

// func (n *CandidateNode) CreateChildInMap(key *CandidateNode) *CandidateNode {
// 	var value interface{}
// 	if key != nil {
// 		value = key.Value
// 	}
// 	return &CandidateNode{
// 		Path:   n.createChildPath(value),
// 		Parent: n,
// 		Key:    key,

// 		Document:  n.Document,
// 		Filename:  n.Filename,
// 		FileIndex: n.FileIndex,
// 	}
// }

// func (n *CandidateNode) CreateChildInArray(index int) *CandidateNode {
// 	return &CandidateNode{
// 		Path:      n.createChildPath(index),
// 		Parent:    n,
// 		Key:       createIntegerScalarNode(index),
// 		Document:  n.Document,
// 		Filename:  n.Filename,
// 		FileIndex: n.FileIndex,
// 	}
// }

func (n *CandidateNode) CreateReplacement() *CandidateNode {
	return &CandidateNode{
		Path:      n.createChildPath(nil),
		Parent:    n.Parent,
		Key:       n.Key,
		IsMapKey:  n.IsMapKey,
		Document:  n.Document,
		Filename:  n.Filename,
		FileIndex: n.FileIndex,
	}
}

// func (n *CandidateNode) CreateReplacementWithDocWrappers(node *yaml.Node) *CandidateNode {
// 	replacement := n.CreateReplacement(node)
// 	replacement.LeadingContent = n.LeadingContent
// 	replacement.TrailingContent = n.TrailingContent
// 	return replacement
// }

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

func (n *CandidateNode) CopyChildren() []*CandidateNode {
	clonedKids := make([]*CandidateNode, len(n.Content))
	for i, child := range n.Content {
		clonedKids[i] = child.Copy()
	}
	return clonedKids
}

func (n *CandidateNode) Copy() *CandidateNode {
	return n.doCopy(true)
}

func (n *CandidateNode) CopyWithoutContent() *CandidateNode {
	return n.doCopy(false)
}

func (n *CandidateNode) doCopy(cloneContent bool) *CandidateNode {
	var content []*CandidateNode
	if cloneContent {
		content = n.CopyChildren()
	}

	return &CandidateNode{
		Kind:  n.Kind,
		Style: n.Style,

		Tag:    n.Tag,
		Value:  n.Value,
		Anchor: n.Anchor,

		// ok not to clone this,
		// as its a reference to somewhere else.
		Alias:   n.Alias,
		Content: content,

		HeadComment: n.HeadComment,
		LineComment: n.LineComment,
		FootComment: n.FootComment,

		Parent: n.Parent,
		Key:    n.Key.Copy(),

		LeadingContent:  n.LeadingContent,
		TrailingContent: n.TrailingContent,

		Path:     n.Path,
		Document: n.Document,
		Filename: n.Filename,

		Line:   n.Line,
		Column: n.Column,

		FileIndex:        n.FileIndex,
		EvaluateTogether: n.EvaluateTogether,
		IsMapKey:         n.IsMapKey,
	}
}

// updates this candidate from the given candidate node
func (n *CandidateNode) UpdateFrom(other *CandidateNode, prefs assignPreferences) {

	// if this is an empty map or empty array, use the style of other node.
	if (n.Kind != ScalarNode && len(n.Content) == 0) ||
		// if the tag has changed (e.g. from str to bool)
		(n.guessTagFromCustomType() != other.guessTagFromCustomType()) {
		n.Style = other.Style
	}

	n.Content = other.CopyChildren()
	n.Kind = other.Kind
	n.Value = other.Value

	n.UpdateAttributesFrom(other, prefs)

}

func (n *CandidateNode) UpdateAttributesFrom(other *CandidateNode, prefs assignPreferences) {
	log.Debug("UpdateAttributesFrom: n: %v other: %v", n.GetKey(), other.GetKey())
	if n.Kind != other.Kind {
		// clear out the contents when switching to a different type
		// e.g. map to array
		n.Content = make([]*CandidateNode, 0)
		n.Value = ""
	}
	n.Kind = other.Kind

	// don't clobber custom tags...
	if prefs.ClobberCustomTags || strings.HasPrefix(n.Tag, "!!") || n.Tag == "" {
		n.Tag = other.Tag
	}

	n.Alias = other.Alias

	if !prefs.DontOverWriteAnchor {
		n.Anchor = other.Anchor
	}

	// merge will pickup the style of the new thing
	// when autocreating nodes

	if n.Style == 0 {
		n.Style = other.Style
	}

	if other.FootComment != "" {
		n.FootComment = other.FootComment
	}
	if other.TrailingContent != "" {
		n.TrailingContent = other.TrailingContent
	}
	if other.HeadComment != "" {
		n.HeadComment = other.HeadComment
	}
	if other.LineComment != "" {
		n.LineComment = other.LineComment
	}
}
