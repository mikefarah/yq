package cmd

import (
	"bufio"
	"container/list"
	"errors"
	"io"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/mikefarah/yq/v4/pkg/yqlib/treeops"
	yaml "gopkg.in/yaml.v3"
)

func readStream(filename string) (*yaml.Decoder, error) {
	if filename == "" {
		return nil, errors.New("Must provide filename")
	}

	var stream io.Reader
	if filename == "-" {
		stream = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(filename) // nolint gosec
		if err != nil {
			return nil, err
		}
		defer safelyCloseFile(file)
		stream = file
	}
	return yaml.NewDecoder(stream), nil
}

func evaluate(filename string, node *treeops.PathTreeNode) (*list.List, error) {

	var treeNavigator = treeops.NewDataTreeNavigator(treeops.NavigationPrefs{})

	var matchingNodes = list.New()

	var currentIndex uint = 0
	var decoder, err = readStream(filename)
	if err != nil {
		return nil, err
	}

	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			return matchingNodes, nil
		} else if errorReading != nil {
			return nil, errorReading
		}
		candidateNode := &treeops.CandidateNode{
			Document: currentIndex,
			Filename: filename,
			Node:     &dataBucket,
		}
		inputList := list.New()
		inputList.PushBack(candidateNode)

		newMatches, errorParsing := treeNavigator.GetMatchingNodes(inputList, node)
		if errorParsing != nil {
			return nil, errorParsing
		}
		matchingNodes.PushBackList(newMatches)
		currentIndex = currentIndex + 1
	}

	return matchingNodes, nil
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

func removeComments(matchingNodes *list.List) {
	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*treeops.CandidateNode)
		removeCommentOfNode(candidate.Node)
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

func setStyle(matchingNodes *list.List, style yaml.Style) {
	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*treeops.CandidateNode)
		updateStyleOfNode(candidate.Node, style)
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

func explode(matchingNodes *list.List) error {
	log.Debug("exploding nodes")
	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		nodeContext := el.Value.(*treeops.CandidateNode)
		log.Debugf("exploding %v", nodeContext.GetKey())
		errorExplodingNode := explodeNode(nodeContext.Node)
		if errorExplodingNode != nil {
			return errorExplodingNode
		}
	}
	return nil
}

func printResults(matchingNodes *list.List, writer io.Writer) error {
	if prettyPrint {
		setStyle(matchingNodes, 0)
	}

	if stripComments {
		removeComments(matchingNodes)
	}

	fileInfo, _ := os.Stdout.Stat()

	if forceColor || (!forceNoColor && (fileInfo.Mode()&os.ModeCharDevice) != 0) {
		colorsEnabled = true
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

	if matchingNodes.Len() == 0 {
		log.Debug("no matching results, nothing to print")
		if defaultValue != "" {
			return writeString(bufferedWriter, defaultValue)
		}
		return nil
	}
	var errorWriting error

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		mappedDoc := el.Value.(*treeops.CandidateNode)

		switch printMode {
		case "p":
			errorWriting = writeString(bufferedWriter, mappedDoc.PathStackToString()+"\n")
			if errorWriting != nil {
				return errorWriting
			}
		case "pv", "vp":
			// put it into a node and print that.
			var parentNode = yaml.Node{Kind: yaml.MappingNode}
			parentNode.Content = make([]*yaml.Node, 2)
			parentNode.Content[0] = &yaml.Node{Kind: yaml.ScalarNode, Value: mappedDoc.PathStackToString()}
			if mappedDoc.Node.Kind == yaml.DocumentNode {
				parentNode.Content[1] = mappedDoc.Node.Content[0]
			} else {
				parentNode.Content[1] = mappedDoc.Node
			}
			if err := printNode(&parentNode, bufferedWriter); err != nil {
				return err
			}
		default:
			if err := printNode(mappedDoc.Node, bufferedWriter); err != nil {
				return err
			}
		}
	}

	return nil
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
