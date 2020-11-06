package cmd

import (
	logging "gopkg.in/op/go-logging.v1"
)

var unwrapScalar = true
var writeInplace = false
var outputToJSON = false
var exitStatus = false
var forceColor = false
var forceNoColor = false
var colorsEnabled = false
var indent = 2
var printDocSeparators = true
var nullInput = false
var verbose = false
var version = false
var shellCompletion = ""
var log = logging.MustGetLogger("yq")
