// Example of using yq as a Go library.
// Run with: go run ./examples/library
package main

import (
	"log"
	"os"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

func main() {
	decoder := yqlib.NewYamlDecoder(yqlib.NewDefaultYamlPreferences())
	if err := decoder.Init(strings.NewReader("a:\n  b: val\n")); err != nil {
		log.Fatal(err)
	}

	node, err := decoder.Decode()
	if err != nil {
		log.Fatal(err)
	}

	results, err := yqlib.NewAllAtOnceEvaluator().EvaluateNodes(`.a.b = "updated"`, node)
	if err != nil {
		log.Fatal(err)
	}

	encoder := yqlib.NewYamlEncoder(yqlib.NewDefaultYamlPreferences())
	printer := yqlib.NewPrinter(encoder, yqlib.NewSinglePrinterWriter(os.Stdout))
	if err := printer.PrintResults(results); err != nil {
		log.Fatal(err)
	}
}
