package yqlib

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/magiconair/properties"
)

type propertiesEncoder struct {
	prefs PropertiesPreferences
}

func NewPropertiesEncoder(prefs PropertiesPreferences) Encoder {
	return &propertiesEncoder{
		prefs: prefs,
	}
}

func (pe *propertiesEncoder) CanHandleAliases() bool {
	return false
}

func (pe *propertiesEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (pe *propertiesEncoder) PrintLeadingContent(writer io.Writer, content string) error {
	reader := bufio.NewReader(strings.NewReader(content))
	for {

		readline, errReading := reader.ReadString('\n')
		if errReading != nil && !errors.Is(errReading, io.EOF) {
			return errReading
		}
		if strings.Contains(readline, "$yqDocSeparator$") {

			if err := pe.PrintDocumentSeparator(writer); err != nil {
				return err
			}

		} else {
			if err := writeString(writer, readline); err != nil {
				return err
			}
		}

		if errors.Is(errReading, io.EOF) {
			if readline != "" {
				// the last comment we read didn't have a newline, put one in
				if err := writeString(writer, "\n"); err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}

func (pe *propertiesEncoder) Encode(writer io.Writer, node *CandidateNode) error {

	if node.Kind == ScalarNode {
		return writeString(writer, node.Value+"\n")
	}

	mapKeysToStrings(node)
	p := properties.NewProperties()
	p.WriteSeparator = pe.prefs.KeyValueSeparator
	err := pe.doEncode(p, node, "", nil)
	if err != nil {
		return err
	}

	_, err = p.WriteComment(writer, "#", properties.UTF8)
	return err
}

func (pe *propertiesEncoder) doEncode(p *properties.Properties, node *CandidateNode, path string, keyNode *CandidateNode) error {

	comments := ""
	if keyNode != nil {
		// include the key node comments if present
		comments = headAndLineComment(keyNode)
	}
	comments = comments + headAndLineComment(node)
	commentsWithSpaces := strings.ReplaceAll(comments, "\n", "\n ")
	p.SetComments(path, strings.Split(commentsWithSpaces, "\n"))

	switch node.Kind {
	case ScalarNode:
		var nodeValue string
		if pe.prefs.UnwrapScalar || !strings.Contains(node.Value, " ") {
			nodeValue = node.Value
		} else {
			nodeValue = fmt.Sprintf("%q", node.Value)
		}
		_, _, err := p.Set(path, nodeValue)
		return err
	case SequenceNode:
		return pe.encodeArray(p, node.Content, path)
	case MappingNode:
		return pe.encodeMap(p, node.Content, path)
	case AliasNode:
		return pe.doEncode(p, node.Alias, path, nil)
	default:
		return fmt.Errorf("unsupported node %v", node.Tag)
	}
}

func (pe *propertiesEncoder) appendPath(path string, key interface{}) string {
	if path == "" {
		return fmt.Sprintf("%v", key)
	}
	switch key.(type) {
	case int:
		if pe.prefs.UseArrayBrackets {
			return fmt.Sprintf("%v[%v]", path, key)
		}

	}
	return fmt.Sprintf("%v.%v", path, key)
}

func (pe *propertiesEncoder) encodeArray(p *properties.Properties, kids []*CandidateNode, path string) error {
	for index, child := range kids {
		err := pe.doEncode(p, child, pe.appendPath(path, index), nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pe *propertiesEncoder) encodeMap(p *properties.Properties, kids []*CandidateNode, path string) error {
	for index := 0; index < len(kids); index = index + 2 {
		key := kids[index]
		value := kids[index+1]
		err := pe.doEncode(p, value, pe.appendPath(path, key.Value), key)
		if err != nil {
			return err
		}
	}
	return nil
}
