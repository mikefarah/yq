package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/mikefarah/yq/v3/pkg/marshal"
	"github.com/mikefarah/yq/v3/pkg/yqlib"

	errors "github.com/pkg/errors"

	"github.com/spf13/cobra"
	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

var rawOutput = false
var trimOutput = true
var writeInplace = false
var writeScript = ""
var outputToJSON = false
var overwriteFlag = false
var allowEmptyFlag = false
var appendFlag = false
var verbose = false
var version = false
var docIndex = "0"
var log = logging.MustGetLogger("yq")
var lib = yqlib.NewYqLib(log)
var jsonConverter = marshal.NewJsonConverter()
var yamlConverter = marshal.NewYamlConverter()
var valueParser = yqlib.NewValueParser()

func main() {
	cmd := newCommandCLI()
	if err := cmd.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func newCommandCLI() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "yq",
		Short: "yq is a lightweight and portable command-line YAML processor.",
		Long:  `yq is a lightweight and portable command-line YAML processor. It aims to be the jq or sed of yaml files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				cmd.Print(GetVersionDisplay())
				return nil
			}
			cmd.Println(cmd.UsageString())

			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var format = logging.MustStringFormatter(
				`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
			)
			var backend = logging.AddModuleLevel(
				logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

			if verbose {
				backend.SetLevel(logging.DEBUG, "")
			} else {
				backend.SetLevel(logging.ERROR, "")
			}

			logging.SetBackend(backend)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&trimOutput, "trim", "t", true, "trim yaml output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and quit")

	rootCmd.AddCommand(
		createReadCmd(),
		createWriteCmd(),
		// createPrefixCmd(),
		// createDeleteCmd(),
		// createNewCmd(),
		// createMergeCmd(),
	)
	rootCmd.SetOutput(os.Stdout)

	return rootCmd
}

func createReadCmd() *cobra.Command {
	var cmdRead = &cobra.Command{
		Use:     "read [yaml_file] [path]",
		Aliases: []string{"r"},
		Short:   "yq r [--doc/-d index] sample.yaml a.b.c",
		Example: `
yq read things.yaml a.b.c
yq r - a.b.c (reads from stdin)
yq r things.yaml a.*.c
yq r -d1 things.yaml a.array[0].blah
yq r things.yaml a.array[*].blah
yq r -- things.yaml --key-starting-with-dashes
      `,
		Long: "Outputs the value of the given path in the yaml file to STDOUT",
		RunE: readProperty,
	}
	cmdRead.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	cmdRead.PersistentFlags().BoolVarP(&rawOutput, "raw", "r", false, "raw yaml output - prints out values instead of yaml")
	cmdRead.PersistentFlags().BoolVarP(&outputToJSON, "tojson", "j", false, "output as json")
	return cmdRead
}

func createWriteCmd() *cobra.Command {
	var cmdWrite = &cobra.Command{
		Use:     "write [yaml_file] [path] [value]",
		Aliases: []string{"w"},
		Short:   "yq w [--inplace/-i] [--script/-s script_file] [--doc/-d index] sample.yaml a.b.c newValue",
		Example: `
yq write things.yaml a.b.c cat
yq write --inplace -- things.yaml a.b.c --cat
yq w -i things.yaml a.b.c cat
yq w --script update_script.yaml things.yaml
yq w -i -s update_script.yaml things.yaml
yq w --doc 2 things.yaml a.b.d[+] foo
yq w -d2 things.yaml a.b.d[+] foo
      `,
		Long: `Updates the yaml file w.r.t the given path and value.
Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.

Append value to array adds the value to the end of array.

Update Scripts:
Note that you can give an update script to perform more sophisticated updated. Update script
format is a yaml map where the key is the path and the value is..well the value. e.g.:
---
a.b.c: true,
a.b.e:
  - name: bob
`,
		RunE: writeProperty,
	}
	cmdWrite.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
	cmdWrite.PersistentFlags().StringVarP(&writeScript, "script", "s", "", "yaml script for updating yaml")
	cmdWrite.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	return cmdWrite
}

