package yqlib

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	yaml "gopkg.in/yaml.v3"
)

type multiplyPreferences struct {
	AppendArrays    bool
	DeepMergeArrays bool
	TraversePrefs   traversePreferences
	AssignPrefs     assignPreferences
}

func createMultiplyOp(prefs interface{}) func(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode {
	return func(lhs *ExpressionNode, rhs *ExpressionNode) *ExpressionNode {
		return &ExpressionNode{Operation: &Operation{OperationType: multiplyOpType, Preferences: prefs},
			LHS: lhs,
			RHS: rhs}
	}
}

func multiplyAssignOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var multiplyPrefs = expressionNode.Operation.Preferences

	return compoundAssignFunction(d, context, expressionNode, createMultiplyOp(multiplyPrefs))
}

func multiplyOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- MultiplyOperator")
	return crossFunction(d, context, expressionNode, multiply(expressionNode.Operation.Preferences.(multiplyPreferences)), false)
}

func getComments(lhs *CandidateNode, rhs *CandidateNode) (leadingContent string, headComment string, footComment string) {
	leadingContent = rhs.LeadingContent
	headComment = rhs.Node.HeadComment
	footComment = rhs.Node.FootComment
	if lhs.Node.HeadComment != "" || lhs.LeadingContent != "" {
		headComment = lhs.Node.HeadComment
		leadingContent = lhs.LeadingContent
	}

	if lhs.Node.FootComment != "" {
		footComment = lhs.Node.FootComment
	}
	return leadingContent, headComment, footComment
}

func multiply(preferences multiplyPreferences) func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		// need to do this before unWrapping the potential document node
		leadingContent, headComment, footComment := getComments(lhs, rhs)
		lhs.Node = unwrapDoc(lhs.Node)
		rhs.Node = unwrapDoc(rhs.Node)
		log.Debugf("Multiplying LHS: %v", lhs.Node.Tag)
		log.Debugf("-          RHS: %v", rhs.Node.Tag)

		if rhs.Node.Tag == "!!null" {
			return lhs.Copy()
		}

		if (lhs.Node.Kind == yaml.MappingNode && rhs.Node.Kind == yaml.MappingNode) ||
			(lhs.Node.Tag == "!!null" && rhs.Node.Kind == yaml.MappingNode) ||
			(lhs.Node.Kind == yaml.SequenceNode && rhs.Node.Kind == yaml.SequenceNode) ||
			(lhs.Node.Tag == "!!null" && rhs.Node.Kind == yaml.SequenceNode) {
			var newBlank = CandidateNode{}
			err := copier.CopyWithOption(&newBlank, lhs, copier.Option{IgnoreEmpty: true, DeepCopy: true})
			if err != nil {
				return nil, err
			}
			newBlank.LeadingContent = leadingContent
			newBlank.Node.HeadComment = headComment
			newBlank.Node.FootComment = footComment

			return mergeObjects(d, context.WritableClone(), &newBlank, rhs, preferences)
		}
		return multiplyScalars(lhs, rhs)
	}
}

func multiplyScalars(lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhsTag := lhs.Node.Tag
	rhsTag := guessTagFromCustomType(rhs.Node)
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = guessTagFromCustomType(lhs.Node)
		lhsIsCustom = true
	}

	if lhsTag == "!!int" && rhsTag == "!!int" {
		return multiplyIntegers(lhs, rhs)
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		return multiplyFloats(lhs, rhs, lhsIsCustom)
	}
	return nil, fmt.Errorf("Cannot multiply %v with %v", lhs.Node.Tag, rhs.Node.Tag)
}

func multiplyFloats(lhs *CandidateNode, rhs *CandidateNode, lhsIsCustom bool) (*CandidateNode, error) {
	target := lhs.CreateReplacement(&yaml.Node{})
	target.Node.Kind = yaml.ScalarNode
	target.Node.Style = lhs.Node.Style
	if lhsIsCustom {
		target.Node.Tag = lhs.Node.Tag
	} else {
		target.Node.Tag = "!!float"
	}

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
	target := lhs.CreateReplacement(&yaml.Node{})
	target.Node.Kind = yaml.ScalarNode
	target.Node.Style = lhs.Node.Style
	target.Node.Tag = lhs.Node.Tag

	format, lhsNum, err := parseInt64(lhs.Node.Value)
	if err != nil {
		return nil, err
	}
	_, rhsNum, err := parseInt64(rhs.Node.Value)
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
	err := recursiveDecent(results, context.SingleChildContext(rhs), prefs)
	if err != nil {
		return nil, err
	}

	var pathIndexToStartFrom int
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

	assignmentOp := &Operation{OperationType: assignAttributesOpType, Preferences: preferences.AssignPrefs}
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
	rhsOp := &Operation{OperationType: referenceOpType, CandidateNode: rhs}

	assignmentOpNode := &ExpressionNode{
		Operation: assignmentOp,
		LHS:       createTraversalTree(lhsPath, preferences.TraversePrefs, rhs.IsMapKey),
		RHS:       &ExpressionNode{Operation: rhsOp},
	}

	_, err := d.GetMatchingNodes(context.SingleChildContext(lhs), assignmentOpNode)

	return err
}
