package cmd

import (
	"github.com/mikefarah/yq/v3/pkg/yqlib"
	logging "gopkg.in/op/go-logging.v1"
)

var customTag = ""
var printMode = "v"
var writeInplace = false
var writeScript = ""
var outputToJSON = false
var overwriteFlag = false
var autoCreateFlag = true
var allowEmptyFlag = false
var appendFlag = false
var verbose = false
var version = false
var docIndex = "0"
var log = logging.MustGetLogger("yq")
var lib = yqlib.NewYqLib()
var valueParser = yqlib.NewValueParser()
