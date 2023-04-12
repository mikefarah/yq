package yqlib

import (
	"container/list"
	"fmt"
	"strconv"
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

	Document uint // the document index of this node
	Filename string

	Line   int
	Column int

	FileIndex int
	// when performing op against all nodes given, this will treat all the nodes as one
	// (e.g. top level cross document merge). This property does not propagate to child nodes.
	EvaluateTogether bool
	IsMapKey         bool
}

func (n *CandidateNode) CreateChild() *CandidateNode {
	return &CandidateNode{
		Parent:    n,
		Document:  n.Document,
		Filename:  n.Filename,
		FileIndex: n.FileIndex,
	}
}

func (n *CandidateNode) GetKey() string {
	keyPrefix := ""
	if n.IsMapKey {
		keyPrefix = "key-"
	}
	key := ""
	if n.Key != nil {
		key = n.Key.Value
	}
	return fmt.Sprintf("%v%v - %v", keyPrefix, n.Document, key)
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

func (n *CandidateNode) getParsedKey() interface{} {
	if n.Key == nil {
		return nil
	}
	if n.Key.Tag == "!!str" {
		return n.Key.Value
	}
	index, err := parseInt(n.Key.Value)
	if err != nil {
		return n.Key.Value
	}
	return index

}

func (n *CandidateNode) GetPath() []interface{} {
	if n.Parent != nil {
		return append(n.Parent.GetPath(), n.getParsedKey())
	}
	key := n.getParsedKey()
	if key != nil {
		return []interface{}{key}
	}
	return make([]interface{}, 0)
}

func (n *CandidateNode) GetNicePath() string {
	var sb strings.Builder
	path := n.GetPath()
	for i, element := range path {
		elementStr := fmt.Sprintf("%v", element)
		switch element.(type) {
		case int:
			sb.WriteString("[" + elementStr + "]")
		default:
			if i == 0 {
				sb.WriteString(elementStr)
			} else if strings.ContainsRune(elementStr, '.') {
				sb.WriteString("[" + elementStr + "]")
			} else {
				sb.WriteString("." + elementStr)
			}
		}
	}
	return sb.String()
}

func (n *CandidateNode) AsList() *list.List {
	elMap := list.New()
	elMap.PushBack(n)
	return elMap
}

func (n *CandidateNode) GetValueRep() (interface{}, error) {
	// TODO: handle booleans, ints, etc
	realTag := n.guessTagFromCustomType()

	switch realTag {
	case "!!int":
		_, val, err := parseInt64(n.Value)
		return val, err
	case "!!float":
		// need to test this
		return strconv.ParseFloat(n.Value, 64)
	case "!!bool":
		return isTruthyNode(n)
	}

	return n.Value, nil
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

func (n *CandidateNode) CreateReplacement(kind Kind, tag string, value string) *CandidateNode {
	node := &CandidateNode{
		Kind:  kind,
		Tag:   tag,
		Value: value,
	}
	return n.CopyAsReplacement(node)
}

func (n *CandidateNode) CopyAsReplacement(replacement *CandidateNode) *CandidateNode {
	newCopy := replacement.Copy()
	newCopy.Parent = n.Parent
	newCopy.Key = n.Key
	newCopy.IsMapKey = n.IsMapKey
	newCopy.Document = n.Document
	newCopy.Filename = n.Filename
	newCopy.FileIndex = n.FileIndex
	return newCopy
}

func (n *CandidateNode) CreateReplacementWithDocWrappers(kind Kind, tag string, style Style) *CandidateNode {
	replacement := n.CreateReplacement(kind, tag, "")
	replacement.LeadingContent = n.LeadingContent
	replacement.TrailingContent = n.TrailingContent
	replacement.Style = style
	return replacement
}

func (n *CandidateNode) CopyChildren() []*CandidateNode {
	log.Debug("n? %v", n)
	log.Debug("n.Content %v", n.Content)
	log.Debug("n.Content %v", len(n.Content))
	clonedKids := make([]*CandidateNode, len(n.Content))
	log.Debug("created clone")
	for i, child := range n.Content {
		clonedKids[i] = child.Copy()
	}
	log.Debug("finishing clone")
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

	var copyKey *CandidateNode
	if n.Key != nil {
		copyKey = n.Key.Copy()
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
		Key:    copyKey,

		LeadingContent:  n.LeadingContent,
		TrailingContent: n.TrailingContent,

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
	log.Debug("UpdateAttributesFrom: n: %v other: %v", n.Key.Value, other.Key.Value)
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
