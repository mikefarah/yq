package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func TestPathTreeNoArgsForTwoArgOp(t *testing.T) {
	_, err := NewPathTreeCreator().ParsePath("=")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 0", err.Error())
}

func TestPathTreeOneLhsArgsForTwoArgOp(t *testing.T) {
	_, err := NewPathTreeCreator().ParsePath(".a =")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestPathTreeOneRhsArgsForTwoArgOp(t *testing.T) {
	_, err := NewPathTreeCreator().ParsePath("= .a")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestPathTreeTwoArgsForTwoArgOp(t *testing.T) {
	_, err := NewPathTreeCreator().ParsePath(".a = .b")
	test.AssertResultComplex(t, nil, err)
}

func TestPathTreeNoArgsForOneArgOp(t *testing.T) {
	_, err := NewPathTreeCreator().ParsePath("explode")
	test.AssertResultComplex(t, "'explode' expects 1 arg but received none", err.Error())
}

func TestPathTreeOneArgForOneArgOp(t *testing.T) {
	_, err := NewPathTreeCreator().ParsePath("explode(.)")
	test.AssertResultComplex(t, nil, err)
}

func TestPathTreeExtraArgs(t *testing.T) {
	_, err := NewPathTreeCreator().ParsePath("sortKeys(.) explode(.)")
	test.AssertResultComplex(t, "expected end of expression but found 'explode', please check expression syntax", err.Error())
}