// func createPrefixCmd() *cobra.Command {
// 	var cmdWrite = &cobra.Command{
// 		Use:     "prefix [yaml_file] [path]",
// 		Aliases: []string{"p"},
// 		Short:   "yq p [--inplace/-i] [--doc/-d index] sample.yaml a.b.c",
// 		Example: `
// yq prefix things.yaml a.b.c
// yq prefix --inplace things.yaml a.b.c
// yq prefix --inplace -- things.yaml --key-starting-with-dash
// yq p -i things.yaml a.b.c
// yq p --doc 2 things.yaml a.b.d
// yq p -d2 things.yaml a.b.d
//       `,
// 		Long: `Prefixes w.r.t to the yaml file at the given path.
// Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.
// `,
// 		RunE: prefixProperty,
// 	}
// 	cmdWrite.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
// 	cmdWrite.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
// 	return cmdWrite
// }

// func createDeleteCmd() *cobra.Command {
// 	var cmdDelete = &cobra.Command{
// 		Use:     "delete [yaml_file] [path]",
// 		Aliases: []string{"d"},
// 		Short:   "yq d [--inplace/-i] [--doc/-d index] sample.yaml a.b.c",
// 		Example: `
// yq delete things.yaml a.b.c
// yq delete --inplace things.yaml a.b.c
// yq delete --inplace -- things.yaml --key-starting-with-dash
// yq d -i things.yaml a.b.c
// yq d things.yaml a.b.c
// 	`,
// 		Long: `Deletes the given path from the YAML file.
// Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.
// `,
// 		RunE: deleteProperty,
// 	}
// 	cmdDelete.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
// 	cmdDelete.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
// 	return cmdDelete
// }

// func createNewCmd() *cobra.Command {
// 	var cmdNew = &cobra.Command{
// 		Use:     "new [path] [value]",
// 		Aliases: []string{"n"},
// 		Short:   "yq n [--script/-s script_file] a.b.c newValue",
// 		Example: `
// yq new a.b.c cat
// yq n a.b.c cat
// yq n -- --key-starting-with-dash cat
// yq n --script create_script.yaml
//       `,
// 		Long: `Creates a new yaml w.r.t the given path and value.
// Outputs to STDOUT

// Create Scripts:
// Note that you can give a create script to perform more sophisticated yaml. This follows the same format as the update script.
// `,
// 		RunE: newProperty,
// 	}
// 	cmdNew.PersistentFlags().StringVarP(&writeScript, "script", "s", "", "yaml script for updating yaml")
// 	return cmdNew
// }

// func createMergeCmd() *cobra.Command {
// 	var cmdMerge = &cobra.Command{
// 		Use:     "merge [initial_yaml_file] [additional_yaml_file]...",
// 		Aliases: []string{"m"},
// 		Short:   "yq m [--inplace/-i] [--doc/-d index] [--overwrite/-x] [--append/-a] sample.yaml sample2.yaml",
// 		Example: `
// yq merge things.yaml other.yaml
// yq merge --inplace things.yaml other.yaml
// yq m -i things.yaml other.yaml
// yq m --overwrite things.yaml other.yaml
// yq m -i -x things.yaml other.yaml
// yq m -i -a things.yaml other.yaml
//       `,
// 		Long: `Updates the yaml file by adding/updating the path(s) and value(s) from additional yaml file(s).
// Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.

// If overwrite flag is set then existing values will be overwritten using the values from each additional yaml file.
// If append flag is set then existing arrays will be merged with the arrays from each additional yaml file.

// Note that if you set both flags only overwrite will take effect.
// `,
// 		RunE: mergeProperties,
// 	}
// 	cmdMerge.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
// 	cmdMerge.PersistentFlags().BoolVarP(&overwriteFlag, "overwrite", "x", false, "update the yaml file by overwriting existing values")
// 	cmdMerge.PersistentFlags().BoolVarP(&appendFlag, "append", "a", false, "update the yaml file by appending array values")
// 	cmdMerge.PersistentFlags().BoolVarP(&allowEmptyFlag, "allow-empty", "e", false, "allow empty yaml files")
// 	cmdMerge.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
// 	return cmdMerge
// }

