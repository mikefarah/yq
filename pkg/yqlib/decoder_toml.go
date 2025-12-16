//go:build !yq_notoml

package yqlib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	toml "github.com/pelletier/go-toml/v2/unstable"
)

type tomlDecoder struct {
	parser           toml.Parser
	finished         bool
	d                DataTreeNavigator
	rootMap          *CandidateNode
	pendingComments  []string // Head comments collected from Comment nodes
	firstContentSeen bool     // Track if we've processed the first non-comment node
}

func NewTomlDecoder() Decoder {
	return &tomlDecoder{
		finished: false,
		d:        NewDataTreeNavigator(),
	}
}

func (dec *tomlDecoder) Init(reader io.Reader) error {
	dec.parser = toml.Parser{KeepComments: true}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return err
	}
	dec.parser.Reset(buf.Bytes())
	dec.rootMap = &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
	}
	dec.pendingComments = make([]string, 0)
	dec.firstContentSeen = false
	return nil
}

func (dec *tomlDecoder) getFullPath(tomlNode *toml.Node) []interface{} {
	path := make([]interface{}, 0)
	for {
		path = append(path, string(tomlNode.Data))
		tomlNode = tomlNode.Next()
		if tomlNode == nil {
			return path
		}
	}
}

func (dec *tomlDecoder) processKeyValueIntoMap(rootMap *CandidateNode, tomlNode *toml.Node) error {
	value := tomlNode.Value()
	path := dec.getFullPath(value.Next())

	valueNode, err := dec.decodeNode(value)
	if err != nil {
		return err
	}

	// Attach pending head comments
	if len(dec.pendingComments) > 0 {
		valueNode.HeadComment = strings.Join(dec.pendingComments, "\n")
		dec.pendingComments = make([]string, 0)
	}

	// Check for inline comment chained to the KeyValue node
	nextNode := tomlNode.Next()
	if nextNode != nil && nextNode.Kind == toml.Comment {
		valueNode.LineComment = string(nextNode.Data)
	}

	context := Context{}
	context = context.SingleChildContext(rootMap)

	return dec.d.DeeplyAssign(context, path, valueNode)
}

func (dec *tomlDecoder) decodeKeyValuesIntoMap(rootMap *CandidateNode, tomlNode *toml.Node) (bool, error) {
	log.Debug("decodeKeyValuesIntoMap -- processing first (current) entry")
	if err := dec.processKeyValueIntoMap(rootMap, tomlNode); err != nil {
		return false, err
	}

	for dec.parser.NextExpression() {
		nextItem := dec.parser.Expression()
		log.Debug("decodeKeyValuesIntoMap -- next exp, its a %v", nextItem.Kind)

		switch nextItem.Kind {
		case toml.KeyValue:
			if err := dec.processKeyValueIntoMap(rootMap, nextItem); err != nil {
				return false, err
			}
		case toml.Comment:
			// Standalone comment - add to pending for next element
			dec.pendingComments = append(dec.pendingComments, string(nextItem.Data))
		default:
			// run out of key values
			log.Debug("done in decodeKeyValuesIntoMap, gota a %v", nextItem.Kind)
			return true, nil
		}
	}
	log.Debug("no more things to read in")
	return false, nil
}

func (dec *tomlDecoder) createInlineTableMap(tomlNode *toml.Node) (*CandidateNode, error) {
	content := make([]*CandidateNode, 0)
	log.Debug("createInlineTableMap")

	iterator := tomlNode.Children()
	for iterator.Next() {
		child := iterator.Node()
		if child.Kind != toml.KeyValue {
			return nil, fmt.Errorf("only keyvalue pairs are supported in inlinetables, got %v instead", child.Kind)
		}

		keyValues := &CandidateNode{
			Kind: MappingNode,
			Tag:  "!!map",
		}

		if err := dec.processKeyValueIntoMap(keyValues, child); err != nil {
			return nil, err
		}

		content = append(content, keyValues.Content...)
	}

	return &CandidateNode{
		Kind:    MappingNode,
		Tag:     "!!map",
		Content: content,
	}, nil
}

func (dec *tomlDecoder) createArray(tomlNode *toml.Node) (*CandidateNode, error) {
	content := make([]*CandidateNode, 0)
	iterator := tomlNode.Children()
	for iterator.Next() {
		child := iterator.Node()
		yamlNode, err := dec.decodeNode(child)
		if err != nil {
			return nil, err
		}
		content = append(content, yamlNode)
	}

	return &CandidateNode{
		Kind:    SequenceNode,
		Tag:     "!!seq",
		Content: content,
	}, nil

}

func (dec *tomlDecoder) createStringScalar(tomlNode *toml.Node) (*CandidateNode, error) {
	content := string(tomlNode.Data)
	return createScalarNode(content, content), nil
}

