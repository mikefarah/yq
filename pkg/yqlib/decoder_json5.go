//go:build !yq_nojson5

package yqlib

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
)

type json5Decoder struct {
	parser *json5Parser
}

func NewJSON5Decoder() Decoder {
	return &json5Decoder{}
}

func (dec *json5Decoder) Init(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	dec.parser = newJSON5Parser(string(data))
	return nil
}

func (dec *json5Decoder) Decode() (*CandidateNode, error) {
	if dec.parser == nil {
		return nil, io.EOF
	}
	if err := dec.parser.skipWhitespaceAndComments(); err != nil {
		return nil, err
	}
	if dec.parser.eof() {
		return nil, io.EOF
	}
	return dec.parser.parseValue()
}

type json5Parser struct {
	input           []rune
	pos             int
	line            int
	col             int
	pendingComments []json5Comment
}

func newJSON5Parser(s string) *json5Parser {
	return &json5Parser{
		input: []rune(s),
		pos:   0,
		line:  1,
		col:   1,
	}
}

type json5Comment struct {
	text            string
	startsOnNewLine bool
}

func (p *json5Parser) eof() bool {
	return p.pos >= len(p.input)
}

func (p *json5Parser) peek() rune {
	if p.eof() {
		return 0
	}
	return p.input[p.pos]
}

func (p *json5Parser) peekNext() rune {
	i := p.pos + 1
	if i >= len(p.input) {
		return 0
	}
	return p.input[i]
}

func (p *json5Parser) next() rune {
	if p.eof() {
		return 0
	}
	r := p.input[p.pos]
	p.pos++
	if r == '\n' {
		p.line++
		p.col = 1
	} else {
		p.col++
	}
	return r
}

func (p *json5Parser) errorf(format string, args ...interface{}) error {
	return fmt.Errorf("json5: %s at line %d, column %d", fmt.Sprintf(format, args...), p.line, p.col)
}

func (p *json5Parser) skipWhitespaceAndComments() error {
	sawNewline := false

	for !p.eof() {
		r := p.peek()
		if unicode.IsSpace(r) {
			if r == '\n' || r == '\r' || r == '\u2028' || r == '\u2029' {
				sawNewline = true
			}
			p.next()
			continue
		}
		if r == '/' && p.peekNext() == '/' {
			startsOnNewLine := sawNewline
			p.next() // /
			p.next() // /
			var sb strings.Builder
			for !p.eof() && p.peek() != '\n' {
				sb.WriteRune(p.next())
			}
			p.pendingComments = append(p.pendingComments, json5Comment{
				text:            strings.TrimSpace(sb.String()),
				startsOnNewLine: startsOnNewLine,
			})
			continue
		}
		if r == '/' && p.peekNext() == '*' {
			startsOnNewLine := sawNewline
			p.next() // /
			p.next() // *
			var sb strings.Builder
			for {
				if p.eof() {
					return p.errorf("unterminated block comment")
				}
				if p.peek() == '*' && p.peekNext() == '/' {
					p.next()
					p.next()
					break
				}
				sb.WriteRune(p.next())
			}
			normalised := normaliseJSON5BlockComment(sb.String())
			if strings.Contains(normalised, "\n") {
				sawNewline = true
			}
			p.pendingComments = append(p.pendingComments, json5Comment{
				text:            normalised,
				startsOnNewLine: startsOnNewLine,
			})
			continue
		}
		break
	}
	return nil
}

func (p *json5Parser) takePendingComments() []json5Comment {
	if len(p.pendingComments) == 0 {
		return nil
	}
	comments := p.pendingComments
	p.pendingComments = nil
	return comments
}

func commentsToText(comments []json5Comment) string {
	if len(comments) == 0 {
		return ""
	}
	parts := make([]string, 0, len(comments))
	for _, c := range comments {
		if strings.TrimSpace(c.text) == "" {
			continue
		}
		parts = append(parts, c.text)
	}
	return strings.Join(parts, "\n")
}

