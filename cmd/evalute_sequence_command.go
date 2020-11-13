package cmd

import (
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createEvaluateSequenceCommand() *cobra.Command {
	var cmdEvalSequence = &cobra.Command{
		Use:     "eval [expression] [yaml_file1]...",
		Aliases: []string{"e"},
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

	var err error
	stat, _ := os.Stdin.Stat()
	pipingStdIn := (stat.Mode() & os.ModeCharDevice) == 0

	out := cmd.OutOrStdout()

	fileInfo, _ := os.Stdout.Stat()

	if forceColor || (!forceNoColor && (fileInfo.Mode()&os.ModeCharDevice) != 0) {
		colorsEnabled = true
	}
	printer := yqlib.NewPrinter(out, outputToJSON, unwrapScalar, colorsEnabled, indent, !noDocSeparators)

	switch len(args) {
	case 0:
		if pipingStdIn {
			err = yqlib.EvaluateFileStreamsSequence("", []string{"-"}, printer)
		} else {
			cmd.Println(cmd.UsageString())
			return nil
		}
	case 1:
		if nullInput {
			err = yqlib.EvaluateAllFileStreams(args[0], []string{}, printer)
		} else {
			err = yqlib.EvaluateFileStreamsSequence("", []string{args[0]}, printer)
		}
	default:
		err = yqlib.EvaluateFileStreamsSequence(args[0], args[1:], printer)
	}

	cmd.SilenceUsage = true
	return err
}
