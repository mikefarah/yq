package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v3/test"
)

// var tokeniser = NewPathTokeniser()
var postFixer = NewPathPostFixer()

func testExpression(expression string) (string, error) {
	tokens, err := tokeniser.Tokenise(expression)
	if err != nil {
		return "", err
	}
	results, errorP := postFixer.ConvertToPostfix(tokens)
	if errorP != nil {
		return "", errorP
	}
	formatted := ""
	for _, path := range results {
		formatted = formatted + path.toString() + "--------\n"
	}
	return formatted, nil
}

func TestPostFixSimple(t *testing.T) {
	var infix = "a"
	var expectedOutput = `Type: PathKey - a
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimplePath(t *testing.T) {
	var infix = "apples.bananas*.cat"
	var expectedOutput = `Type: PathKey - apples.bananas*.cat
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixOr(t *testing.T) {
	var infix = "a OR b"
	var expectedOutput = `Type: PathKey - a
--------
Type: PathKey - b
--------
Type: Operation - OR
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixOrWithEquals(t *testing.T) {
	var infix = "a==thing OR b==thongs"
	var expectedOutput = `Type: PathKey - a
--------
Type: PathKey - thing
--------
Type: Operation - EQUALS
--------
Type: PathKey - b
--------
Type: PathKey - thongs
--------
Type: Operation - EQUALS
--------
Type: Operation - OR
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixOrWithEqualsPath(t *testing.T) {
	var infix = "apples.monkeys==thing OR bogs.bobos==thongs"
	var expectedOutput = `Type: PathKey - apples.monkeys
--------
Type: PathKey - thing
--------
Type: Operation - EQUALS
--------
Type: PathKey - bogs.bobos
--------
Type: PathKey - thongs
--------
Type: Operation - EQUALS
--------
Type: Operation - OR
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}
