package yqlib

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func getStringParamter(parameterName string, d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (string, error) {
	result, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode)

	if err != nil {
		return "", err
	} else if result.MatchingNodes.Len() == 0 {
		return "", fmt.Errorf("could not find %v for format_time", parameterName)
	}

	return result.MatchingNodes.Front().Value.(*CandidateNode).Node.Value, nil
}

func withDateTimeFormat(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	if expressionNode.RHS.Operation.OperationType == blockOpType || expressionNode.RHS.Operation.OperationType == unionOpType {
		layout, err := getStringParamter("layout", d, context, expressionNode.RHS.LHS)
		if err != nil {
			return Context{}, fmt.Errorf("could not get date time format: %w", err)
		}
		context.SetDateTimeLayout(layout)
		return d.GetMatchingNodes(context, expressionNode.RHS.RHS)

	}
	return Context{}, errors.New(`must provide a date time format string and an expression, e.g. with_dtf("Monday, 02-Jan-06 at 3:04PM MST"; <exp>)`)

}

// for unit tests
var Now = time.Now

func nowOp(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	node := &yaml.Node{
		Tag:   "!!timestamp",
		Kind:  yaml.ScalarNode,
		Value: Now().Format(time.RFC3339),
	}

	return context.SingleChildContext(&CandidateNode{Node: node}), nil

}

func formatDateTime(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	format, err := getStringParamter("format", d, context, expressionNode.RHS)
	layout := context.GetDateTimeLayout()
	decoder := NewYamlDecoder()

	if err != nil {
		return Context{}, err
	}
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		parsedTime, err := time.Parse(layout, candidate.Node.Value)
		if err != nil {
			return Context{}, fmt.Errorf("could not parse datetime of [%v]: %w", candidate.GetNicePath(), err)
		}
		formattedTimeStr := parsedTime.Format(format)
		decoder.Init(strings.NewReader(formattedTimeStr))
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)
		var node *yaml.Node
		if errorReading != nil {
			log.Debugf("could not parse %v - lets just leave it as a string", formattedTimeStr)
			node = &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: formattedTimeStr,
			}
		} else {
			node = unwrapDoc(&dataBucket)
		}

		results.PushBack(candidate.CreateReplacement(node))
	}

	return context.ChildContext(results), nil
}

func tzOp(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	timezoneStr, err := getStringParamter("timezone", d, context, expressionNode.RHS)
	layout := context.GetDateTimeLayout()

	if err != nil {
		return Context{}, err
	}
	var results = list.New()

	timezone, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return Context{}, fmt.Errorf("could not load tz [%v]: %w", timezoneStr, err)
	}

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		parsedTime, err := time.Parse(layout, candidate.Node.Value)
		if err != nil {
			return Context{}, fmt.Errorf("could not parse datetime of [%v] using layout [%v]: %w", candidate.GetNicePath(), layout, err)
		}
		tzTime := parsedTime.In(timezone)

		node := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   candidate.Node.Tag,
			Value: tzTime.Format(layout),
		}

		results.PushBack(candidate.CreateReplacement(node))
	}

	return context.ChildContext(results), nil
}
