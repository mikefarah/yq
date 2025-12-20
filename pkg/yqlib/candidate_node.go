package yqlib

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

type Kind uint32

const (
	SequenceNode Kind = 1 << iota
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

type NodeInfo struct {
	Kind        string      `yaml:"kind"`
	Style       string      `yaml:"style,omitempty"`
	Anchor      string      `yaml:"anchor,omitempty"`
	Tag         string      `yaml:"tag,omitempty"`
	HeadComment string      `yaml:"headComment,omitempty"`
	LineComment string      `yaml:"lineComment,omitempty"`
	FootComment string      `yaml:"footComment,omitempty"`
	Value       string      `yaml:"value,omitempty"`
	Line        int         `yaml:"line,omitempty"`
	Column      int         `yaml:"column,omitempty"`
	Content     []*NodeInfo `yaml:"content,omitempty"`
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

	LeadingContent string

	document uint // the document index of this node
	filename string

	Line   int
	Column int

	fileIndex int
	// when performing op against all nodes given, this will treat all the nodes as one
	// (e.g. top level cross document merge). This property does not propagate to child nodes.
	EvaluateTogether bool
	IsMapKey         bool
	// For formats like HCL and TOML: indicates that child entries should be emitted as separate blocks/tables
	// rather than consolidated into nested mappings (default behaviour)
	EncodeSeparate bool
}

func (n *CandidateNode) CreateChild() *CandidateNode {
	return &CandidateNode{
		Parent: n,
	}
}

func (n *CandidateNode) SetDocument(idx uint) {
	n.document = idx
}

func (n *CandidateNode) GetDocument() uint {
	// defer to parent
	if n.Parent != nil {
		return n.Parent.GetDocument()
	}
	return n.document
}

func (n *CandidateNode) SetFilename(name string) {
	n.filename = name
}

func (n *CandidateNode) GetFilename() string {
	if n.Parent != nil {
		return n.Parent.GetFilename()
	}
	return n.filename
}

func (n *CandidateNode) SetFileIndex(idx int) {
	n.fileIndex = idx
}

func (n *CandidateNode) GetFileIndex() int {
	if n.Parent != nil {
		return n.Parent.GetFileIndex()
	}
	return n.fileIndex
}

func (n *CandidateNode) GetKey() string {
	keyPrefix := ""
	if n.IsMapKey {
		keyPrefix = fmt.Sprintf("key-%v-", n.Value)
	}
	key := ""
	if n.Key != nil {
		key = n.Key.Value
	}
	return fmt.Sprintf("%v%v - %v", keyPrefix, n.GetDocument(), key)
}

func (n *CandidateNode) getParsedKey() interface{} {
	if n.IsMapKey {
		return n.Value
	}
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

func (n *CandidateNode) FilterMapContentByKey(keyPredicate func(*CandidateNode) bool) []*CandidateNode {
	var result []*CandidateNode
	for index := 0; index < len(n.Content); index = index + 2 {
		keyNode := n.Content[index]
		valueNode := n.Content[index+1]
		if keyPredicate(keyNode) {
			result = append(result, keyNode, valueNode)
		}
	}
	return result
}

func (n *CandidateNode) GetPath() []interface{} {
	key := n.getParsedKey()
	if n.Parent != nil && key != nil {
		return append(n.Parent.GetPath(), key)
	}

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

func (n *CandidateNode) SetParent(parent *CandidateNode) {
	n.Parent = parent
}

type ValueVisitor func(*CandidateNode) error

func (n *CandidateNode) VisitValues(visitor ValueVisitor) error {
	switch n.Kind {
	case MappingNode:
		for i := 1; i < len(n.Content); i = i + 2 {
			if err := visitor(n.Content[i]); err != nil {
				return err
			}
		}
	case SequenceNode:
		for i := 0; i < len(n.Content); i = i + 1 {
			if err := visitor(n.Content[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *CandidateNode) CanVisitValues() bool {
	return n.Kind == MappingNode || n.Kind == SequenceNode
}

func (n *CandidateNode) AddKeyValueChild(rawKey *CandidateNode, rawValue *CandidateNode) (*CandidateNode, *CandidateNode) {
	key := rawKey.Copy()
	key.SetParent(n)
	key.IsMapKey = true

	value := rawValue.Copy()
	value.SetParent(n)
	value.IsMapKey = false // force this, incase we are creating a value from a key
	value.Key = key

	n.Content = append(n.Content, key, value)
	return key, value
}

func (n *CandidateNode) AddChild(rawChild *CandidateNode) {
	value := rawChild.Copy()
	value.SetParent(n)
	value.IsMapKey = false

	index := len(n.Content)
	keyNode := createScalarNode(index, fmt.Sprintf("%v", index))
	keyNode.SetParent(n)
	value.Key = keyNode

	n.Content = append(n.Content, value)
}

func (n *CandidateNode) AddChildren(children []*CandidateNode) {
	if n.Kind == MappingNode {
		for i := 0; i < len(children); i += 2 {
			key := children[i]
			value := children[i+1]
			n.AddKeyValueChild(key, value)
		}

	} else {
		for _, rawChild := range children {
			n.AddChild(rawChild)
		}
	}
}

func (n *CandidateNode) GetValueRep() (interface{}, error) {
	log.Debugf("GetValueRep for %v value: %v", n.GetNicePath(), n.Value)
	realTag := n.guessTagFromCustomType()

	switch realTag {
	case "!!int":
		_, val, err := parseInt64(n.Value)
		return val, err
	case "!!float":
		// need to test this
		return strconv.ParseFloat(n.Value, 64)
	case "!!bool":
		return isTruthyNode(n), nil
	case "!!null":
		return nil, nil
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
	guessedTag := dataBucket.Tag
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

	if n.IsMapKey {
		newCopy.Key = n
	} else {
		newCopy.Key = n.Key
	}

	return newCopy
}

func (n *CandidateNode) CreateReplacementWithComments(kind Kind, tag string, style Style) *CandidateNode {
	replacement := n.CreateReplacement(kind, tag, "")
	replacement.LeadingContent = n.LeadingContent
	replacement.HeadComment = n.HeadComment
	replacement.LineComment = n.LineComment
	replacement.FootComment = n.FootComment
	replacement.Style = style
	return replacement
}

func (n *CandidateNode) Copy() *CandidateNode {
	return n.doCopy(true)
}

func (n *CandidateNode) CopyWithoutContent() *CandidateNode {
	return n.doCopy(false)
}

func (n *CandidateNode) doCopy(cloneContent bool) *CandidateNode {
	var content []*CandidateNode

	var copyKey *CandidateNode
	if n.Key != nil {
		copyKey = n.Key.Copy()
	}

	clone := &CandidateNode{
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

		LeadingContent: n.LeadingContent,

		document:  n.document,
		filename:  n.filename,
		fileIndex: n.fileIndex,

		Line:   n.Line,
		Column: n.Column,

		EvaluateTogether: n.EvaluateTogether,
		IsMapKey:         n.IsMapKey,

		EncodeSeparate: n.EncodeSeparate,
	}

	if cloneContent {
		clone.AddChildren(n.Content)
	}

	return clone
}

// updates this candidate from the given candidate node
func (n *CandidateNode) UpdateFrom(other *CandidateNode, prefs assignPreferences) {
	if n == other {
		log.Debugf("UpdateFrom, no need to update from myself.")
		return
	}
	// if this is an empty map or empty array, use the style of other node.
	if (n.Kind != ScalarNode && len(n.Content) == 0) ||
		// if the tag has changed (e.g. from str to bool)
		(n.guessTagFromCustomType() != other.guessTagFromCustomType()) {
		n.Style = other.Style
	}

	n.Content = make([]*CandidateNode, 0)
	n.Kind = other.Kind
	n.AddChildren(other.Content)

	n.Value = other.Value

	n.UpdateAttributesFrom(other, prefs)

}

func (n *CandidateNode) UpdateAttributesFrom(other *CandidateNode, prefs assignPreferences) {
	log.Debug("UpdateAttributesFrom: n: %v other: %v", NodeToString(n), NodeToString(other))
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

	// Preserve EncodeSeparate flag for format-specific encoding hints
	n.EncodeSeparate = other.EncodeSeparate

	// merge will pickup the style of the new thing
	// when autocreating nodes

	if n.Style == 0 {
		n.Style = other.Style
	}

	if other.FootComment != "" {
		n.FootComment = other.FootComment
	}
	if other.HeadComment != "" {
		n.HeadComment = other.HeadComment
	}
	if other.LineComment != "" {
		n.LineComment = other.LineComment
	}
}

func (n *CandidateNode) ConvertToNodeInfo() *NodeInfo {
	info := &NodeInfo{
		Kind:        kindToString(n.Kind),
		Style:       styleToString(n.Style),
		Anchor:      n.Anchor,
		Tag:         n.Tag,
		HeadComment: n.HeadComment,
		LineComment: n.LineComment,
		FootComment: n.FootComment,
		Value:       n.Value,
		Line:        n.Line,
		Column:      n.Column,
	}
	if len(n.Content) > 0 {
		info.Content = make([]*NodeInfo, len(n.Content))
		for i, child := range n.Content {
			info.Content[i] = child.ConvertToNodeInfo()
		}
	}
	return info
}

// Helper functions to convert Kind and Style to string for NodeInfo
func kindToString(k Kind) string {
	switch k {
	case SequenceNode:
		return "SequenceNode"
	case MappingNode:
		return "MappingNode"
	case ScalarNode:
		return "ScalarNode"
	case AliasNode:
		return "AliasNode"
	default:
		return "Unknown"
	}
}

func styleToString(s Style) string {
	var styles []string
	if s&TaggedStyle != 0 {
		styles = append(styles, "TaggedStyle")
	}
	if s&DoubleQuotedStyle != 0 {
		styles = append(styles, "DoubleQuotedStyle")
	}
	if s&SingleQuotedStyle != 0 {
		styles = append(styles, "SingleQuotedStyle")
	}
	if s&LiteralStyle != 0 {
		styles = append(styles, "LiteralStyle")
	}
	if s&FoldedStyle != 0 {
		styles = append(styles, "FoldedStyle")
	}
	if s&FlowStyle != 0 {
		styles = append(styles, "FlowStyle")
	}
	return strings.Join(styles, ",")
}
