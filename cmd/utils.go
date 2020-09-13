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
	yaml "gopkg.in/yaml.v3"
)

type readDataFn func(dataBucket *yaml.Node) ([]*yqlib.NodeContext, error)

func createReadFunction(path string) func(*yaml.Node) ([]*yqlib.NodeContext, error) {
	return func(dataBucket *yaml.Node) ([]*yqlib.NodeContext, error) {
		return lib.Get(dataBucket, path)
	}
}

func readYamlFile(filename string, path string, updateAll bool, docIndexInt int) ([]*yqlib.NodeContext, error) {
	return doReadYamlFile(filename, createReadFunction(path), updateAll, docIndexInt)
}

func doReadYamlFile(filename string, readFn readDataFn, updateAll bool, docIndexInt int) ([]*yqlib.NodeContext, error) {
	var matchingNodes []*yqlib.NodeContext

	var currentIndex = 0
	var errorReadingStream = readStream(filename, func(decoder *yaml.Decoder) error {
		for {
			var dataBucket yaml.Node
			errorReading := decoder.Decode(&dataBucket)

			if errorReading == io.EOF {
				return handleEOF(updateAll, docIndexInt, currentIndex)
			} else if errorReading != nil {
				return errorReading
			}

			var errorParsing error
			matchingNodes, errorParsing = appendDocument(matchingNodes, dataBucket, readFn, updateAll, docIndexInt, currentIndex)
			if errorParsing != nil {
				return errorParsing
			}
			if !updateAll && currentIndex == docIndexInt {
				log.Debug("all done")
				return nil
			}
			currentIndex = currentIndex + 1
		}
	})
	return matchingNodes, errorReadingStream
}

func handleEOF(updateAll bool, docIndexInt int, currentIndex int) error {
	log.Debugf("done %v / %v", currentIndex, docIndexInt)
	if !updateAll && currentIndex <= docIndexInt && docIndexInt != 0 {
		return fmt.Errorf("Could not process document index %v as there are only %v document(s)", docIndex, currentIndex)
	}
	return nil
}

func appendDocument(originalMatchingNodes []*yqlib.NodeContext, dataBucket yaml.Node, readFn readDataFn, updateAll bool, docIndexInt int, currentIndex int) ([]*yqlib.NodeContext, error) {
	log.Debugf("processing document %v - requested index %v", currentIndex, docIndexInt)
	yqlib.DebugNode(&dataBucket)
	if !updateAll && currentIndex != docIndexInt {
		return originalMatchingNodes, nil
	}
	log.Debugf("reading in document %v", currentIndex)
	matchingNodes, errorParsing := readFn(&dataBucket)
	if errorParsing != nil {
		return nil, errors.Wrapf(errorParsing, "Error reading path in document index %v", currentIndex)
	}
	return append(originalMatchingNodes, matchingNodes...), nil
}

func lengthOf(node *yaml.Node) int {
	kindToCheck := node.Kind
	if node.Kind == yaml.DocumentNode && len(node.Content) == 1 {
		log.Debugf("length of document node, calculating length of child")
		kindToCheck = node.Content[0].Kind
	}
	switch kindToCheck {
	case yaml.ScalarNode:
		return len(node.Value)
	case yaml.MappingNode:
		return len(node.Content) / 2
	default:
		return len(node.Content)
	}
}

// transforms node before printing, if required
func transformNode(node *yaml.Node) *yaml.Node {
	if printLength {
		return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", lengthOf(node))}
	}
	return node
}

func printNode(node *yaml.Node, writer io.Writer) error {
	var encoder yqlib.Encoder
	if node.Kind == yaml.ScalarNode && unwrapScalar && !outputToJSON {
		return writeString(writer, node.Value+"\n")
	}
	if outputToJSON {
		encoder = yqlib.NewJsonEncoder(writer, prettyPrint, indent)
	} else {
		encoder = yqlib.NewYamlEncoder(writer, indent, colorsEnabled)
	}
	return encoder.Encode(node)
}

func removeComments(matchingNodes []*yqlib.NodeContext) {
	for _, nodeContext := range matchingNodes {
		removeCommentOfNode(nodeContext.Node)
	}
}

func removeCommentOfNode(node *yaml.Node) {
	node.HeadComment = ""
	node.LineComment = ""
	node.FootComment = ""

	for _, child := range node.Content {
		removeCommentOfNode(child)
	}
}

func setStyle(matchingNodes []*yqlib.NodeContext, style yaml.Style) {
	for _, nodeContext := range matchingNodes {
		updateStyleOfNode(nodeContext.Node, style)
	}
}

