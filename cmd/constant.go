package cmd

var leadingContentPreProcessing = true
var unwrapScalar = true

var writeInplace = false
var outputToJSON = false
var outputFormat = "yaml"
var inputFormat = "yaml"

var xmlAttributePrefix = "+"
var xmlContentName = "+content"
var xmlStrictMode = false

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

var completedSuccessfully = false

var forceExpression = ""

var expressionFile = ""
