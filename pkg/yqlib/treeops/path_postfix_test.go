package treeops

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
		formatted = formatted + path.toString() + "\n--------\n"
	}
	return formatted, nil
}

func TestPostFixTraverseBar(t *testing.T) {
	var infix = ".animals | [.]"
	var expectedOutput = `PathKey - animals
--------
SELF
--------
Operation - COLLECT
--------
Operation - PIPE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixPipeEquals(t *testing.T) {
	var infix = `.animals |   (. == "cat") `
	var expectedOutput = `PathKey - animals
--------
SELF
--------
Value - cat (string)
--------
Operation - EQUALS
--------
Operation - PIPE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixCollect(t *testing.T) {
	var infix = "[.a]"
	var expectedOutput = `PathKey - a
--------
Operation - COLLECT
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSplatSearch(t *testing.T) {
	var infix = `.a | (.[].b == "apple")`
	var expectedOutput = `PathKey - a
--------
PathKey - []
--------
PathKey - b
--------
Operation - PIPE
--------
Value - apple (string)
--------
Operation - EQUALS
--------
Operation - PIPE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixCollectWithExpression(t *testing.T) {
	var infix = `[ (.a == "fred") | (.d, .f)]`
	var expectedOutput = `PathKey - a
--------
Value - fred (string)
--------
Operation - EQUALS
--------
PathKey - d
--------
PathKey - f
--------
Operation - OR
--------
Operation - PIPE
--------
Operation - COLLECT
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixLength(t *testing.T) {
	var infix = ".a | length"
	var expectedOutput = `PathKey - a
--------
Operation - LENGTH
--------
Operation - PIPE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimpleExample(t *testing.T) {
	var infix = ".a"
	var expectedOutput = `PathKey - a
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimplePathExample(t *testing.T) {
	var infix = ".apples.bananas*.cat"
	var expectedOutput = `PathKey - apples
--------
PathKey - bananas*
--------
Operation - PIPE
--------
PathKey - cat
--------
Operation - PIPE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimpleAssign(t *testing.T) {
	var infix = ".a.b |= \"frog\""
	var expectedOutput = `PathKey - a
--------
PathKey - b
--------
Operation - PIPE
--------
Value - frog (string)
--------
Operation - ASSIGN
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimplePathNumbersExample(t *testing.T) {
	var infix = ".apples[0].cat"
	var expectedOutput = `PathKey - apples
--------
PathKey - 0
--------
Operation - PIPE
--------
PathKey - cat
--------
Operation - PIPE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixSimplePathSplatArrayExample(t *testing.T) {
	var infix = ".apples[].cat"
	var expectedOutput = `PathKey - apples
--------
PathKey - []
--------
Operation - PIPE
--------
PathKey - cat
--------
Operation - PIPE
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixOrExample(t *testing.T) {
	var infix = ".a, .b"
	var expectedOutput = `PathKey - a
--------
PathKey - b
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

func TestPostFixEqualsNumberExample(t *testing.T) {
	var infix = ".animal == 3"
	var expectedOutput = `PathKey - animal
--------
Value - 3 (int64)
--------
Operation - EQUALS
--------
`

	actual, err := testExpression(infix)
	if err != nil {
		t.Error(err)
	}

	test.AssertResultComplex(t, expectedOutput, actual)
}

func TestPostFixOrWithEqualsExample(t *testing.T) {
	var infix = ".a==\"thing\", .b==.thongs"
	var expectedOutput = `PathKey - a
--------
Value - thing (string)
--------
Operation - EQUALS
--------
PathKey - b
--------
PathKey - thongs
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
	var infix = ".apples.monkeys==\"thing\", .bogs.bobos==true"
	var expectedOutput = `PathKey - apples
--------
PathKey - monkeys
--------
Operation - PIPE
--------
Value - thing (string)
--------
Operation - EQUALS
--------
PathKey - bogs
--------
PathKey - bobos
--------
Operation - PIPE
--------
Value - true (bool)
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