func updateStyleOfNode(node *yaml.Node, style yaml.Style) {
	node.Style = style

	for _, child := range node.Content {
		updateStyleOfNode(child, style)
	}
}

func writeString(writer io.Writer, txt string) error {
	_, errorWriting := writer.Write([]byte(txt))
	return errorWriting
}

func setIfNotThere(node *yaml.Node, key string, value *yaml.Node) {
	for index := 0; index < len(node.Content); index = index + 2 {
		keyNode := node.Content[index]
		if keyNode.Value == key {
			return
		}
	}
	// need to add it to the map
	mapEntryKey := yaml.Node{Value: key, Kind: yaml.ScalarNode}
	node.Content = append(node.Content, &mapEntryKey)
	node.Content = append(node.Content, value)
}

func applyAlias(node *yaml.Node, alias *yaml.Node) {
	if alias == nil {
		return
	}
	for index := 0; index < len(alias.Content); index = index + 2 {
		keyNode := alias.Content[index]
		log.Debugf("applying alias key %v", keyNode.Value)
		valueNode := alias.Content[index+1]
		setIfNotThere(node, keyNode.Value, valueNode)
	}
}

func explodeNode(node *yaml.Node) error {
	node.Anchor = ""
	switch node.Kind {
	case yaml.SequenceNode, yaml.DocumentNode:
		for index, contentNode := range node.Content {
			log.Debugf("exploding index %v", index)
			errorInContent := explodeNode(contentNode)
			if errorInContent != nil {
				return errorInContent
			}
		}
		return nil
	case yaml.AliasNode:
		log.Debugf("its an alias!")
		if node.Alias != nil {
			node.Kind = node.Alias.Kind
			node.Style = node.Alias.Style
			node.Tag = node.Alias.Tag
			node.Content = node.Alias.Content
			node.Value = node.Alias.Value
			node.Alias = nil
		}
		return nil
	case yaml.MappingNode:
		for index := 0; index < len(node.Content); index = index + 2 {
			keyNode := node.Content[index]
			valueNode := node.Content[index+1]
			log.Debugf("traversing %v", keyNode.Value)
			if keyNode.Value != "<<" {
				errorInContent := explodeNode(valueNode)
				if errorInContent != nil {
					return errorInContent
				}
				errorInContent = explodeNode(keyNode)
				if errorInContent != nil {
					return errorInContent
				}
			} else {
				if valueNode.Kind == yaml.SequenceNode {
					log.Debugf("an alias merge list!")
					for index := len(valueNode.Content) - 1; index >= 0; index = index - 1 {
						aliasNode := valueNode.Content[index]
						applyAlias(node, aliasNode.Alias)
					}
				} else {
					log.Debugf("an alias merge!")
					applyAlias(node, valueNode.Alias)
				}
				node.Content = append(node.Content[:index], node.Content[index+2:]...)
				//replay that index, since the array is shorter now.
				index = index - 2
			}
		}

		return nil
	default:
		return nil
	}
}

func explode(matchingNodes []*yqlib.NodeContext) error {
	log.Debug("exploding nodes")
	for _, nodeContext := range matchingNodes {
		log.Debugf("exploding %v", nodeContext.Head)
		errorExplodingNode := explodeNode(nodeContext.Node)
		if errorExplodingNode != nil {
			return errorExplodingNode
		}
	}
	return nil
}

