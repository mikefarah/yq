//go:build !yq_nodotenv

package yqlib

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type dotenvDecoder struct {
	reader   io.Reader
	finished bool
}

func NewDotenvDecoder() Decoder {
	return &dotenvDecoder{finished: false}
}

func (dec *dotenvDecoder) Init(reader io.Reader) error {
	dec.reader = reader
	dec.finished = false
	return nil
}

// Decode reads a dotenv file into a flat map. The grammar is the one Node's
// --env-file and npm dotenv implement: `KEY=value`, an optional `export` prefix,
// `#` comments on their own line or starting a value, single quotes taken
// literally, double quotes where `\n` is the only escape sequence, and quoted
// values that may span several lines.
func (dec *dotenvDecoder) Decode() (*CandidateNode, error) {
	if dec.finished {
		return nil, io.EOF
	}

	root := &CandidateNode{Kind: MappingNode, Tag: "!!map"}

	scanner := bufio.NewScanner(dec.reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := trimByteOrderMark(scanner.Text())

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = trimExportPrefix(line)

		key, rawValue, found := strings.Cut(line, "=")
		if !found {
			return nil, fmt.Errorf("line %d: expected KEY=value, got %q", lineNumber, line)
		}

		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNumber)
		}

		// A quoted value keeps reading until its closing quote, so that a
		// multi-line value (a PEM key, say) is one value and not a parse error.
		for quoteNeedsMoreLines(rawValue) && scanner.Scan() {
			lineNumber++
			rawValue = rawValue + "\n" + trimByteOrderMark(scanner.Text())
		}

		value, err := parseDotenvValue(rawValue)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		root.AddKeyValueChild(createStringScalarNode(key), createStringScalarNode(value))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read dotenv content: %w", err)
	}

	dec.finished = true
	return root, nil
}

// trimByteOrderMark drops the BOM a Windows editor writes at the start of the
// file: it is not whitespace, so it would otherwise become part of the first key.
func trimByteOrderMark(line string) string {
	return strings.TrimPrefix(line, "\uFEFF")
}

func trimExportPrefix(line string) string {
	if !strings.HasPrefix(line, "export") {
		return line
	}
	rest := line[len("export"):]
	if rest == "" || (rest[0] != ' ' && rest[0] != '\t') {
		return line
	}
	return strings.TrimSpace(rest)
}

func quoteNeedsMoreLines(raw string) bool {
	value := strings.TrimSpace(raw)
	if len(value) < 2 {
		return false
	}
	switch value[0] {
	case '"':
		return !hasClosingQuote(value, '"')
	case '\'':
		return !hasClosingQuote(value, '\'')
	default:
		return false
	}
}

func hasClosingQuote(value string, quote byte) bool {
	for i := 1; i < len(value); i++ {
		// Only `\n` is an escape, so a backslash does not shield a quote.
		if quote == '"' && value[i] == '\\' && i+1 < len(value) && value[i+1] == 'n' {
			i++
			continue
		}
		if value[i] == quote {
			return true
		}
	}
	return false
}

func parseDotenvValue(raw string) (string, error) {
	value := strings.TrimSpace(raw)

	switch {
	case strings.HasPrefix(value, `"`):
		return parseQuotedDotenvValue(value, '"')
	case strings.HasPrefix(value, "'"):
		return parseQuotedDotenvValue(value, '\'')
	default:
		return stripUnquotedDotenvComment(value), nil
	}
}

// parseQuotedDotenvValue returns the contents of a quoted value. In a
// double-quoted value `\n` is a newline and any other backslash is kept, which
// is what Node's parser does; a single-quoted value is taken literally.
func parseQuotedDotenvValue(value string, quote byte) (string, error) {
	var result strings.Builder

	for i := 1; i < len(value); i++ {
		char := value[i]

		if char == quote {
			return result.String(), nil
		}

		// Only `\n` is an escape: every other backslash is kept as it is, so a
		// `\"` still ends the value, exactly as Node reads it.
		if quote == '"' && char == '\\' && i+1 < len(value) && value[i+1] == 'n' {
			result.WriteByte('\n')
			i++
			continue
		}

		result.WriteByte(char)
	}

	name := "double-quoted"
	if quote == '\'' {
		name = "single-quoted"
	}
	return "", fmt.Errorf("unterminated %s value %q", name, value)
}

// stripUnquotedDotenvComment removes a trailing comment. Node and npm dotenv
// both treat the first `#` as the start of one, so `a#b` is `a`, not `a#b`.
func stripUnquotedDotenvComment(value string) string {
	if index := strings.IndexByte(value, '#'); index >= 0 {
		return strings.TrimSpace(value[:index])
	}
	return strings.TrimSpace(value)
}