func readProperty(cmd *cobra.Command, args []string) error {
	var path = ""

	if len(args) < 1 {
		return errors.New("Must provide filename")
	} else if len(args) > 1 {
		path = args[1]
	}

	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	var mappedDocs []*yaml.Node
	var dataBucket yaml.Node
	var currentIndex = 0
	var errorReadingStream = readStream(args[0], func(decoder *yaml.Decoder) error {
		for {
			errorReading := decoder.Decode(&dataBucket)
			if errorReading == io.EOF {
				log.Debugf("done %v / %v", currentIndex, docIndexInt)
				if !updateAll && currentIndex <= docIndexInt {
					return fmt.Errorf("asked to process document index %v but there are only %v document(s)", docIndex, currentIndex)
				}
				return nil
			}
			log.Debugf("processing document %v - requested index %v", currentIndex, docIndexInt)
			if updateAll || currentIndex == docIndexInt {
				log.Debugf("reading %v in document %v", path, currentIndex)
				mappedDoc, errorParsing := lib.Get(&dataBucket, path)
				if errorParsing != nil {
					return errors.Wrapf(errorParsing, "Error reading path in document index %v", currentIndex)
				}
				mappedDocs = append(mappedDocs, mappedDoc)
			}
			currentIndex = currentIndex + 1
		}
	})

	if errorReadingStream != nil {
		return errorReadingStream
	}

	var encoder = yaml.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent(2)
	var err error

	if rawOutput {
		for _, mappedDoc := range mappedDocs {
			if mappedDoc != nil {
				cmd.Println(mappedDoc.Value)
			}
		}
	} else if !updateAll {
		err = encoder.Encode(mappedDocs[0])
	} else {
		err = encoder.Encode(&yaml.Node{Kind: yaml.SequenceNode, Content: mappedDocs})
	}
	if err != nil {
		return err
	}
	encoder.Close()
	return nil
}

// func newProperty(cmd *cobra.Command, args []string) error {
// 	updatedData, err := newYaml(args)
// 	if err != nil {
// 		return err
// 	}
// 	dataStr, err := toString(updatedData)
// 	if err != nil {
// 		return err
// 	}
// 	cmd.Println(dataStr)
// 	return nil
// }

// func newYaml(args []string) (interface{}, error) {
// 	var writeCommands, writeCommandsError = readWriteCommands(args, 2, "Must provide <path_to_update> <value>")
// 	if writeCommandsError != nil {
// 		return nil, writeCommandsError
// 	}

// 	var dataBucket interface{}
// 	var isArray = strings.HasPrefix(writeCommands[0].Key.(string), "[")
// 	if isArray {
// 		dataBucket = make([]interface{}, 0)
// 	} else {
// 		dataBucket = make(yaml.MapSlice, 0)
// 	}

// 	for _, entry := range writeCommands {
// 		path := entry.Key.(string)
// 		value := entry.Value
// 		log.Debugf("setting %v to %v", path, value)
// 		dataBucket = lib.WritePath(dataBucket, path, value)
// 	}

// 	return dataBucket, nil
// }

func parseDocumentIndex() (bool, int, error) {
	if docIndex == "*" {
		return true, -1, nil
	}
	docIndexInt64, err := strconv.ParseInt(docIndex, 10, 32)
	if err != nil {
		return false, -1, errors.Wrapf(err, "Document index %v is not a integer or *", docIndex)
	}
	return false, int(docIndexInt64), nil
}

type updateDataFn func(dataBucket *yaml.Node, currentIndex int) error

func mapYamlDecoder(updateData updateDataFn, encoder *yaml.Encoder) yamlDecoderFn {
	return func(decoder *yaml.Decoder) error {
		var dataBucket yaml.Node
		var errorReading error
		var errorWriting error
		var errorUpdating error
		var currentIndex = 0

		var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
		if errorParsingDocIndex != nil {
			return errorParsingDocIndex
		}

		for {
			log.Debugf("Read doc %v", currentIndex)
			errorReading = decoder.Decode(&dataBucket)

			if errorReading == io.EOF {
				if !updateAll && currentIndex <= docIndexInt {
					return fmt.Errorf("asked to process document index %v but there are only %v document(s)", docIndex, currentIndex)
				}
				return nil
			} else if errorReading != nil {
				return errors.Wrapf(errorReading, "Error reading document at index %v, %v", currentIndex, errorReading)
			}
			errorUpdating = updateData(&dataBucket, currentIndex)
			if errorUpdating != nil {
				return errors.Wrapf(errorUpdating, "Error updating document at index %v", currentIndex)
			}

			errorWriting = encoder.Encode(&dataBucket)

			if errorWriting != nil {
				return errors.Wrapf(errorWriting, "Error writing document at index %v, %v", currentIndex, errorWriting)
			}
			currentIndex = currentIndex + 1
		}
	}
}

