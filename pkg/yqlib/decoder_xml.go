package yqlib

import (
	"io"

	xj "github.com/basgys/goxml2json"
	yaml "gopkg.in/yaml.v3"
)

type xmlDecoder struct {
	reader          io.Reader
	attributePrefix string
	contentPrefix   string
	finished        bool
}

func NewXmlDecoder(reader io.Reader, attributePrefix string, contentPrefix string) Decoder {
	return &xmlDecoder{reader: reader, attributePrefix: attributePrefix, contentPrefix: contentPrefix, finished: false}
}

func (dec *xmlDecoder) createSequence(nodes xj.Nodes) (*yaml.Node, error) {
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

func (dec *xmlDecoder) createMap(n *xj.Node) (*yaml.Node, error) {
	yamlNode := &yaml.Node{Kind: yaml.MappingNode}

	if len(n.Data) > 0 {
		label := dec.contentPrefix + "content"
		yamlNode.Content = append(yamlNode.Content, createScalarNode(label, label), createScalarNode(n.Data, n.Data))
	}

	for label, children := range n.Children {
		labelNode := createScalarNode(label, label)
		var valueNode *yaml.Node
		var err error
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

func (dec *xmlDecoder) convertToYamlNode(n *xj.Node) (*yaml.Node, error) {
	if n.IsComplex() {
		return dec.createMap(n)
	}
	return createScalarNode(n.Data, n.Data), nil
}

func (dec *xmlDecoder) Decode(rootYamlNode *yaml.Node) error {
	if dec.finished {
		return io.EOF
	}
	root := &xj.Node{}
	// cant use xj - it doesn't keep map order.
	err := xj.NewDecoder(dec.reader).Decode(root)

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