func (dec *tomlDecoder) createBoolScalar(tomlNode *toml.Node) (*CandidateNode, error) {
	content := string(tomlNode.Data)
	return createScalarNode(content == "true", content), nil
}

func (dec *tomlDecoder) createIntegerScalar(tomlNode *toml.Node) (*CandidateNode, error) {
	content := string(tomlNode.Data)
	_, num, err := parseInt64(content)
	return createScalarNode(num, content), err
}

func (dec *tomlDecoder) createDateTimeScalar(tomlNode *toml.Node) (*CandidateNode, error) {
	content := string(tomlNode.Data)
	val, err := parseDateTime(time.RFC3339, content)
	return createScalarNode(val, content), err
}

func (dec *tomlDecoder) createFloatScalar(tomlNode *toml.Node) (*CandidateNode, error) {
	content := string(tomlNode.Data)
	num, err := strconv.ParseFloat(content, 64)
	return createScalarNode(num, content), err
}

func (dec *tomlDecoder) decodeNode(tomlNode *toml.Node) (*CandidateNode, error) {
	switch tomlNode.Kind {
	case toml.Key, toml.String:
		return dec.createStringScalar(tomlNode)
	case toml.Bool:
		return dec.createBoolScalar(tomlNode)
	case toml.Integer:
		return dec.createIntegerScalar(tomlNode)
	case toml.DateTime:
		return dec.createDateTimeScalar(tomlNode)
	case toml.Float:
		return dec.createFloatScalar(tomlNode)
	case toml.Array:
		return dec.createArray(tomlNode)
	case toml.InlineTable:
		return dec.createInlineTableMap(tomlNode)
	default:
		return nil, fmt.Errorf("unsupported type %v", tomlNode.Kind)
	}

}

func (dec *tomlDecoder) Decode() (*CandidateNode, error) {
	if dec.finished {
		return nil, io.EOF
	}
	//
	// toml library likes to panic
	var deferredError error
	defer func() { //catch or finally
		if r := recover(); r != nil {
			var ok bool
			deferredError, ok = r.(error)
			if !ok {
				deferredError = fmt.Errorf("pkg: %v", r)
			}
		}
	}()

	log.Debug("ok here we go")
	var runAgainstCurrentExp = false
	var err error
	for runAgainstCurrentExp || dec.parser.NextExpression() {

		if runAgainstCurrentExp {
			log.Debug("running against current exp")
		}

		currentNode := dec.parser.Expression()

		log.Debug("currentNode: %v ", currentNode.Kind)
		runAgainstCurrentExp, err = dec.processTopLevelNode(currentNode)
		if err != nil {
			return dec.rootMap, err
		}

	}

	err = dec.parser.Error()
	if err != nil {
		return nil, err
	}

	// must have finished
	dec.finished = true

	if len(dec.rootMap.Content) == 0 {
		return nil, io.EOF
	}

	return dec.rootMap, deferredError

}

func (dec *tomlDecoder) processTopLevelNode(currentNode *toml.Node) (bool, error) {
	var runAgainstCurrentExp bool
	var err error
	log.Debug("processTopLevelNode: Going to process %v state is current %v", currentNode.Kind, NodeToString(dec.rootMap))
	switch currentNode.Kind {
	case toml.Comment:
		// Collect comment to attach to next element
		commentText := string(currentNode.Data)
		// If we haven't seen any content yet, accumulate comments for root
		if !dec.firstContentSeen {
			if dec.rootMap.HeadComment == "" {
				dec.rootMap.HeadComment = commentText
			} else {
				dec.rootMap.HeadComment = dec.rootMap.HeadComment + "\n" + commentText
			}
		} else {
			// We've seen content, so these comments are for the next element
			dec.pendingComments = append(dec.pendingComments, commentText)
		}
		return false, nil
	case toml.Table:
		dec.firstContentSeen = true
		runAgainstCurrentExp, err = dec.processTable(currentNode)
	case toml.ArrayTable:
		dec.firstContentSeen = true
		runAgainstCurrentExp, err = dec.processArrayTable(currentNode)
	default:
		dec.firstContentSeen = true
		runAgainstCurrentExp, err = dec.decodeKeyValuesIntoMap(dec.rootMap, currentNode)
	}

	log.Debug("processTopLevelNode: DONE Processing state is now %v", NodeToString(dec.rootMap))
	return runAgainstCurrentExp, err
}

