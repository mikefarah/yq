//go:build !yq_nodotenv

package yqlib

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

type dotenvEncoder struct{}

func NewDotenvEncoder() Encoder {
	return &dotenvEncoder{}
}

func (de *dotenvEncoder) CanHandleAliases() bool {
	return false
}

func (de *dotenvEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (de *dotenvEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (de *dotenvEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	mapKeysToStrings(node)
	return de.doEncode(&writer, node, "")
}

// doEncode flattens the tree into `KEY=value` lines: nested maps are joined with
// an underscore and sequences are indexed, as shell variables cannot be nested
// either. Comments are dropped, as they cannot be attached to a `KEY=value` line
// without changing its value.
func (de *dotenvEncoder) doEncode(w *io.Writer, node *CandidateNode, path string) error {
	switch node.Kind {
	case ScalarNode:
		nonemptyPath := path
		if nonemptyPath == "" {
			// A bare scalar has no key to be assigned to; `.` of a one-line file
			// is the common case, so give it a name instead of failing.
			nonemptyPath = "value"
		}
		quoted, err := quoteDotenvValue(node.Value)
		if err != nil {
			return fmt.Errorf("cannot encode %v: %w", nonemptyPath, err)
		}
		_, err = io.WriteString(*w, nonemptyPath+"="+quoted+"\n")
		return err
	case SequenceNode:
		for index, child := range node.Content {
			if err := de.doEncode(w, child, dotenvAppendPath(path, index)); err != nil {
				return err
			}
		}
		return nil
	case MappingNode:
		for index := 0; index < len(node.Content); index = index + 2 {
			key := node.Content[index]
			value := node.Content[index+1]
			if err := de.doEncode(w, value, dotenvAppendPath(path, key.Value)); err != nil {
				return err
			}
		}
		return nil
	case AliasNode:
		return de.doEncode(w, node.Alias, path)
	default:
		return fmt.Errorf("unsupported node %v", node.Tag)
	}
}

// dotenvAppendPath builds an environment variable name. Names have to match
// [a-zA-Z_][a-zA-Z0-9_]*, so a key is decomposed with NFKD, stripped of anything
// outside printable ASCII and with its remaining punctuation turned into an
// underscore. This is not a 1:1 mapping, and it matches what the shell variables
// output does.
func dotenvAppendPath(cookedPath string, rawKey interface{}) string {
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

// quoteDotenvValue picks the quoting dotenv readers agree on: bare when the
// value survives unquoted, single quotes otherwise, because they are literal in
// Node's parser, npm dotenv and in this decoder. Double quotes are only used for
// a value that spans lines, where `\n` is the one escape every reader undoes.
func quoteDotenvValue(value string) (string, error) {
	if value == "" {
		return `""`, nil
	}
	if !dotenvValueNeedsQuoting(value) {
		return value, nil
	}

	if strings.ContainsAny(value, "\n\r") {
		if strings.Contains(value, `"`) {
			return "", fmt.Errorf(`value %q has both a newline and a double quote, which dotenv cannot express`, value)
		}
		return `"` + strings.ReplaceAll(strings.ReplaceAll(value, "\r\n", `\n`), "\n", `\n`) + `"`, nil
	}

	if strings.Contains(value, "'") {
		if strings.Contains(value, `"`) {
			return "", fmt.Errorf(`value %q has both quote characters, which dotenv cannot express`, value)
		}
		return `"` + value + `"`, nil
	}

	return "'" + value + "'", nil
}

func dotenvValueNeedsQuoting(value string) bool {
	for _, char := range value {
		switch char {
		case '\'', '"', '$', '\\', '#', ' ', '\t', '\n', '\r':
			return true
		}
	}
	return false
}