func writeProperty(cmd *cobra.Command, args []string) error {
	var writeCommands, writeCommandsError = readWriteCommands(args, 3, "Must provide <filename> <path_to_update> <value>")
	if writeCommandsError != nil {
		return writeCommandsError
	}
	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	var updateData = func(dataBucket *yaml.Node, currentIndex int) error {
		if updateAll || currentIndex == docIndexInt {
			log.Debugf("Updating doc %v", currentIndex)
			for _, entry := range writeCommands {
				path := entry.Key
				changesToApply := entry.Value
				var paths = parsePath(path)

				errorUpdating := updateChild(dataBucket, paths, changesToApply)
				if errorUpdating != nil {
					return errorUpdating
				}
			}
		}
		return nil
	}
	return readAndUpdate(cmd.OutOrStdout(), args[0], updateData)
}

// func prefixProperty(cmd *cobra.Command, args []string) error {
// 	if len(args) != 2 {
// 		return errors.New("Must provide <filename> <prefixed_path>")
// 	}
// 	prefixPath := args[1]

// 	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
// 	if errorParsingDocIndex != nil {
// 		return errorParsingDocIndex
// 	}

// 	var updateData = func(dataBucket interface{}, currentIndex int) (interface{}, error) {

// 		if updateAll || currentIndex == docIndexInt {
// 			log.Debugf("Prefixing %v to doc %v", prefixPath, currentIndex)
// 			var mapDataBucket = lib.PrefixPath(dataBucket, prefixPath)
// 			return mapDataBucket, nil
// 		}
// 		return dataBucket, nil
// 	}
// 	return readAndUpdate(cmd.OutOrStdout(), args[0], updateData)
// }

func readAndUpdate(stdOut io.Writer, inputFile string, updateData updateDataFn) error {
	var destination io.Writer
	var destinationName string
	if writeInplace {
		info, err := os.Stat(inputFile)
		if err != nil {
			return err
		}
		tempFile, err := ioutil.TempFile("", "temp")
		if err != nil {
			return err
		}
		destinationName = tempFile.Name()
		err = os.Chmod(destinationName, info.Mode())
		if err != nil {
			return err
		}
		destination = tempFile
		defer func() {
			safelyCloseFile(tempFile)
			safelyRenameFile(tempFile.Name(), inputFile)
		}()
	} else {
		var writer = bufio.NewWriter(stdOut)
		destination = writer
		destinationName = "Stdout"
		defer safelyFlush(writer)
	}
	var encoder = yaml.NewEncoder(destination)
	encoder.SetIndent(2)
	log.Debugf("Writing to %v from %v", destinationName, inputFile)
	return readStream(inputFile, mapYamlDecoder(updateData, encoder))
}

// func deleteProperty(cmd *cobra.Command, args []string) error {
// 	if len(args) < 2 {
// 		return errors.New("Must provide <filename> <path_to_delete>")
// 	}
// 	var deletePath = args[1]
// 	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
// 	if errorParsingDocIndex != nil {
// 		return errorParsingDocIndex
// 	}

// 	var updateData = func(dataBucket interface{}, currentIndex int) (interface{}, error) {
// 		if updateAll || currentIndex == docIndexInt {
// 			log.Debugf("Deleting path in doc %v", currentIndex)
// 			return lib.DeletePath(dataBucket, deletePath)
// 		}
// 		return dataBucket, nil
// 	}

// 	return readAndUpdate(cmd.OutOrStdout(), args[0], updateData)
// }

// func mergeProperties(cmd *cobra.Command, args []string) error {
// 	if len(args) < 2 {
// 		return errors.New("Must provide at least 2 yaml files")
// 	}
// 	var input = args[0]
// 	var filesToMerge = args[1:]
// 	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
// 	if errorParsingDocIndex != nil {
// 		return errorParsingDocIndex
// 	}

// 	var updateData = func(dataBucket interface{}, currentIndex int) (interface{}, error) {
// 		if updateAll || currentIndex == docIndexInt {
// 			log.Debugf("Merging doc %v", currentIndex)
// 			var mergedData map[interface{}]interface{}
// 			// merge only works for maps, so put everything in a temporary
// 			// map
// 			var mapDataBucket = make(map[interface{}]interface{})
// 			mapDataBucket["root"] = dataBucket
// 			if err := lib.Merge(&mergedData, mapDataBucket, overwriteFlag, appendFlag); err != nil {
// 				return nil, err
// 			}
// 			for _, f := range filesToMerge {
// 				var fileToMerge interface{}
// 				if err := readData(f, 0, &fileToMerge); err != nil {
// 					if allowEmptyFlag && err == io.EOF {
// 						continue
// 					}
// 					return nil, err
// 				}
// 				mapDataBucket["root"] = fileToMerge
// 				if err := lib.Merge(&mergedData, mapDataBucket, overwriteFlag, appendFlag); err != nil {
// 					return nil, err
// 				}
// 			}
// 			return mergedData["root"], nil
// 		}
// 		return dataBucket, nil
// 	}
// 	return readAndUpdate(cmd.OutOrStdout(), input, updateData)
// }

