//go:build !yq_noprops

package yqlib

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/magiconair/properties"
)

type propertiesDecoder struct {
	reader   io.Reader
	finished bool
	d        DataTreeNavigator
	prefs    PropertiesPreferences
}

func NewPropertiesDecoder() Decoder {
	return &propertiesDecoder{d: NewDataTreeNavigator(), finished: false, prefs: ConfiguredPropertiesPreferences.Copy()}
}

func (dec *propertiesDecoder) Init(reader io.Reader) error {
	dec.reader = reader
	dec.finished = false
	return nil
}

func parsePropKey(key string, prefs PropertiesPreferences) []interface{} {
	pathStrArray := strings.Split(key, ".")
	path := make([]interface{}, 0, len(pathStrArray))
	for _, pathStr := range pathStrArray {
		path = appendPropKeySegment(path, pathStr, prefs.UseArrayBrackets)
	}
	return path
}

func appendPropKeySegment(path []interface{}, segment string, useArrayBrackets bool) []interface{} {
	if useArrayBrackets && strings.Contains(segment, "[") {
		bracketPath, ok := parsePropKeyArrayBracketSegment(segment)
		if ok {
			return append(path, bracketPath...)
		}
	}

	num, err := strconv.ParseInt(segment, 10, 32)
	if err == nil {
		return append(path, num)
	}
	return append(path, segment)
}

func parsePropKeyArrayBracketSegment(segment string) ([]interface{}, bool) {
	path := []interface{}{}
	bracketIndex := strings.Index(segment, "[")
	if bracketIndex > 0 {
		path = append(path, segment[:bracketIndex])
	}

	remaining := segment[bracketIndex:]
	for remaining != "" {
		if !strings.HasPrefix(remaining, "[") {
			return nil, false
		}
		closingBracket := strings.Index(remaining, "]")
		if closingBracket < 0 {
			return nil, false
		}
		arrayIndex, err := strconv.ParseInt(remaining[1:closingBracket], 10, 32)
		if err != nil {
			return nil, false
		}
		path = append(path, arrayIndex)
		remaining = remaining[closingBracket+1:]
	}
	return path, true
}

func (dec *propertiesDecoder) processComment(c string) string {
	if c == "" {
		return ""
	}
	return "# " + c
}

func (dec *propertiesDecoder) applyPropertyComments(context Context, path []interface{}, comments []string) error {
	assignmentOp := &Operation{OperationType: assignOpType, Preferences: assignPreferences{}}

	rhsCandidateNode := &CandidateNode{
		Tag:         "!!str",
		Value:       fmt.Sprintf("%v", path[len(path)-1]),
		HeadComment: dec.processComment(strings.Join(comments, "\n")),
		Kind:        ScalarNode,
	}

	rhsCandidateNode.Tag = rhsCandidateNode.guessTagFromCustomType()

	rhsOp := &Operation{OperationType: referenceOpType, CandidateNode: rhsCandidateNode}

	assignmentOpNode := &ExpressionNode{
		Operation: assignmentOp,
		LHS:       createTraversalTree(path, traversePreferences{}, true),
		RHS:       &ExpressionNode{Operation: rhsOp},
	}

	_, err := dec.d.GetMatchingNodes(context, assignmentOpNode)
	return err
}

func (dec *propertiesDecoder) applyProperty(context Context, properties *properties.Properties, key string) error {
	value, _ := properties.Get(key)
	path := parsePropKey(key, dec.prefs)

	propertyComments := properties.GetComments(key)
	if len(propertyComments) > 0 {
		err := dec.applyPropertyComments(context, path, propertyComments)
		if err != nil {
			return nil
		}
	}

	rhsNode := createStringScalarNode(value)
	rhsNode.Tag = rhsNode.guessTagFromCustomType()

	return dec.d.DeeplyAssign(context, path, rhsNode)
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
		Kind: MappingNode,
		Tag:  "!!map",
	}

	context := Context{}
	context = context.SingleChildContext(rootMap)

	for _, key := range properties.Keys() {
		if err := dec.applyProperty(context, properties, key); err != nil {
			return nil, err
		}

	}
	dec.finished = true

	return rootMap, nil

}
