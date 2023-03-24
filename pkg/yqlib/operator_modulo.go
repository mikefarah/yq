package yqlib

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func moduloOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("Modulo operator")

	return crossFunction(d, context.ReadOnlyClone(), expressionNode, modulo, false)
}

func modulo(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	lhs.Node = unwrapDoc(lhs.Node)
	rhs.Node = unwrapDoc(rhs.Node)

	lhsNode := lhs.Node

	if lhsNode.Tag == "!!null" {
		return nil, fmt.Errorf("%v (%v) cannot modulo by %v (%v)", lhsNode.Tag, lhs.GetNicePath(), rhs.Node.Tag, rhs.GetNicePath())
	}

	target := &yaml.Node{}

	if lhsNode.Kind == yaml.ScalarNode && rhs.Node.Kind == yaml.ScalarNode {
		if err := moduloScalars(target, lhsNode, rhs.Node); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%v (%v) cannot modulo by %v (%v)", lhsNode.Tag, lhs.GetNicePath(), rhs.Node.Tag, rhs.GetNicePath())
	}

	return lhs.CreateReplacement(target), nil
}

func moduloScalars(target *yaml.Node, lhs *yaml.Node, rhs *yaml.Node) error {
	lhsTag := lhs.Tag
	rhsTag := guessTagFromCustomType(rhs)
	lhsIsCustom := false
	if !strings.HasPrefix(lhsTag, "!!") {
		// custom tag - we have to have a guess
		lhsTag = guessTagFromCustomType(lhs)
		lhsIsCustom = true
	}

	if lhsTag == "!!int" && rhsTag == "!!int" {
		target.Kind = yaml.ScalarNode
		target.Style = lhs.Style

		format, lhsNum, err := parseInt64(lhs.Value)
		if err != nil {
			return err
		}
		_, rhsNum, err := parseInt64(rhs.Value)
		if err != nil {
			return err
		}
		if rhsNum == 0 {
			return fmt.Errorf("cannot modulo by 0")
		}
		remainder := lhsNum % rhsNum

		target.Tag = lhs.Tag
		target.Value = fmt.Sprintf(format, remainder)
	} else if (lhsTag == "!!int" || lhsTag == "!!float") && (rhsTag == "!!int" || rhsTag == "!!float") {
		target.Kind = yaml.ScalarNode
		target.Style = lhs.Style

		lhsNum, err := strconv.ParseFloat(lhs.Value, 64)
		if err != nil {
			return err
		}
		rhsNum, err := strconv.ParseFloat(rhs.Value, 64)
		if err != nil {
			return err
		}
		remainder := math.Mod(lhsNum, rhsNum)
		if lhsIsCustom {
			target.Tag = lhs.Tag
		} else {
			target.Tag = "!!float"
		}
		target.Value = fmt.Sprintf("%v", remainder)
	} else {
		return fmt.Errorf("%v cannot modulo by %v", lhsTag, rhsTag)
	}
	return nil
}