type rawWriteCommand struct {
	// Command string TODO
	Key   string
	Value yaml.Node
}

func readWriteCommands(args []string, expectedArgs int, badArgsMessage string) ([]rawWriteCommand, error) {
	var writeCommands []rawWriteCommand
	if writeScript != "" {
		var rawCommands yaml.Node
		if err := readData(writeScript, 0, &rawCommands); err != nil {
			return nil, err
		}
		log.Debugf("Read write commands file '%v'", rawCommands)
		var key string
		for index, content := range rawCommands.Content[0].Content {
			if index%2 == 0 { // must be the key
				key = content.Value
			} else { // its the value
				writeCommands = append(writeCommands, rawWriteCommand{Key: key, Value: *content})
			}
		}
		log.Debugf("Read write commands '%v'", writeCommands)
	} else if len(args) < expectedArgs {
		return nil, errors.New(badArgsMessage)
	} else {
		writeCommands = make([]rawWriteCommand, 1)
		writeCommands[0] = rawWriteCommand{Key: args[expectedArgs-2], Value: parseValue(args[expectedArgs-1])}
	}
	return writeCommands, nil
}

func parseValue(argument string) yaml.Node {
	var err interface{}
	var tag = customTag

	var inQuotes = len(argument) > 0 && argument[0] == '"'
	if tag == "" && !inQuotes {

		_, err = strconv.ParseBool(argument)
		if err == nil {
			tag = "!!bool"
		}
		_, err = strconv.ParseFloat(argument, 64)
		if err == nil {
			tag = "!!float"
		}
		_, err = strconv.ParseInt(argument, 10, 64)
		if err == nil {
			tag = "!!int"
		}

		if argument == "null" {
			tag = "!!null"
		}
		if argument == "[]" {
			return yaml.Node{Tag: "!!seq", Kind: yaml.SequenceNode}
		}
	}
	log.Debugf("Updating node to value '%v', tag: '%v'", argument, tag)
	return yaml.Node{Value: argument, Tag: tag, Kind: yaml.ScalarNode}
}

func toString(context interface{}) (string, error) {
	if outputToJSON {
		return jsonConverter.JsonToString(context)
	}
	return yamlConverter.YamlToString(context, trimOutput)
}

func safelyRenameFile(from string, to string) {
	if renameError := os.Rename(from, to); renameError != nil {
		log.Debugf("Error renaming from %v to %v, attempting to copy contents", from, to)
		log.Debug(renameError.Error())
		// can't do this rename when running in docker to a file targeted in a mounted volume,
		// so gracefully degrade to copying the entire contents.
		if copyError := copyFileContents(from, to); copyError != nil {
			log.Errorf("Failed copying from %v to %v", from, to)
			log.Error(copyError.Error())
		} else {
			removeErr := os.Remove(from)
			if removeErr != nil {
				log.Errorf("failed removing original file: %s", from)
			}
		}
	}
}

// thanks https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src) // nolint gosec
	if err != nil {
		return err
	}
	defer safelyCloseFile(in)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer safelyCloseFile(out)
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func safelyFlush(writer *bufio.Writer) {
	if err := writer.Flush(); err != nil {
		log.Error("Error flushing writer!")
		log.Error(err.Error())
	}

}
func safelyCloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Error("Error closing file!")
		log.Error(err.Error())
	}
}

type yamlDecoderFn func(*yaml.Decoder) error

func readStream(filename string, yamlDecoder yamlDecoderFn) error {
	if filename == "" {
		return errors.New("Must provide filename")
	}

	var stream io.Reader
	if filename == "-" {
		stream = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(filename) // nolint gosec
		if err != nil {
			return err
		}
		defer safelyCloseFile(file)
		stream = file
	}
	return yamlDecoder(yaml.NewDecoder(stream))
}

func readData(filename string, indexToRead int, parsedData interface{}) error {
	return readStream(filename, func(decoder *yaml.Decoder) error {
		for currentIndex := 0; currentIndex < indexToRead; currentIndex++ {
			errorSkipping := decoder.Decode(parsedData)
			if errorSkipping != nil {
				return errors.Wrapf(errorSkipping, "Error processing document at index %v, %v", currentIndex, errorSkipping)
			}
		}
		return decoder.Decode(parsedData)
	})
}