func (dec *tomlDecoder) processTable(currentNode *toml.Node) (bool, error) {
	log.Debug("Enter processTable")
	child := currentNode.Child()
	fullPath := dec.getFullPath(child)
	log.Debug("fullpath: %v", fullPath)

	c := Context{}
	c = c.SingleChildContext(dec.rootMap)

	fullPath, err := getPathToUse(fullPath, dec, c)
	if err != nil {
		return false, err
	}

	tableNodeValue := &CandidateNode{
		Kind:           MappingNode,
		Tag:            "!!map",
		Content:        make([]*CandidateNode, 0),
		EncodeSeparate: true,
	}

	// Attach pending head comments to the table
	if len(dec.pendingComments) > 0 {
		tableNodeValue.HeadComment = strings.Join(dec.pendingComments, "\n")
		dec.pendingComments = make([]string, 0)
	}

	var tableValue *toml.Node
	runAgainstCurrentExp := false
	hasValue := dec.parser.NextExpression()
	// check to see if there is any table data
	if hasValue {
		tableValue = dec.parser.Expression()
		// next expression is not table data, so we are done
		if tableValue.Kind != toml.KeyValue {
			log.Debug("got an empty table")
			runAgainstCurrentExp = true
		} else {
			runAgainstCurrentExp, err = dec.decodeKeyValuesIntoMap(tableNodeValue, tableValue)
			if err != nil && !errors.Is(err, io.EOF) {
				return false, err
			}
		}
	}

	err = dec.d.DeeplyAssign(c, fullPath, tableNodeValue)
	if err != nil {
		return false, err
	}
	return runAgainstCurrentExp, nil
}

func (dec *tomlDecoder) arrayAppend(context Context, path []interface{}, rhsNode *CandidateNode) error {
	log.Debug("arrayAppend to path: %v,%v", path, NodeToString(rhsNode))
	rhsCandidateNode := &CandidateNode{
		Kind:    SequenceNode,
		Tag:     "!!seq",
		Content: []*CandidateNode{rhsNode},
	}

	assignmentOp := &Operation{OperationType: addAssignOpType}

	rhsOp := &Operation{OperationType: valueOpType, CandidateNode: rhsCandidateNode}

	assignmentOpNode := &ExpressionNode{
		Operation: assignmentOp,
		LHS:       createTraversalTree(path, traversePreferences{}, false),
		RHS:       &ExpressionNode{Operation: rhsOp},
	}

	_, err := dec.d.GetMatchingNodes(context, assignmentOpNode)
	return err
}

func (dec *tomlDecoder) processArrayTable(currentNode *toml.Node) (bool, error) {
	log.Debug("Enter processArrayTable")
	child := currentNode.Child()
	fullPath := dec.getFullPath(child)
	log.Debug("Fullpath: %v", fullPath)

	c := Context{}
	c = c.SingleChildContext(dec.rootMap)

	fullPath, err := getPathToUse(fullPath, dec, c)
	if err != nil {
		return false, err
	}

	// need to use the array append exp to add another entry to
	// this array: fullpath += [ thing ]
	hasValue := dec.parser.NextExpression()

	tableNodeValue := &CandidateNode{
		Kind:           MappingNode,
		Tag:            "!!map",
		EncodeSeparate: true,
	}

	// Attach pending head comments to the array table
	if len(dec.pendingComments) > 0 {
		tableNodeValue.HeadComment = strings.Join(dec.pendingComments, "\n")
		dec.pendingComments = make([]string, 0)
	}

	runAgainstCurrentExp := false
	// if the next value is a ArrayTable or Table, then its not part of this declaration (not a key value pair)
	// so lets leave that expression for the next round of parsing
	if hasValue && (dec.parser.Expression().Kind == toml.ArrayTable || dec.parser.Expression().Kind == toml.Table) {
		runAgainstCurrentExp = true
	} else if hasValue {
		// otherwise, if there is a value, it must be some key value pairs of the
		// first object in the array!
		tableValue := dec.parser.Expression()
		runAgainstCurrentExp, err = dec.decodeKeyValuesIntoMap(tableNodeValue, tableValue)
		if err != nil && !errors.Is(err, io.EOF) {
			return false, err
		}
	}

	// += function
	err = dec.arrayAppend(c, fullPath, tableNodeValue)

	return runAgainstCurrentExp, err
}

// if fullPath points to an array of maps rather than a map
// then it should set this element into the _last_ element of that array.
// Because TOML. So we'll inject the last index into the path.

func getPathToUse(fullPath []interface{}, dec *tomlDecoder, c Context) ([]interface{}, error) {
	pathToCheck := fullPath
	if len(fullPath) >= 1 {
		pathToCheck = fullPath[:len(fullPath)-1]
	}
	readOp := createTraversalTree(pathToCheck, traversePreferences{DontAutoCreate: true}, false)

	resultContext, err := dec.d.GetMatchingNodes(c, readOp)
	if err != nil {
		return nil, err
	}
	if resultContext.MatchingNodes.Len() >= 1 {
		match := resultContext.MatchingNodes.Front().Value.(*CandidateNode)
		// path refers to an array, we need to add this to the last element in the array
		if match.Kind == SequenceNode {
			fullPath = append(pathToCheck, len(match.Content)-1, fullPath[len(fullPath)-1])
			log.Debugf("Adding to end of %v array, using path: %v", pathToCheck, fullPath)
		}
	}
	return fullPath, err
}
