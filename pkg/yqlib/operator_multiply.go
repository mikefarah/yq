package yqlib

import (
	"fmt"
	"strconv"

	"container/list"

	"github.com/jinzhu/copier"
	yaml "gopkg.in/yaml.v3"
)

type multiplyPreferences struct {
	AppendArrays    bool
	DeepMergeArrays bool
	TraversePrefs   traversePreferences
}

func multiplyOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- MultiplyOperator")
	return crossFunction(d, context, expressionNode, multiply(expressionNode.Operation.Preferences.(multiplyPreferences)), false)
}

func getNewBlankNode(lhs *yaml.Node, rhs *yaml.Node) *yaml.Node {

	blankNode := &yaml.Node{}

	if lhs.HeadComment != "" {
		blankNode.HeadComment = lhs.HeadComment
	} else if rhs.HeadComment != "" {
		blankNode.HeadComment = rhs.HeadComment
	}

	if lhs.FootComment != "" {
		blankNode.FootComment = lhs.FootComment
	} else if rhs.FootComment != "" {
		blankNode.FootComment = rhs.FootComment
	}

	return blankNode
}

func multiply(preferences multiplyPreferences) func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		// need to do this before unWrapping the potential document node
		newBlankNode := getNewBlankNode(lhs.Node, rhs.Node)
		lhs.Node = unwrapDoc(lhs.Node)
		rhs.Node = unwrapDoc(rhs.Node)
		log.Debugf("Multipling LHS: %v", lhs.Node.Tag)
		log.Debugf("-          RHS: %v", rhs.Node.Tag)

		if lhs.Node.Kind == yaml.MappingNode && rhs.Node.Kind == yaml.MappingNode ||
			(lhs.Node.Kind == yaml.SequenceNode && rhs.Node.Kind == yaml.SequenceNode) {
			var newBlank = CandidateNode{}
			err := copier.CopyWithOption(&newBlank, lhs, copier.Option{IgnoreEmpty: true, DeepCopy: true})
			if err != nil {
				return nil, err
			}
			newBlank.Node.HeadComment = newBlankNode.HeadComment
			newBlank.Node.FootComment = newBlankNode.FootComment

			// var newBlank = lhs.CreateChild(nil, newBlankNode)
			// log.Debugf("merge - merge lhs into blank")
			// var newThing, err = mergeObjects(d, context.WritableClone(), newBlank, lhs, multiplyPreferences{})
			// if err != nil {
			// 	return nil, err
			// }
			return mergeObjects(d, context.WritableClone(), &newBlank, rhs, preferences)
		} else if lhs.Node.Tag == "!!int" && rhs.Node.Tag == "!!int" {
			return multiplyIntegers(lhs, rhs)
		} else if (lhs.Node.Tag == "!!int" || lhs.Node.Tag == "!!float") && (rhs.Node.Tag == "!!int" || rhs.Node.Tag == "!!float") {
			return multiplyFloats(lhs, rhs)
		}
		return nil, fmt.Errorf("Cannot multiply %v with %v", lhs.Node.Tag, rhs.Node.Tag)
	}
}

func multiplyFloats(lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	target := lhs.CreateChild(nil, &yaml.Node{})
	target.Node.Kind = yaml.ScalarNode
	target.Node.Style = lhs.Node.Style
	target.Node.Tag = "!!float"

	lhsNum, err := strconv.ParseFloat(lhs.Node.Value, 64)
	if err != nil {
		return nil, err
	}
	rhsNum, err := strconv.ParseFloat(rhs.Node.Value, 64)
	if err != nil {
		return nil, err
	}
	target.Node.Value = fmt.Sprintf("%v", lhsNum*rhsNum)
	return target, nil
}

func multiplyIntegers(lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	target := lhs.CreateChild(nil, &yaml.Node{})
	target.Node.Kind = yaml.ScalarNode
	target.Node.Style = lhs.Node.Style
	target.Node.Tag = "!!int"

	format, lhsNum, err := parseInt(lhs.Node.Value)
	if err != nil {
		return nil, err
	}
	_, rhsNum, err := parseInt(rhs.Node.Value)
	if err != nil {
		return nil, err
	}
	target.Node.Value = fmt.Sprintf(format, lhsNum*rhsNum)
	return target, nil
}

func mergeObjects(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode, preferences multiplyPreferences) (*CandidateNode, error) {
	var results = list.New()

	// only need to recurse the array if we are doing a deep merge
	prefs := recursiveDescentPreferences{RecurseArray: preferences.DeepMergeArrays,
		TraversePreferences: traversePreferences{DontFollowAlias: true, IncludeMapKeys: true}}
	log.Debugf("merge - preferences.DeepMergeArrays %v", preferences.DeepMergeArrays)
	log.Debugf("merge - preferences.AppendArrays %v", preferences.AppendArrays)
	err := recursiveDecent(d, results, context.SingleChildContext(rhs), prefs)
	if err != nil {
		return nil, err
	}

	var pathIndexToStartFrom int = 0
	if results.Front() != nil {
		pathIndexToStartFrom = len(results.Front().Value.(*CandidateNode).Path)
	}

	for el := results.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		if candidate.Node.Tag == "!!merge" {
			continue
		}

		err := applyAssignment(d, context, pathIndexToStartFrom, lhs, candidate, preferences)
		if err != nil {
			return nil, err
		}
	}
	return lhs, nil
}

func applyAssignment(d *dataTreeNavigator, context Context, pathIndexToStartFrom int, lhs *CandidateNode, rhs *CandidateNode, preferences multiplyPreferences) error {
	shouldAppendArrays := preferences.AppendArrays
	log.Debugf("merge - applyAssignment lhs %v, rhs: %v", lhs.GetKey(), rhs.GetKey())

	lhsPath := rhs.Path[pathIndexToStartFrom:]
	log.Debugf("merge - lhsPath %v", lhsPath)

	assignmentOp := &Operation{OperationType: assignAttributesOpType}
	if shouldAppendArrays && rhs.Node.Kind == yaml.SequenceNode {
		assignmentOp.OperationType = addAssignOpType
		log.Debugf("merge - assignmentOp.OperationType = addAssignOpType")
	} else if !preferences.DeepMergeArrays && rhs.Node.Kind == yaml.SequenceNode ||
		(rhs.Node.Kind == yaml.ScalarNode || rhs.Node.Kind == yaml.AliasNode) {
		assignmentOp.OperationType = assignOpType
		assignmentOp.UpdateAssign = false
		log.Debugf("merge - rhs.Node.Kind == yaml.SequenceNode: %v", rhs.Node.Kind == yaml.SequenceNode)
		log.Debugf("merge - rhs.Node.Kind == yaml.ScalarNode: %v", rhs.Node.Kind == yaml.ScalarNode)
		log.Debugf("merge - rhs.Node.Kind == yaml.AliasNode: %v", rhs.Node.Kind == yaml.AliasNode)
		log.Debugf("merge - assignmentOp.OperationType = assignOpType, no updateassign")
	} else {
		log.Debugf("merge - assignmentOp := &Operation{OperationType: assignAttributesOpType}")
	}
	rhsOp := &Operation{OperationType: valueOpType, CandidateNode: rhs}

	assignmentOpNode := &ExpressionNode{Operation: assignmentOp, Lhs: createTraversalTree(lhsPath, preferences.TraversePrefs, rhs.IsMapKey), Rhs: &ExpressionNode{Operation: rhsOp}}

	_, err := d.GetMatchingNodes(context.SingleChildContext(lhs), assignmentOpNode)

	return err
}