func normaliseJSON5BlockComment(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.Trim(content, "\n\r\t ")
	if content == "" {
		return ""
	}

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimLeft(line, " \t")
		if strings.HasPrefix(line, "*") {
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimLeft(line, " \t")
		}
		lines[i] = line
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func (p *json5Parser) parseValue() (*CandidateNode, error) {
	if err := p.skipWhitespaceAndComments(); err != nil {
		return nil, err
	}
	leading := commentsToText(p.takePendingComments())
	if p.eof() {
		return nil, io.EOF
	}

	var node *CandidateNode
	var err error

	switch r := p.peek(); r {
	case '{':
		node, err = p.parseObject()
	case '[':
		node, err = p.parseArray()
	case '"', '\'':
		s, err := p.parseString()
		if err != nil {
			return nil, err
		}
		node = createScalarNode(s, s)
	case '-', '+', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		node, err = p.parseNumber()
	default:
		if isIdentifierStart(r) {
			ident, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			switch ident {
			case "true":
				node = createScalarNode(true, "true")
			case "false":
				node = createScalarNode(false, "false")
			case "null":
				node = createScalarNode(nil, "null")
			case "Infinity":
				node = &CandidateNode{Kind: ScalarNode, Tag: "!!float", Value: "+Inf"}
			case "NaN":
				node = &CandidateNode{Kind: ScalarNode, Tag: "!!float", Value: "NaN"}
			default:
				return nil, p.errorf("unexpected identifier %q", ident)
			}
		} else {
			return nil, p.errorf("unexpected character %q", r)
		}
	}

	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, p.errorf("invalid value")
	}
	if leading != "" {
		if node.HeadComment != "" {
			node.HeadComment = leading + "\n" + node.HeadComment
		} else {
			node.HeadComment = leading
		}
	}
	return node, nil
}

func (p *json5Parser) parseObject() (*CandidateNode, error) {
	if p.next() != '{' {
		return nil, p.errorf("expected '{'")
	}

	node := &CandidateNode{Kind: MappingNode, Tag: "!!map"}
	if err := p.skipWhitespaceAndComments(); err != nil {
		return nil, err
	}
	node.HeadComment = commentsToText(p.takePendingComments())
	if p.peek() == '}' {
		p.next()
		return node, nil
	}

	for {
		if err := p.skipWhitespaceAndComments(); err != nil {
			return nil, err
		}
		pendingBeforeKey := commentsToText(p.takePendingComments())

		key, err := p.parseObjectKey()
		if err != nil {
			return nil, err
		}

		if err := p.skipWhitespaceAndComments(); err != nil {
			return nil, err
		}
		pendingAfterKey := commentsToText(p.takePendingComments())
		if p.next() != ':' {
			return nil, p.errorf("expected ':' after object key")
		}

		if err := p.skipWhitespaceAndComments(); err != nil {
			return nil, err
		}
		pendingAfterColon := commentsToText(p.takePendingComments())

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		childKey := node.CreateChild()
		childKey.IsMapKey = true
		childKey.Value = key
		childKey.Kind = ScalarNode
		childKey.Tag = "!!str"
		childKey.HeadComment = pendingBeforeKey
		childKey.LineComment = pendingAfterKey

		if pendingAfterColon != "" {
			if value.HeadComment != "" {
				value.HeadComment = pendingAfterColon + "\n" + value.HeadComment
			} else {
				value.HeadComment = pendingAfterColon
			}
		}

		value.Parent = node
		value.Key = childKey
		node.Content = append(node.Content, childKey, value)

		if err := p.skipWhitespaceAndComments(); err != nil {
			return nil, err
		}
		value.LineComment = commentsToText(p.takePendingComments())

		switch p.peek() {
		case ',':
			p.next()
			if err := p.skipWhitespaceAndComments(); err != nil {
				return nil, err
			}
			if p.peek() == '}' {
				p.next()
				return node, nil
			}
		case '}':
			p.next()
			return node, nil
		default:
			return nil, p.errorf("expected ',' or '}' after object entry")
		}
	}
}

func (p *json5Parser) parseObjectKey() (string, error) {
	if err := p.skipWhitespaceAndComments(); err != nil {
		return "", err
	}

	switch p.peek() {
	case '"', '\'':
		return p.parseString()
	default:
		r := p.peek()
		if !isIdentifierStart(r) && (r != '\\' || p.peekNext() != 'u') {
			return "", p.errorf("expected object key")
		}
		return p.parseIdentifier()
	}
}

