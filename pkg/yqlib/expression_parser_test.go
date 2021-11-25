package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func TestParserNoMatchingCloseCollect(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("[1,2")
	test.AssertResultComplex(t, "Bad expression, could not find matching `]`", err.Error())
}

func TestParserNoMatchingCloseObjectInCollect(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression(`[{"b": "c"]`)
	test.AssertResultComplex(t, "Bad expression, could not find matching `}`", err.Error())
}

func TestParserNoMatchingCloseInCollect(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression(`[(.a]`)
	test.AssertResultComplex(t, "Bad expression, could not find matching `)`", err.Error())
}

func TestParserNoMatchingCloseCollectObject(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression(`{"a": "b"`)
	test.AssertResultComplex(t, "Bad expression, could not find matching `}`", err.Error())
}

func TestParserNoMatchingCloseCollectInCollectObject(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression(`{"b": [1}`)
	test.AssertResultComplex(t, "Bad expression, could not find matching `]`", err.Error())
}

func TestParserNoMatchingCloseBracketInCollectObject(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression(`{"b": (1}`)
	test.AssertResultComplex(t, "Bad expression, could not find matching `)`", err.Error())
}

func TestParserNoArgsForTwoArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("=")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 0", err.Error())
}

func TestParserOneLhsArgsForTwoArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression(".a =")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestParserOneRhsArgsForTwoArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("= .a")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestParserTwoArgsForTwoArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression(".a = .b")
	test.AssertResultComplex(t, nil, err)
}

func TestParserNoArgsForOneArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("explode")
	test.AssertResultComplex(t, "'explode' expects 1 arg but received none", err.Error())
}

func TestParserOneArgForOneArgOp(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("explode(.)")
	test.AssertResultComplex(t, nil, err)
}

func TestParserExtraArgs(t *testing.T) {
	_, err := NewExpressionParser().ParseExpression("sortKeys(.) explode(.)")
	test.AssertResultComplex(t, "Bad expression, please check expression syntax", err.Error())
}
