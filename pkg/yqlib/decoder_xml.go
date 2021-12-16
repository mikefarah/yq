package yqlib

import (
	"encoding/xml"
	"io"
	"unicode"

	"golang.org/x/net/html/charset"
	yaml "gopkg.in/yaml.v3"
)

type xmlDecoder struct {
	reader          io.Reader
	attributePrefix string
	contentPrefix   string
	finished        bool
}

func NewXmlDecoder(reader io.Reader, attributePrefix string, contentPrefix string) Decoder {
	return &xmlDecoder{reader: reader, attributePrefix: attributePrefix, contentPrefix: "c", finished: false}
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

func (dec *xmlDecoder) createMap(n *xmlNode) (*yaml.Node, error) {
	yamlNode := &yaml.Node{Kind: yaml.MappingNode, HeadComment: n.Comment}

	if len(n.Data) > 0 {
		label := dec.contentPrefix
		yamlNode.Content = append(yamlNode.Content, createScalarNode(label, label), createScalarNode(n.Data, n.Data))
	}

	for _, keyValuePair := range n.Children {
		label := keyValuePair.K
		children := keyValuePair.V
		labelNode := createScalarNode(label, label)
		var valueNode *yaml.Node
		var err error
		log.Debug("len of children in %v is %v", label, len(children))
		if len(children) > 1 {
			valueNode, err = dec.createSequence(children)
			if err != nil {
				return nil, err
			}
		} else {
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
	scalar.HeadComment = n.Comment
	return scalar, nil
}

func (dec *xmlDecoder) Decode(rootYamlNode *yaml.Node) error {
	if dec.finished {
		return io.EOF
	}
	root := &xmlNode{}
	// cant use xj - it doesn't keep map order.
	err := dec.decodeXml(root)

	if err != nil {
		return err
	}
	firstNode, err := dec.convertToYamlNode(root)

	if err != nil {
		return err
	}
	rootYamlNode.Kind = yaml.DocumentNode
	rootYamlNode.Content = []*yaml.Node{firstNode}
	dec.finished = true
	return nil
}

type xmlNode struct {
	Children []*xmlChildrenKv
	Comment  string
	Data     string
}

type xmlChildrenKv struct {
	K string
	V []*xmlNode
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
}

// this code is heavily based on https://github.com/basgys/goxml2json
// main changes are to decode into a structure that preserves the original order
// of the map keys.
func (dec *xmlDecoder) decodeXml(root *xmlNode) error {
	xmlDec := xml.NewDecoder(dec.reader)

	// That will convert the charset if the provided XML is non-UTF-8
	xmlDec.CharsetReader = charset.NewReaderLabel

	// Create first element from the root node
	elem := &element{
		parent: nil,
		n:      root,
	}

	for {
		t, _ := xmlDec.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
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
			elem.n.Data = trimNonGraphic(string(xml.CharData(se)))
		case xml.EndElement:
			// And add it to its parent list
			if elem.parent != nil {
				elem.parent.n.AddChild(elem.label, elem.n)
			}

			// Then change the current element to its parent
			elem = elem.parent
		case xml.Comment:
			elem.n.Comment = trimNonGraphic(string(xml.CharData(se)))
		}
	}

	return nil
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
