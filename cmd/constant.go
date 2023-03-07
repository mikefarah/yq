package cmd

var unwrapScalarFlag = newUnwrapFlag()

var unwrapScalar = false

var writeInplace = false
var outputToJSON = false
var outputFormat = "yaml"
var inputFormatDefault = "yaml"
var inputFormat = ""

var exitStatus = false
var forceColor = false
var forceNoColor = false
var colorsEnabled = false
var indent = 2
var noDocSeparators = false
var nullInput = false
var verbose = false
var version = false
var prettyPrint = false

// can be either "" (off), "extract" or "process"
var frontMatter = ""

var splitFileExp = ""
var splitFileExpFile = ""

var completedSuccessfully = false

var forceExpression = ""

var expressionFile = ""
