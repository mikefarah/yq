package cmd

import (
	"errors"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createEvaluateAllCommand() *cobra.Command {
	var cmdEvalAll = &cobra.Command{
		Use:     "eval-all [expression] [yaml_file1]...",
		Aliases: []string{"ea"},
		Short:   "Loads _all_ yaml documents of _all_ yaml files and runs expression once",
		Example: `
# Merge f2.yml into f1.yml (inplace)
yq eval-all --inplace 'select(fileIndex == 0) * select(fileIndex == 1)' f1.yml f2.yml
## the same command and expression using shortened names:
yq ea -i 'select(fi == 0) * select(fi == 1)' f1.yml f2.yml


# Merge all given files
yq ea '. as $item ireduce ({}; . * $item )' file1.yml file2.yml ...

# Pipe from STDIN
## use '-' as a filename to pipe from STDIN
cat file2.yml | yq ea '.a.b' file1.yml - file3.yml
`,
		Long: `yq is a portable command-line YAML processor (https://github.com/mikefarah/yq/) 
See https://mikefarah.gitbook.io/yq/ for detailed documentation and examples.

## Evaluate All ##
This command loads _all_ yaml documents of _all_ yaml files and runs expression once
Useful when you need to run an expression across several yaml documents or files (like merge).
Note that it consumes more memory than eval.
`,
		RunE: evaluateAll,
	}
	return cmdEvalAll
}
func evaluateAll(cmd *cobra.Command, args []string) (cmdError error) {
	// 0 args, read std in
	// 1 arg, null input, process expression
	// 1 arg, read file in sequence
	// 2+ args, [0] = expression, file the rest

	var err error

	firstFileIndex, err := initCommand(cmd, args)
	if err != nil {
		return err
	}

	stat, _ := os.Stdin.Stat()
	pipingStdIn := (stat.Mode() & os.ModeCharDevice) == 0
	yqlib.GetLogger().Debug("pipingStdIn: %v", pipingStdIn)

	yqlib.GetLogger().Debug("stat.Mode(): %v", stat.Mode())
	yqlib.GetLogger().Debug("ModeDir: %v", stat.Mode()&os.ModeDir)
	yqlib.GetLogger().Debug("ModeAppend: %v", stat.Mode()&os.ModeAppend)
	yqlib.GetLogger().Debug("ModeExclusive: %v", stat.Mode()&os.ModeExclusive)
	yqlib.GetLogger().Debug("ModeTemporary: %v", stat.Mode()&os.ModeTemporary)
	yqlib.GetLogger().Debug("ModeSymlink: %v", stat.Mode()&os.ModeSymlink)
	yqlib.GetLogger().Debug("ModeDevice: %v", stat.Mode()&os.ModeDevice)
	yqlib.GetLogger().Debug("ModeNamedPipe: %v", stat.Mode()&os.ModeNamedPipe)
	yqlib.GetLogger().Debug("ModeSocket: %v", stat.Mode()&os.ModeSocket)
	yqlib.GetLogger().Debug("ModeSetuid: %v", stat.Mode()&os.ModeSetuid)
	yqlib.GetLogger().Debug("ModeSetgid: %v", stat.Mode()&os.ModeSetgid)
	yqlib.GetLogger().Debug("ModeCharDevice: %v", stat.Mode()&os.ModeCharDevice)
	yqlib.GetLogger().Debug("ModeSticky: %v", stat.Mode()&os.ModeSticky)
	yqlib.GetLogger().Debug("ModeIrregular: %v", stat.Mode()&os.ModeIrregular)

	out := cmd.OutOrStdout()

	if writeInplace {
		// only use colors if its forced
		colorsEnabled = forceColor
		writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(args[firstFileIndex])
		out, err = writeInPlaceHandler.CreateTempFile()
		if err != nil {
			return err
		}
		// need to indirectly call the function so  that completedSuccessfully is
		// passed when we finish execution as opposed to now
		defer func() {
			if cmdError == nil {
				cmdError = writeInPlaceHandler.FinishWriteInPlace(completedSuccessfully)
			}
		}()
	}

	format, err := yqlib.OutputFormatFromString(outputFormat)
	if err != nil {
		return err
	}

	decoder, err := configureDecoder()
	if err != nil {
		return err
	}

	printerWriter, err := configurePrinterWriter(format, out)
	if err != nil {
		return err
	}
	encoder := configureEncoder(format)

	printer := yqlib.NewPrinter(encoder, printerWriter)

	if frontMatter != "" {
		frontMatterHandler := yqlib.NewFrontMatterHandler(args[firstFileIndex])
		err = frontMatterHandler.Split()
		if err != nil {
			return err
		}
		args[firstFileIndex] = frontMatterHandler.GetYamlFrontMatterFilename()

		if frontMatter == "process" {
			reader := frontMatterHandler.GetContentReader()
			printer.SetAppendix(reader)
			defer yqlib.SafelyCloseReader(reader)
		}
		defer frontMatterHandler.CleanUp()
	}

	allAtOnceEvaluator := yqlib.NewAllAtOnceEvaluator()

	expression, args, err := processArgs(pipingStdIn, args)
	if err != nil {
		return err
	}
	yqlib.GetLogger().Debugf("processed args: %v", args)

	switch len(args) {
	case 0:
		if nullInput {
			err = yqlib.NewStreamEvaluator().EvaluateNew(processExpression(expression), printer, "")
		} else {
			cmd.Println(cmd.UsageString())
			return nil
		}
	default:
		err = allAtOnceEvaluator.EvaluateFiles(processExpression(expression), args, printer, leadingContentPreProcessing, decoder)
	}

	completedSuccessfully = err == nil

	if err == nil && exitStatus && !printer.PrintedAnything() {
		return errors.New("no matches found")
	}

	return err
}
