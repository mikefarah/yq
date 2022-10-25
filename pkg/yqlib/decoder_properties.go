package yqlib

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/magiconair/properties"
	"gopkg.in/yaml.v3"
)

type propertiesDecoder struct {
	reader   io.Reader
	finished bool
	d        DataTreeNavigator
}

func NewPropertiesDecoder() Decoder {
	return &propertiesDecoder{d: NewDataTreeNavigator(), finished: false}
}

func (dec *propertiesDecoder) Init(reader io.Reader) error {
	dec.reader = reader
	dec.finished = false
	return nil
}

func parsePropKey(key string) []interface{} {
	pathStrArray := strings.Split(key, ".")
	path := make([]interface{}, len(pathStrArray))
	for i, pathStr := range pathStrArray {
		num, err := strconv.ParseInt(pathStr, 10, 32)
		if err == nil {
			path[i] = num
		} else {
			path[i] = pathStr
		}
	}
	return path
}

func (dec *propertiesDecoder) processComment(c string) string {
	if c == "" {
		return ""
	}
	return "# " + c
}

func (dec *propertiesDecoder) applyProperty(properties *properties.Properties, context Context, key string) error {
	value, _ := properties.Get(key)
	path := parsePropKey(key)

	rhsNode := &yaml.Node{
		Value:       value,
		Tag:         "!!str",
		Kind:        yaml.ScalarNode,
		LineComment: dec.processComment(properties.GetComment(key)),
	}

	rhsNode.Tag = guessTagFromCustomType(rhsNode)

	rhsCandidateNode := &CandidateNode{
		Path: path,
		Node: rhsNode,
	}

	assignmentOp := &Operation{OperationType: assignOpType, Preferences: assignPreferences{}}

	rhsOp := &Operation{OperationType: valueOpType, CandidateNode: rhsCandidateNode}

	assignmentOpNode := &ExpressionNode{
		Operation: assignmentOp,
		LHS:       createTraversalTree(path, traversePreferences{}, false),
		RHS:       &ExpressionNode{Operation: rhsOp},
	}

	_, err := dec.d.GetMatchingNodes(context, assignmentOpNode)
	return err
}

func (dec *propertiesDecoder) Decode() (*CandidateNode, error) {
	if dec.finished {
		return nil, io.EOF
	}
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(dec.reader); err != nil {
		return nil, err
	}
	if buf.Len() == 0 {
		dec.finished = true
		return nil, io.EOF
	}
	properties, err := properties.LoadString(buf.String())
	if err != nil {
		return nil, err
	}
	properties.DisableExpansion = true

	rootMap := &CandidateNode{
		Node: &yaml.Node{
			Kind: yaml.MappingNode,
			Tag:  "!!map",
		},
	}

	context := Context{}
	context = context.SingleChildContext(rootMap)

	for _, key := range properties.Keys() {
		if err := dec.applyProperty(properties, context, key); err != nil {
			return nil, err
		}

	}
	dec.finished = true

	return &CandidateNode{
		Node: &yaml.Node{
			Kind:    yaml.DocumentNode,
			Content: []*yaml.Node{rootMap.Node},
		},
	}, nil

}