func (p *json5Parser) parseArray() (*CandidateNode, error) {
	if p.next() != '[' {
		return nil, p.errorf("expected '['")
	}

	node := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}
	if err := p.skipWhitespaceAndComments(); err != nil {
		return nil, err
	}
	node.HeadComment = commentsToText(p.takePendingComments())
	if p.peek() == ']' {
		p.next()
		return node, nil
	}

	index := 0
	for {
		if err := p.skipWhitespaceAndComments(); err != nil {
			return nil, err
		}
		pendingBeforeElement := commentsToText(p.takePendingComments())

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		if pendingBeforeElement != "" {
			if value.HeadComment != "" {
				value.HeadComment = pendingBeforeElement + "\n" + value.HeadComment
			} else {
				value.HeadComment = pendingBeforeElement
			}
		}

		childKey := node.CreateChild()
		childKey.Kind = ScalarNode
		childKey.Tag = "!!int"
		childKey.Value = fmt.Sprintf("%v", index)
		childKey.IsMapKey = true

		value.Parent = node
		value.Key = childKey
		node.Content = append(node.Content, value)
		index++

		if err := p.skipWhitespaceAndComments(); err != nil {
			return nil, err
		}
		value.LineComment = commentsToText(p.takePendingComments())

		switch p.peek() {
		case ',':
			p.next()
			if err := p.skipWhitespaceAndComments(); err != nil {
				return nil, err
			}
			if p.peek() == ']' {
				p.next()
				return node, nil
			}
		case ']':
			p.next()
			return node, nil
		default:
			return nil, p.errorf("expected ',' or ']' after array element")
		}
	}
}

func (p *json5Parser) parseString() (string, error) {
	quote := p.next()
	var sb strings.Builder

	for {
		if p.eof() {
			return "", p.errorf("unterminated string")
		}

		r := p.next()
		if r == quote {
			return sb.String(), nil
		}

		if r == '\n' || r == '\r' || r == '\u2028' || r == '\u2029' {
			return "", p.errorf("unterminated string")
		}

		if r != '\\' {
			sb.WriteRune(r)
			continue
		}

		if p.eof() {
			return "", p.errorf("unterminated escape sequence")
		}

		esc := p.next()

		if esc == '\n' || esc == '\u2028' || esc == '\u2029' {
			continue
		}
		if esc == '\r' {
			if p.peek() == '\n' {
				p.next()
			}
			continue
		}

		switch esc {
		case '\\', '/', '"', '\'':
			sb.WriteRune(esc)
		case 'b':
			sb.WriteByte('\b')
		case 'f':
			sb.WriteByte('\f')
		case 'n':
			sb.WriteByte('\n')
		case 'r':
			sb.WriteByte('\r')
		case 't':
			sb.WriteByte('\t')
		case 'v':
			sb.WriteByte('\v')
		case '0':
			if isDigit(p.peek()) {
				return "", p.errorf("invalid escape sequence \\0 followed by digit")
			}
			sb.WriteByte(0)
		case 'x':
			r, err := p.parseHexEscape(2)
			if err != nil {
				return "", err
			}
			sb.WriteRune(r)
		case 'u':
			r, err := p.parseUnicodeEscape()
			if err != nil {
				return "", err
			}
			sb.WriteRune(r)
		default:
			sb.WriteRune(esc)
		}
	}
}

func (p *json5Parser) parseHexEscape(length int) (rune, error) {
	if p.pos+length > len(p.input) {
		return 0, p.errorf("invalid hex escape")
	}
	var value rune
	for i := 0; i < length; i++ {
		d := p.next()
		h, ok := hexDigitValue(d)
		if !ok {
			return 0, p.errorf("invalid hex escape")
		}
		value = (value << 4) | h
	}
	return value, nil
}

func (p *json5Parser) parseUnicodeEscape() (rune, error) {
	r, err := p.parseHexEscape(4)
	if err != nil {
		return 0, err
	}

	if utf16.IsSurrogate(r) {
		originalPos, originalLine, originalCol := p.pos, p.line, p.col
		if p.peek() == '\\' && p.peekNext() == 'u' {
			p.next()
			p.next()
			r2, err := p.parseHexEscape(4)
			if err == nil && utf16.IsSurrogate(r2) {
				return utf16.DecodeRune(r, r2), nil
			}
		}
		p.pos, p.line, p.col = originalPos, originalLine, originalCol
	}

	return r, nil
}