func printResults(matchingNodes []*yqlib.NodeContext, writer io.Writer) error {
	if prettyPrint {
		setStyle(matchingNodes, 0)
	}

	if stripComments {
		removeComments(matchingNodes)
	}

	//always explode anchors when printing json
	if explodeAnchors || outputToJSON {
		errorExploding := explode(matchingNodes)
		if errorExploding != nil {
			return errorExploding
		}
	}

	bufferedWriter := bufio.NewWriter(writer)
	defer safelyFlush(bufferedWriter)

	if len(matchingNodes) == 0 {
		log.Debug("no matching results, nothing to print")
		if defaultValue != "" {
			return writeString(bufferedWriter, defaultValue)
		}
		return nil
	}
	var errorWriting error

	var arrayCollection = yaml.Node{Kind: yaml.SequenceNode}

	for _, mappedDoc := range matchingNodes {
		switch printMode {
		case "p":
			errorWriting = writeString(bufferedWriter, lib.PathStackToString(mappedDoc.PathStack)+"\n")
			if errorWriting != nil {
				return errorWriting
			}
		case "pv", "vp":
			// put it into a node and print that.
			var parentNode = yaml.Node{Kind: yaml.MappingNode}
			parentNode.Content = make([]*yaml.Node, 2)
			parentNode.Content[0] = &yaml.Node{Kind: yaml.ScalarNode, Value: lib.PathStackToString(mappedDoc.PathStack)}
			parentNode.Content[1] = transformNode(mappedDoc.Node)
			if collectIntoArray {
				arrayCollection.Content = append(arrayCollection.Content, &parentNode)
			} else if err := printNode(&parentNode, bufferedWriter); err != nil {
				return err
			}
		default:
			if collectIntoArray {
				arrayCollection.Content = append(arrayCollection.Content, mappedDoc.Node)
			} else if err := printNode(transformNode(mappedDoc.Node), bufferedWriter); err != nil {
				return err
			}
		}
	}

	if collectIntoArray {
		if err := printNode(transformNode(&arrayCollection), bufferedWriter); err != nil {
			return err
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

func isNullDocument(dataBucket *yaml.Node) bool {
	return dataBucket.Kind == yaml.DocumentNode && (len(dataBucket.Content) == 0 ||
		dataBucket.Content[0].Kind == yaml.ScalarNode && dataBucket.Content[0].Tag == "!!null")
}

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

			if errorReading == io.EOF && docIndexInt == 0 && currentIndex == 0 {
				//empty document, lets just make one
				dataBucket = yaml.Node{Kind: yaml.DocumentNode, Content: make([]*yaml.Node, 1)}
				child := yaml.Node{Kind: yaml.MappingNode}
				dataBucket.Content[0] = &child
			} else if isNullDocument(&dataBucket) && (updateAll || docIndexInt == currentIndex) {
				child := yaml.Node{Kind: yaml.MappingNode}
				dataBucket.Content[0] = &child
			} else if errorReading == io.EOF {
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

			if prettyPrint {
				updateStyleOfNode(&dataBucket, 0)
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
	var completedSuccessfully = false
	if writeInplace {
		info, err := os.Stat(inputFile)
		if err != nil {
			return err
		}
		// mkdir temp dir as some docker images does not have temp dir
		_, err = os.Stat(os.TempDir())
		if os.IsNotExist(err) {
			err = os.Mkdir(os.TempDir(), 0700)
			if err != nil {
				return err
			}
		} else if err != nil {
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
			if completedSuccessfully {
				safelyRenameFile(tempFile.Name(), inputFile)
			}
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
		encoder = yqlib.NewJsonEncoder(bufferedWriter, prettyPrint, indent)
	} else {
		encoder = yqlib.NewYamlEncoder(bufferedWriter, indent, colorsEnabled)
	}

	var errorProcessing = readStream(inputFile, mapYamlDecoder(updateData, encoder))
	completedSuccessfully = errorProcessing == nil
	return errorProcessing
}

type updateCommandParsed struct {
	Command string
	Path    string
	Value   yaml.Node
}

func readUpdateCommands(args []string, expectedArgs int, badArgsMessage string, allowNoValue bool) ([]yqlib.UpdateCommand, error) {
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
	} else if sourceYamlFile != "" && len(args) == expectedArgs-1 {
		log.Debugf("Reading value from %v", sourceYamlFile)
		var value yaml.Node
		err := readData(sourceYamlFile, 0, &value)
		if err != nil && err != io.EOF {
			return nil, err
		}
		log.Debug("args %v", args[expectedArgs-2])
		updateCommands = make([]yqlib.UpdateCommand, 1)
		updateCommands[0] = yqlib.UpdateCommand{Command: "update", Path: args[expectedArgs-2], Value: value.Content[0], Overwrite: true}
	} else if len(args) == expectedArgs {
		updateCommands = make([]yqlib.UpdateCommand, 1)
		log.Debug("args %v", args)
		log.Debug("path %v", args[expectedArgs-2])
		log.Debug("Value %v", args[expectedArgs-1])
		value := valueParser.Parse(args[expectedArgs-1], customTag, customStyle, anchorName, makeAlias)
		updateCommands[0] = yqlib.UpdateCommand{Command: "update", Path: args[expectedArgs-2], Value: value, Overwrite: true, CommentsMergeStrategy: yqlib.IgnoreCommentsMergeStrategy}
	} else if len(args) == expectedArgs-1 && allowNoValue {
		// don't update the value
		updateCommands = make([]yqlib.UpdateCommand, 1)
		log.Debug("args %v", args)
		log.Debug("path %v", args[expectedArgs-2])
		updateCommands[0] = yqlib.UpdateCommand{Command: "update", Path: args[expectedArgs-2], Value: valueParser.Parse("", customTag, customStyle, anchorName, makeAlias), Overwrite: true, DontUpdateNodeValue: true}
	} else {
		return nil, errors.New(badArgsMessage)
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
