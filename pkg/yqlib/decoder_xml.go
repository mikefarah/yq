package yqlib

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"
	"unicode"

	"golang.org/x/net/html/charset"
	yaml "gopkg.in/yaml.v3"
)

type xmlDecoder struct {
	reader          io.Reader
	attributePrefix string
	contentName     string
	strictMode      bool
	finished        bool
}

func NewXMLDecoder(attributePrefix string, contentName string, strictMode bool) Decoder {
	if contentName == "" {
		contentName = "content"
	}
	return &xmlDecoder{attributePrefix: attributePrefix, contentName: contentName, finished: false, strictMode: strictMode}
}

func (dec *xmlDecoder) Init(reader io.Reader) {
	dec.reader = reader
	dec.finished = false
}

func (dec *xmlDecoder) createSequence(nodes []*xmlNode) (*yaml.Node, error) {
	yamlNode := &yaml.Node{Kind: yaml.SequenceNode}
	for _, child := range nodes {
		yamlChild, err := dec.convertToYamlNode(child)
		if err != nil {
			return nil, err
		}
		yamlNode.Content = append(yamlNode.Content, yamlChild)
	}

	return yamlNode, nil
}

func (dec *xmlDecoder) processComment(c string) string {
	if c == "" {
		return ""
	}
	return "#" + strings.TrimRight(c, " ")
}

func (dec *xmlDecoder) createMap(n *xmlNode) (*yaml.Node, error) {
	log.Debug("createMap: headC: %v, footC: %v", n.HeadComment, n.FootComment)
	yamlNode := &yaml.Node{Kind: yaml.MappingNode}

	if len(n.Data) > 0 {
		label := dec.contentName
		labelNode := createScalarNode(label, label)
		labelNode.HeadComment = dec.processComment(n.HeadComment)
		labelNode.FootComment = dec.processComment(n.FootComment)
		yamlNode.Content = append(yamlNode.Content, labelNode, createScalarNode(n.Data, n.Data))
	}

	for i, keyValuePair := range n.Children {
		label := keyValuePair.K
		children := keyValuePair.V
		labelNode := createScalarNode(label, label)
		var valueNode *yaml.Node
		var err error

		if i == 0 {
			labelNode.HeadComment = dec.processComment(n.HeadComment)

		}

		// if i == len(n.Children)-1 {
		labelNode.FootComment = dec.processComment(keyValuePair.FootComment)
		// }

		log.Debug("len of children in %v is %v", label, len(children))
		if len(children) > 1 {
			valueNode, err = dec.createSequence(children)
			if err != nil {
				return nil, err
			}
		} else {
			// comment hack for maps of scalars
			// if the value is a scalar, the head comment of the scalar needs to go on the key?
			// add tests for <z/> as well as multiple <ds> of inputXmlWithComments > yaml
			if len(children[0].Children) == 0 && children[0].HeadComment != "" {
				labelNode.HeadComment = labelNode.HeadComment + "\n" + strings.TrimSpace(children[0].HeadComment)
				children[0].HeadComment = ""
			}
			valueNode, err = dec.convertToYamlNode(children[0])
			if err != nil {
				return nil, err
			}
		}
		yamlNode.Content = append(yamlNode.Content, labelNode, valueNode)
	}

	return yamlNode, nil
}

func (dec *xmlDecoder) convertToYamlNode(n *xmlNode) (*yaml.Node, error) {
	if len(n.Children) > 0 {
		return dec.createMap(n)
	}
	scalar := createScalarNode(n.Data, n.Data)
	if n.Data == "" {
		scalar = createScalarNode(nil, "")
	}
	log.Debug("scalar headC: %v, footC: %v", n.HeadComment, n.FootComment)
	scalar.HeadComment = dec.processComment(n.HeadComment)
	scalar.LineComment = dec.processComment(n.LineComment)
	scalar.FootComment = dec.processComment(n.FootComment)

	return scalar, nil
}

func (dec *xmlDecoder) Decode(rootYamlNode *yaml.Node) error {
	if dec.finished {
		return io.EOF
	}
	root := &xmlNode{}
	// cant use xj - it doesn't keep map order.
	err := dec.decodeXML(root)

	if err != nil {
		return err
	}
	firstNode, err := dec.convertToYamlNode(root)

	if err != nil {
		return err
	} else if firstNode.Tag == "!!null" {
		dec.finished = true
		return io.EOF
	}
	rootYamlNode.Kind = yaml.DocumentNode
	rootYamlNode.Content = []*yaml.Node{firstNode}
	dec.finished = true
	return nil
}

type xmlNode struct {
	Children    []*xmlChildrenKv
	HeadComment string
	FootComment string
	LineComment string
	Data        string
}

