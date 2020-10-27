package cmd

import (
	logging "gopkg.in/op/go-logging.v1"
)

var customTag = ""
var printMode = "v"
var printLength = false
var unwrapScalar = true
var customStyle = ""
var anchorName = ""
var makeAlias = false
var stripComments = false
var writeInplace = false
var writeScript = ""
var sourceYamlFile = ""
var outputToJSON = false
var exitStatus = false
var prettyPrint = false
var explodeAnchors = false
var forceColor = false
var forceNoColor = false
var colorsEnabled = false
var defaultValue = ""
var indent = 2
var overwriteFlag = false
var autoCreateFlag = true
var arrayMergeStrategyFlag = "update"
var commentsMergeStrategyFlag = "setWhenBlank"
var verbose = false
var version = false
var shellCompletion = ""
var log = logging.MustGetLogger("yq")
