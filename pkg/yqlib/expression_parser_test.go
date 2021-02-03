package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func TestPathTreeNoArgsForTwoArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("=")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 0", err.Error())
}

func TestPathTreeOneLhsArgsForTwoArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression(".a =")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestPathTreeOneRhsArgsForTwoArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("= .a")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestPathTreeTwoArgsForTwoArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression(".a = .b")
	test.AssertResultComplex(t, nil, err)
}

func TestPathTreeNoArgsForOneArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("explode")
	test.AssertResultComplex(t, "'explode' expects 1 arg but received none", err.Error())
}

func TestPathTreeOneArgForOneArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("explode(.)")
	test.AssertResultComplex(t, nil, err)
}

func TestPathTreeExtraArgs(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("sortKeys(.) explode(.)")
	test.AssertResultComplex(t, "Bad expression, please check expression syntax", err.Error())
}
