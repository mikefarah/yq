package yqlib

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

type shellVariablesEncoder struct {
}

func NewShellVariablesEncoder() Encoder {
	return &shellVariablesEncoder{}
}

func (pe *shellVariablesEncoder) CanHandleAliases() bool {
	return false
}

func (pe *shellVariablesEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (pe *shellVariablesEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (pe *shellVariablesEncoder) Encode(writer io.Writer, node *CandidateNode) error {

	mapKeysToStrings(node)
	err := pe.doEncode(&writer, node, "")
	if err != nil {
		return err
	}

	return err
}

func (pe *shellVariablesEncoder) doEncode(w *io.Writer, node *CandidateNode, path string) error {

	// Note this drops all comments.

	switch node.Kind {
	case ScalarNode:
		nonemptyPath := path
		if path == "" {
			// We can't assign an empty variable "=somevalue" because that would error out if sourced in a shell,
			// nor can we use "_" as a variable name ($_ is a special shell variable that can't be assigned)...
			// let's just pick a fallback key to use if we are encoding a single scalar
			nonemptyPath = "value"
		}
		_, err := io.WriteString(*w, nonemptyPath+"="+quoteValue(node.Value)+"\n")
		return err
	case SequenceNode:
		for index, child := range node.Content {
			err := pe.doEncode(w, child, appendPath(path, index))
			if err != nil {
				return err
			}
		}
		return nil
	case MappingNode:
		for index := 0; index < len(node.Content); index = index + 2 {
			key := node.Content[index]
			value := node.Content[index+1]
			err := pe.doEncode(w, value, appendPath(path, key.Value))
			if err != nil {
				return err
			}
		}
		return nil
	case AliasNode:
		return pe.doEncode(w, node.Alias, path)
	default:
		return fmt.Errorf("unsupported node %v", node.Tag)
	}
}

func appendPath(cookedPath string, rawKey interface{}) string {

	// Shell variable names must match
	//    [a-zA-Z_]+[a-zA-Z0-9_]*
	//
	// While this is not mandated by POSIX, which is quite lenient, it is
	// what shells (for example busybox ash *) allow in practice.
	//
	// Since yaml names can contain basically any character, we will process them according to these steps:
	//
	//     1. apply unicode compatibility decomposition NFKD (this will convert accented
	//        letters to letters followed by accents, split ligatures, replace exponents
	//        with the corresponding digit, etc.
	//
	//     2. discard non-ASCII characters as well as ASCII control characters (ie. anything
	//        with code point < 32 or > 126), this will eg. discard accents but keep the base
	//        unaccented letter because of NFKD above
	//
	//     3. replace all non-alphanumeric characters with _
	//
	// Moreover, for the root key only, we will prepend an underscore if what results from the steps above
	// does not start with [a-zA-Z_] (ie. if the root key starts with a digit).
	//
	// Note this is NOT a 1:1 mapping.
	//
	// (*) see endofname.c from https://git.busybox.net/busybox/tag/?h=1_36_0

	// XXX empty strings

	key := strings.Map(func(r rune) rune {
		if isAlphaNumericOrUnderscore(r) {
			return r
		} else if r < 32 || 126 < r {
			return -1
		}
		return '_'
	}, norm.NFKD.String(fmt.Sprintf("%v", rawKey)))

	if cookedPath == "" {
		firstRune, _ := utf8.DecodeRuneInString(key)
		if !isAlphaOrUnderscore(firstRune) {
			return "_" + key
		}
		return key
	}
	return cookedPath + "_" + key
}

func quoteValue(value string) string {
	needsQuoting := false
	for _, r := range value {
		if !isAlphaNumericOrUnderscore(r) {
			needsQuoting = true
			break
		}
	}
	if needsQuoting {
		return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
	}
	return value
}

func isAlphaOrUnderscore(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || r == '_'
}

func isAlphaNumericOrUnderscore(r rune) bool {
	return isAlphaOrUnderscore(r) || ('0' <= r && r <= '9')
}
