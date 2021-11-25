package cmd

var (
	leadingContentPreProcessing = true
	unwrapScalar                = true
)

var (
	writeInplace = false
	outputToJSON = false
	outputFormat = "yaml"
)

var (
	exitStatus      = false
	forceColor      = false
	forceNoColor    = false
	colorsEnabled   = false
	indent          = 2
	noDocSeparators = false
	nullInput       = false
	verbose         = false
	version         = false
	prettyPrint     = false
)

// can be either "" (off), "extract" or "process"
var frontMatter = ""

var splitFileExp = ""

var completedSuccessfully = false
