package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/mikefarah/yq/v3/pkg/yqlib"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

func readYamlFile(filename string, path string, updateAll bool, docIndexInt int) ([]*yqlib.NodeContext, error) {
	var matchingNodes []*yqlib.NodeContext

	var currentIndex = 0
	var errorReadingStream = readStream(filename, func(decoder *yaml.Decoder) error {
		for {
			var dataBucket yaml.Node
			errorReading := decoder.Decode(&dataBucket)

			if errorReading == io.EOF {
				return handleEOF(updateAll, docIndexInt, currentIndex)
			}
			var errorParsing error
			matchingNodes, errorParsing = appendDocument(matchingNodes, dataBucket, path, updateAll, docIndexInt, currentIndex)
			if errorParsing != nil {
				return errorParsing
			}
			currentIndex = currentIndex + 1
		}
	})
	return matchingNodes, errorReadingStream
}

func handleEOF(updateAll bool, docIndexInt int, currentIndex int) error {
	log.Debugf("done %v / %v", currentIndex, docIndexInt)
	if !updateAll && currentIndex <= docIndexInt {
		return fmt.Errorf("Could not process document index %v as there are only %v document(s)", docIndex, currentIndex)
	}
	return nil
}

func appendDocument(originalMatchingNodes []*yqlib.NodeContext, dataBucket yaml.Node, path string, updateAll bool, docIndexInt int, currentIndex int) ([]*yqlib.NodeContext, error) {
	log.Debugf("processing document %v - requested index %v", currentIndex, docIndexInt)
	yqlib.DebugNode(&dataBucket)
	if !updateAll && currentIndex != docIndexInt {
		return originalMatchingNodes, nil
	}
	log.Debugf("reading %v in document %v", path, currentIndex)
	matchingNodes, errorParsing := lib.Get(&dataBucket, path)
	if errorParsing != nil {
		return nil, errors.Wrapf(errorParsing, "Error reading path in document index %v", currentIndex)
	}
	return append(originalMatchingNodes, matchingNodes...), nil
}

func printValue(node *yaml.Node, cmd *cobra.Command) error {
	if node.Kind == yaml.ScalarNode {
		cmd.Print(node.Value)
		return nil
	}

	bufferedWriter := bufio.NewWriter(cmd.OutOrStdout())
	defer safelyFlush(bufferedWriter)

	var encoder yqlib.Encoder
	if outputToJSON {
		encoder = yqlib.NewJsonEncoder(bufferedWriter)
	} else {
		encoder = yqlib.NewYamlEncoder(bufferedWriter)
	}
	if err := encoder.Encode(node); err != nil {
		return err
	}
	return nil
}

func printResults(matchingNodes []*yqlib.NodeContext, cmd *cobra.Command) error {
	if len(matchingNodes) == 0 {
		log.Debug("no matching results, nothing to print")
		return nil
	}

	for index, mappedDoc := range matchingNodes {
		switch printMode {
		case "p":
			cmd.Print(lib.PathStackToString(mappedDoc.PathStack))
			if index < len(matchingNodes)-1 {
				cmd.Print("\n")
			}
		case "pv", "vp":
			// put it into a node and print that.
			var parentNode = yaml.Node{Kind: yaml.MappingNode}
			parentNode.Content = make([]*yaml.Node, 2)
			parentNode.Content[0] = &yaml.Node{Kind: yaml.ScalarNode, Value: lib.PathStackToString(mappedDoc.PathStack)}
			parentNode.Content[1] = mappedDoc.Node
			if err := printValue(&parentNode, cmd); err != nil {
				return err
			}
		default:
			if err := printValue(mappedDoc.Node, cmd); err != nil {
				return err
			}
			// Printing our Scalars does not print a new line at the end
			// we only want to do that if there are more values (so users can easily script extraction of values in the yaml)
			if index < len(matchingNodes)-1 && mappedDoc.Node.Kind == yaml.ScalarNode {
				cmd.Print("\n")
			}
		}
	}

	return nil
}

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

func mapYamlDecoder(updateData updateDataFn, encoder yqlib.Encoder) yamlDecoderFn {
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

func prefixDocument(updateAll bool, docIndexInt int, currentIndex int, dataBucket *yaml.Node, updateCommand yqlib.UpdateCommand) error {
	if updateAll || currentIndex == docIndexInt {
		log.Debugf("Prefixing document %v", currentIndex)
		yqlib.DebugNode(dataBucket)
		updateCommand.Value = dataBucket.Content[0]
		dataBucket.Content = make([]*yaml.Node, 1)

		newNode := lib.New(updateCommand.Path)
		dataBucket.Content[0] = &newNode

		errorUpdating := lib.Update(dataBucket, updateCommand, true)
		if errorUpdating != nil {
			return errorUpdating
		}
	}
	return nil
}

func updateDoc(inputFile string, updateCommands []yqlib.UpdateCommand, writer io.Writer) error {
	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	var updateData = func(dataBucket *yaml.Node, currentIndex int) error {
		if updateAll || currentIndex == docIndexInt {
			log.Debugf("Updating doc %v", currentIndex)
			for _, updateCommand := range updateCommands {
				log.Debugf("Processing update to Path %v", updateCommand.Path)
				errorUpdating := lib.Update(dataBucket, updateCommand, autoCreateFlag)
				if errorUpdating != nil {
					return errorUpdating
				}
			}
		}
		return nil
	}
	return readAndUpdate(writer, inputFile, updateData)
}

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
		destination = stdOut
		destinationName = "Stdout"
	}

	log.Debugf("Writing to %v from %v", destinationName, inputFile)

	bufferedWriter := bufio.NewWriter(destination)
	defer safelyFlush(bufferedWriter)

	var encoder yqlib.Encoder
	if outputToJSON {
		encoder = yqlib.NewJsonEncoder(bufferedWriter)
	} else {
		encoder = yqlib.NewYamlEncoder(bufferedWriter)
	}
	return readStream(inputFile, mapYamlDecoder(updateData, encoder))
}

type updateCommandParsed struct {
	Command string
	Path    string
	Value   yaml.Node
}

func readUpdateCommands(args []string, expectedArgs int, badArgsMessage string) ([]yqlib.UpdateCommand, error) {
	var updateCommands []yqlib.UpdateCommand = make([]yqlib.UpdateCommand, 0)
	if writeScript != "" {
		var parsedCommands = make([]updateCommandParsed, 0)

		err := readData(writeScript, 0, &parsedCommands)

		if err != nil && err != io.EOF {
			return nil, err
		}

		log.Debugf("Read write commands file '%v'", parsedCommands)
		for index := range parsedCommands {
			parsedCommand := parsedCommands[index]
			updateCommand := yqlib.UpdateCommand{Command: parsedCommand.Command, Path: parsedCommand.Path, Value: &parsedCommand.Value, Overwrite: true}
			updateCommands = append(updateCommands, updateCommand)
		}

		log.Debugf("Read write commands file '%v'", updateCommands)
	} else if len(args) < expectedArgs {
		return nil, errors.New(badArgsMessage)
	} else {
		updateCommands = make([]yqlib.UpdateCommand, 1)
		log.Debug("args %v", args)
		log.Debug("path %v", args[expectedArgs-2])
		log.Debug("Value %v", args[expectedArgs-1])
		updateCommands[0] = yqlib.UpdateCommand{Command: "update", Path: args[expectedArgs-2], Value: valueParser.Parse(args[expectedArgs-1], customTag), Overwrite: true}
	}
	return updateCommands, nil
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
