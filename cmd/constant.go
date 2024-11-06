package cmd

import "os"

var unwrapScalarFlag = newUnwrapFlag()

var unwrapScalar = false

var writeInplace = false
var outputToJSON = false

var outputFormat = ""

var inputFormat = ""

var exitStatus = false
var indent = 2
var noDocSeparators = false
var nullInput = false
var nulSepOutput = false
var verbose = false
var version = false
var prettyPrint = false

var forceColor = false
var forceNoColor = false
var colorsEnabled = false

func init() {
	// when NO_COLOR environment variable presents and not an empty string the colored output should be disabled;
	// refer to no-color.org
	forceNoColor = os.Getenv("NO_COLOR") != ""
}

// can be either "" (off), "extract" or "process"
var frontMatter = ""

var splitFileExp = ""
var splitFileExpFile = ""

var completedSuccessfully = false

var forceExpression = ""

var expressionFile = ""