func (p *json5Parser) parseNumber() (*CandidateNode, error) {
	startPos := p.pos

	sign := rune(0)
	if p.peek() == '+' || p.peek() == '-' {
		sign = p.next()
	}

	if isIdentifierStart(p.peek()) {
		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if ident != "Infinity" {
			return nil, p.errorf("invalid number")
		}
		if sign == '-' {
			return &CandidateNode{Kind: ScalarNode, Tag: "!!float", Value: "-Inf"}, nil
		}
		return &CandidateNode{Kind: ScalarNode, Tag: "!!float", Value: "+Inf"}, nil
	}

	posAfterSign := p.pos

	isFloat := false
	if p.peek() == '0' && (p.peekNext() == 'x' || p.peekNext() == 'X') {
		p.next()
		p.next()
		if !isHexDigit(p.peek()) {
			return nil, p.errorf("invalid hex number")
		}
		for isHexDigit(p.peek()) {
			p.next()
		}
		lit := string(p.input[posAfterSign:p.pos])
		if sign == '-' {
			lit = "-" + lit
		}
		// if sign == '+' { } // drop explicit plus sign
		return &CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: lit}, nil
	}

	if p.peek() == '.' {
		isFloat = true
		p.next()
		if !isDigit(p.peek()) {
			return nil, p.errorf("invalid number")
		}
		for isDigit(p.peek()) {
			p.next()
		}
	} else {
		if !isDigit(p.peek()) {
			return nil, p.errorf("invalid number")
		}
		for isDigit(p.peek()) {
			p.next()
		}
		if p.peek() == '.' {
			isFloat = true
			p.next()
			for isDigit(p.peek()) {
				p.next()
			}
		}
	}

	if p.peek() == 'e' || p.peek() == 'E' {
		isFloat = true
		p.next()
		if p.peek() == '+' || p.peek() == '-' {
			p.next()
		}
		if !isDigit(p.peek()) {
			return nil, p.errorf("invalid number exponent")
		}
		for isDigit(p.peek()) {
			p.next()
		}
	}

	lit := string(p.input[startPos:p.pos])

	lit = strings.TrimPrefix(lit, "+")

	if isFloat {
		if _, err := strconv.ParseFloat(lit, 64); err != nil {
			return nil, p.errorf("invalid float number %q", lit)
		}
		return &CandidateNode{Kind: ScalarNode, Tag: "!!float", Value: lit}, nil
	}

	if _, err := strconv.ParseInt(lit, 10, 64); err != nil {
		return nil, p.errorf("invalid integer %q", lit)
	}
	return &CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: lit}, nil
}

func (p *json5Parser) parseIdentifier() (string, error) {
	var sb strings.Builder

	r := p.peek()
	if !isIdentifierStart(r) && (r != '\\' || p.peekNext() != 'u') {
		return "", p.errorf("expected identifier")
	}

	for !p.eof() {
		r := p.peek()
		if r == '\\' && p.peekNext() == 'u' {
			p.next()
			p.next()
			decoded, err := p.parseUnicodeEscape()
			if err != nil {
				return "", err
			}
			if sb.Len() == 0 {
				if !isIdentifierStart(decoded) {
					return "", p.errorf("invalid identifier start")
				}
			} else if !isIdentifierPart(decoded) {
				return "", p.errorf("invalid identifier part")
			}
			sb.WriteRune(decoded)
			continue
		}

		if sb.Len() == 0 {
			if !isIdentifierStart(r) {
				break
			}
		} else if !isIdentifierPart(r) {
			break
		}

		sb.WriteRune(p.next())
	}

	return sb.String(), nil
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func hexDigitValue(r rune) (rune, bool) {
	switch {
	case r >= '0' && r <= '9':
		return r - '0', true
	case r >= 'a' && r <= 'f':
		return r - 'a' + 10, true
	case r >= 'A' && r <= 'F':
		return r - 'A' + 10, true
	default:
		return 0, false
	}
}

func isIdentifierStart(r rune) bool {
	return r == '$' || r == '_' || unicode.IsLetter(r)
}

func isIdentifierPart(r rune) bool {
	return isIdentifierStart(r) || unicode.IsDigit(r)
}
