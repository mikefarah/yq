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

func TestPostFixSimpleExample(t *testing.T) {
	var infix = "a"
	var expectedOutput = `PathKey - 'a'
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimplePathExample(t *testing.T) {
	var infix = "apples.bananas*.cat"
	var expectedOutput = `PathKey - 'apples'
--------
PathKey - 'bananas*'
--------
Operation - TRAVERSE
--------
PathKey - 'cat'
--------
Operation - TRAVERSE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimplePathNumbersExample(t *testing.T) {
	var infix = "apples[0].cat"
	var expectedOutput = `PathKey - 'apples'
--------
PathKey - '0'
--------
Operation - TRAVERSE
--------
PathKey - 'cat'
--------
Operation - TRAVERSE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimplePathAppendArrayExample(t *testing.T) {
	var infix = "apples[+].cat"
	var expectedOutput = `PathKey - 'apples'
--------
PathKey - '[+]'
--------
Operation - TRAVERSE
--------
PathKey - 'cat'
--------
Operation - TRAVERSE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimplePathSplatArrayExample(t *testing.T) {
	var infix = "apples.[*]cat"
	var expectedOutput = `PathKey - 'apples'
--------
PathKey - '[*]'
--------
Operation - TRAVERSE
--------
PathKey - 'cat'
--------
Operation - TRAVERSE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixDeepMatchExample(t *testing.T) {
	var infix = "apples.**.cat"
	var expectedOutput = `PathKey - 'apples'
--------
PathKey - '**'
--------
Operation - TRAVERSE
--------
PathKey - 'cat'
--------
Operation - TRAVERSE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixOrExample(t *testing.T) {
	var infix = "a OR b"
	var expectedOutput = `PathKey - 'a'
--------
PathKey - 'b'
--------
Operation - OR
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixOrWithEqualsExample(t *testing.T) {
	var infix = "a==thing OR b==thongs"
	var expectedOutput = `PathKey - 'a'
--------
PathKey - 'thing'
--------
Operation - EQUALS
--------
PathKey - 'b'
--------
PathKey - 'thongs'
--------
Operation - EQUALS
--------
Operation - OR
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixOrWithEqualsPathExample(t *testing.T) {
	var infix = "apples.monkeys==thing OR bogs.bobos==thongs"
	var expectedOutput = `PathKey - 'apples'
--------
PathKey - 'monkeys'
--------
Operation - TRAVERSE
--------
PathKey - 'thing'
--------
Operation - EQUALS
--------
PathKey - 'bogs'
--------
PathKey - 'bobos'
--------
Operation - TRAVERSE
--------
PathKey - 'thongs'
--------
Operation - EQUALS
--------
Operation - OR
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}