type xmlChildrenKv struct {
	K           string
	V           []*xmlNode
	FootComment string
}

// AddChild appends a node to the list of children
func (n *xmlNode) AddChild(s string, c *xmlNode) {

	if n.Children == nil {
		n.Children = make([]*xmlChildrenKv, 0)
	}
	log.Debug("looking for %s", s)
	// see if we can find an existing entry to add to
	for _, childEntry := range n.Children {
		if childEntry.K == s {
			log.Debug("found it, appending an entry%s", s)
			childEntry.V = append(childEntry.V, c)
			log.Debug("yay len of children in %v is %v", s, len(childEntry.V))
			return
		}
	}
	log.Debug("not there, making a new one %s", s)
	n.Children = append(n.Children, &xmlChildrenKv{K: s, V: []*xmlNode{c}})
}

type element struct {
	parent *element
	n      *xmlNode
	label  string
	state  string
}

// this code is heavily based on https://github.com/basgys/goxml2json
// main changes are to decode into a structure that preserves the original order
// of the map keys.
func (dec *xmlDecoder) decodeXML(root *xmlNode) error {
	xmlDec := xml.NewDecoder(dec.reader)
	xmlDec.Strict = dec.strictMode
	// That will convert the charset if the provided XML is non-UTF-8
	xmlDec.CharsetReader = charset.NewReaderLabel

	// Create first element from the root node
	elem := &element{
		parent: nil,
		n:      root,
	}

	for {
		t, e := xmlDec.Token()
		if e != nil && !errors.Is(e, io.EOF) {
			return e
		}
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			log.Debug("start element %v", se.Name.Local)
			elem.state = "started"
			// Build new a new current element and link it to its parent
			elem = &element{
				parent: elem,
				n:      &xmlNode{},
				label:  se.Name.Local,
			}

			// Extract attributes as children
			for _, a := range se.Attr {
				elem.n.AddChild(dec.attributePrefix+a.Name.Local, &xmlNode{Data: a.Value})
			}
		case xml.CharData:
			// Extract XML data (if any)
			elem.n.Data = trimNonGraphic(string(se))
			if elem.n.Data != "" {
				elem.state = "chardata"
				log.Debug("chardata [%v] for %v", elem.n.Data, elem.label)
			}
		case xml.EndElement:
			log.Debug("end element %v", elem.label)
			elem.state = "finished"
			// And add it to its parent list
			if elem.parent != nil {
				elem.parent.n.AddChild(elem.label, elem.n)
			}

			// Then change the current element to its parent
			elem = elem.parent
		case xml.Comment:

			commentStr := string(xml.CharData(se))
			if elem.state == "started" {
				applyFootComment(elem, commentStr)

			} else if elem.state == "chardata" {
				log.Debug("got a line comment for (%v) %v: [%v]", elem.state, elem.label, commentStr)
				elem.n.LineComment = joinFilter([]string{elem.n.LineComment, commentStr})
			} else {
				log.Debug("got a head comment for (%v) %v: [%v]", elem.state, elem.label, commentStr)
				elem.n.HeadComment = joinFilter([]string{elem.n.HeadComment, commentStr})
			}

		}
	}

	return nil
}

func applyFootComment(elem *element, commentStr string) {

	// first lets try to put the comment on the last child
	if len(elem.n.Children) > 0 {
		lastChildIndex := len(elem.n.Children) - 1
		childKv := elem.n.Children[lastChildIndex]
		log.Debug("got a foot comment for %v: [%v]", childKv.K, commentStr)
		childKv.FootComment = joinFilter([]string{elem.n.FootComment, commentStr})
	} else {
		log.Debug("got a foot comment for %v: [%v]", elem.label, commentStr)
		elem.n.FootComment = joinFilter([]string{elem.n.FootComment, commentStr})
	}
}

func joinFilter(rawStrings []string) string {
	stringsToJoin := make([]string, 0)
	for _, str := range rawStrings {
		if str != "" {
			stringsToJoin = append(stringsToJoin, str)
		}
	}
	return strings.Join(stringsToJoin, " ")
}

// trimNonGraphic returns a slice of the string s, with all leading and trailing
// non graphic characters and spaces removed.
//
// Graphic characters include letters, marks, numbers, punctuation, symbols,
// and spaces, from categories L, M, N, P, S, Zs.
// Spacing characters are set by category Z and property Pattern_White_Space.
func trimNonGraphic(s string) string {
	if s == "" {
		return s
	}

	var first *int
	var last int
	for i, r := range []rune(s) {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) {
			continue
		}

		if first == nil {
			f := i // copy i
			first = &f
			last = i
		} else {
			last = i
		}
	}

	// If first is nil, it means there are no graphic characters
	if first == nil {
		return ""
	}

	return string([]rune(s)[*first : last+1])
}
