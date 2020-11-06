package cmd

import (
	"container/list"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createEvaluateSequenceCommand() *cobra.Command {
	var cmdEvalSequence = &cobra.Command{
		Use:     "eval-seq [expression] [yaml_file1]...",
		Aliases: []string{"es"},
		Short:   "Apply expression to each document in each yaml file given in sequence",
		Example: `
yq es '.a.b | length' file1.yml file2.yml
yq es < sample.yaml
yq es -n '{"a": "b"}'
`,
		Long: "Evaluate Sequence:\nIterate over each yaml document, apply the expression and print the results, in sequence.",
		RunE: evaluateSequence,
	}
	return cmdEvalSequence
}
func evaluateSequence(cmd *cobra.Command, args []string) error {
	// 0 args, read std in
	// 1 arg, null input, process expression
	// 1 arg, read file in sequence
	// 2+ args, [0] = expression, file the rest

	var matchingNodes *list.List
	var err error
	stat, _ := os.Stdin.Stat()
	pipingStdIn := (stat.Mode() & os.ModeCharDevice) == 0

	switch len(args) {
	case 0:
		if pipingStdIn {
			matchingNodes, err = yqlib.Evaluate("-", "")
		} else {
			cmd.Println(cmd.UsageString())
			return nil
		}
	case 1:
		if nullInput {
			matchingNodes, err = yqlib.EvaluateExpression(args[0])
		} else {
			matchingNodes, err = yqlib.Evaluate(args[0], "")
		}
	}
	cmd.SilenceUsage = true
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()

	fileInfo, _ := os.Stdout.Stat()

	if forceColor || (!forceNoColor && (fileInfo.Mode()&os.ModeCharDevice) != 0) {
		colorsEnabled = true
	}
	printer := yqlib.NewPrinter(outputToJSON, unwrapScalar, colorsEnabled, indent, printDocSeparators)

	return printer.PrintResults(matchingNodes, out)
}
