package yqlib

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
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
	log.Debugf("MultiplyOperator")
	return crossFunction(d, context.ReadOnlyClone(), expressionNode, multiply(expressionNode.Operation.Preferences.(multiplyPreferences)), false)
}

func getComments(lhs *CandidateNode, rhs *CandidateNode) (leadingContent string, headComment string, footComment string) {
	leadingContent = rhs.LeadingContent
	headComment = rhs.HeadComment
	footComment = rhs.FootComment
	if lhs.HeadComment != "" || lhs.LeadingContent != "" {
		headComment = lhs.HeadComment
		leadingContent = lhs.LeadingContent
	}

	if lhs.FootComment != "" {
		footComment = lhs.FootComment
	}

	return leadingContent, headComment, footComment
}

func multiply(preferences multiplyPreferences) func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		// need to do this before unWrapping the potential document node
		leadingContent, headComment, footComment := getComments(lhs, rhs)
		log.Debugf("Multiplying LHS: %v", NodeToString(lhs))
		log.Debugf("-           RHS: %v", NodeToString(rhs))

		if rhs.Tag == "!!null" {
			return lhs.Copy(), nil
		}

		if (lhs.Kind == MappingNode && rhs.Kind == MappingNode) ||
			(lhs.Tag == "!!null" && rhs.Kind == MappingNode) ||
			(lhs.Kind == SequenceNode && rhs.Kind == SequenceNode) ||
			(lhs.Tag == "!!null" && rhs.Kind == SequenceNode) {

			var newBlank = lhs.Copy()

			newBlank.LeadingContent = leadingContent
			newBlank.HeadComment = headComment
			newBlank.FootComment = footComment

			return mergeObjects(d, context.WritableClone(), newBlank, rhs, preferences)
		}
		return multiplyScalars(lhs, rhs)
	}
}

func multiplyScalars(lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhsTag := lhs.Tag
	rhsTag := rhs.guessTagFromCustomType()
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = lhs.guessTagFromCustomType()
		lhsIsCustom = true
	}

	if lhsTag == "!!int" && rhsTag == "!!int" {
		return multiplyIntegers(lhs, rhs)
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		return multiplyFloats(lhs, rhs, lhsIsCustom)
	} else if (lhsTag == "!!str" && rhsTag == "!!int") || (lhsTag == "!!int" && rhsTag == "!!str") {
		return repeatString(lhs, rhs)
	}
	return nil, fmt.Errorf("cannot multiply %v with %v", lhs.Tag, rhs.Tag)
}

func multiplyFloats(lhs *CandidateNode, rhs *CandidateNode, lhsIsCustom bool) (*CandidateNode, error) {
	target := lhs.CopyWithoutContent()
	target.Kind = ScalarNode
	target.Style = lhs.Style
	if lhsIsCustom {
		target.Tag = lhs.Tag
	} else {
		target.Tag = "!!float"
	}

	lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
	if err != nil {
		return nil, err
	}
	rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
	if err != nil {
		return nil, err
	}
	target.Value = fmt.Sprintf("%v", lhsNum*rhsNum)
	return target, nil
}

func multiplyIntegers(lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	target := lhs.CopyWithoutContent()
	target.Kind = ScalarNode
	target.Style = lhs.Style
	target.Tag = lhs.Tag

	format, lhsNum, err := parseInt64(lhs.Value)
	if err != nil {
		return nil, err
	}
	_, rhsNum, err := parseInt64(rhs.Value)
	if err != nil {
		return nil, err
	}
	target.Value = fmt.Sprintf(format, lhsNum*rhsNum)
	return target, nil
}

func repeatString(lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	var stringNode *CandidateNode
	var intNode *CandidateNode
	if lhs.Tag == "!!str" {
		stringNode = lhs
		intNode = rhs
	} else {
		stringNode = rhs
		intNode = lhs
	}
	target := lhs.CopyWithoutContent()
	target.UpdateAttributesFrom(stringNode, assignPreferences{})

	count, err := parseInt(intNode.Value)
	if err != nil {
		return nil, err
	} else if count < 0 {
		return nil, fmt.Errorf("cannot repeat string by a negative number (%v)", count)
	} else if count > 10000000 {
		return nil, fmt.Errorf("cannot repeat string by more than 100 million (%v)", count)
	}
	target.Value = strings.Repeat(stringNode.Value, count)

	return target, nil
}

func mergeObjects(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode, preferences multiplyPreferences) (*CandidateNode, error) {
	var results = list.New()

	// only need to recurse the array if we are doing a deep merge
	prefs := recursiveDescentPreferences{RecurseArray: preferences.DeepMergeArrays,
		TraversePreferences: traversePreferences{DontFollowAlias: true, IncludeMapKeys: true, ExactKeyMatch: true}}
	log.Debugf("merge - preferences.DeepMergeArrays %v", preferences.DeepMergeArrays)
	log.Debugf("merge - preferences.AppendArrays %v", preferences.AppendArrays)
	err := recursiveDecent(results, context.SingleChildContext(rhs), prefs)
	if err != nil {
		return nil, err
	}

	var pathIndexToStartFrom int
	if results.Front() != nil {
		pathIndexToStartFrom = len(results.Front().Value.(*CandidateNode).GetPath())
		log.Debugf("pathIndexToStartFrom: %v", pathIndexToStartFrom)
	}

	for el := results.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		log.Debugf("going to applied assignment to LHS: %v with RHS: %v", NodeToString(lhs), NodeToString(candidate))

		if candidate.Tag == "!!merge" {
			continue
		}

		err := applyAssignment(d, context, pathIndexToStartFrom, lhs, candidate, preferences)
		if err != nil {
			return nil, err
		}

		log.Debugf("applied assignment to LHS: %v", NodeToString(lhs))
	}
	return lhs, nil
}

func applyAssignment(d *dataTreeNavigator, context Context, pathIndexToStartFrom int, lhs *CandidateNode, rhs *CandidateNode, preferences multiplyPreferences) error {
	shouldAppendArrays := preferences.AppendArrays

	lhsPath := rhs.GetPath()[pathIndexToStartFrom:]
	log.Debugf("merge - lhsPath %v", lhsPath)

	assignmentOp := &Operation{OperationType: assignAttributesOpType, Preferences: preferences.AssignPrefs}
	if shouldAppendArrays && rhs.Kind == SequenceNode {
		assignmentOp.OperationType = addAssignOpType
		log.Debugf("merge - assignmentOp.OperationType = addAssignOpType")
	} else if !preferences.DeepMergeArrays && rhs.Kind == SequenceNode ||
		(rhs.Kind == ScalarNode || rhs.Kind == AliasNode) {
		assignmentOp.OperationType = assignOpType
		assignmentOp.UpdateAssign = false
		log.Debugf("merge - rhs.Kind == SequenceNode: %v", rhs.Kind == SequenceNode)
		log.Debugf("merge - rhs.Kind == ScalarNode: %v", rhs.Kind == ScalarNode)
		log.Debugf("merge - rhs.Kind == AliasNode: %v", rhs.Kind == AliasNode)
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
