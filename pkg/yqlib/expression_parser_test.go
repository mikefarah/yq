package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func getExpressionParser() ExpressionParserInterface {
	InitExpressionParser()
	return ExpressionParser
}

func TestParserCreateMapColonOnItsOwn(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(":")
	test.AssertResultComplex(t, "':' expects 2 args but there is 0", err.Error())
}

func TestParserNoMatchingCloseBracket(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(".cat | with(.;.bob")
	test.AssertResultComplex(t, "bad expression - probably missing close bracket on WITH", err.Error())
}

func TestParserNoMatchingCloseCollect(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("[1,2")
	test.AssertResultComplex(t, "bad expression, could not find matching `]`", err.Error())
}
func TestParserNoMatchingCloseObjectInCollect(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`[{"b": "c"]`)
	test.AssertResultComplex(t, "bad expression, could not find matching `}`", err.Error())
}

func TestParserNoMatchingCloseInCollect(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`[(.a]`)
	test.AssertResultComplex(t, "bad expression, could not find matching `)`", err.Error())
}

func TestParserNoMatchingCloseCollectObject(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`{"a": "b"`)
	test.AssertResultComplex(t, "bad expression, could not find matching `}`", err.Error())
}

func TestParserNoMatchingCloseCollectInCollectObject(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`{"b": [1}`)
	test.AssertResultComplex(t, "bad expression, could not find matching `]`", err.Error())
}

func TestParserNoMatchingCloseBracketInCollectObject(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`{"b": (1}`)
	test.AssertResultComplex(t, "bad expression, could not find matching `)`", err.Error())
}

func TestParserNoArgsForTwoArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("=")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 0", err.Error())
}

func TestParserOneLhsArgsForTwoArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(".a =")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestParserOneRhsArgsForTwoArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("= .a")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestParserTwoArgsForTwoArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(".a = .b")
	test.AssertResultComplex(t, nil, err)
}

func TestParserNoArgsForOneArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("explode")
	test.AssertResultComplex(t, "'explode' expects 1 arg but received none", err.Error())
}

func TestParserOneArgForOneArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("explode(.)")
	test.AssertResultComplex(t, nil, err)
}

func TestParserExtraArgs(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("sortKeys(.) explode(.)")
	test.AssertResultComplex(t, "bad expression, please check expression syntax", err.Error())
}
