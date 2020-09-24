package yqlib

import (
	"strconv"
	"strings"

	lex "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

var Literals []string        // The tokens representing literal strings
var ClosingLiterals []string // The tokens representing literal strings
var Keywords []string        // The keyword tokens
var Tokens []string          // All of the tokens (including literals and keywords)
var TokenIds map[string]int  // A map from the token names to their int ids

func initTokens() {
	Literals = []string{ // these need a traverse operator infront
		"(",
		"[+]",
		"[*]",
		"**",
	}
	ClosingLiterals = []string{ // these need a traverse operator after
		")",
	}
	Tokens = []string{
		"BEGIN_SUB_EXPRESSION",
		"END_SUB_EXPRESSION",
		"OR_OPERATOR",
		"AND_OPERATOR",
		"EQUALS_OPERATOR",
		"EQUALS_SELF_OPERATOR",
		"TRAVERSE_OPERATOR",
		"PATH_KEY",    // apples
		"ARRAY_INDEX", // 123
	}
	Tokens = append(Tokens, Literals...)
	Tokens = append(Tokens, ClosingLiterals...)
	TokenIds = make(map[string]int)
	for i, tok := range Tokens {
		TokenIds[tok] = i
	}

	initMaps()
}

func skip(*lex.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

func token(name string) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(TokenIds[name], string(m.Bytes), m), nil
	}
}

func unwrap(value string) string {
	return value[1 : len(value)-1]
}

func wrappedToken(name string) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(TokenIds[name], unwrap(string(m.Bytes)), m), nil
	}
}

func numberToken(name string, wrapped bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		if wrapped {
			numberString = unwrap(numberString)
		}
		var number, errParsingInt = strconv.ParseInt(numberString, 10, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}
		return s.Token(TokenIds[name], number, m), nil
	}
}

// Creates the lexer object and compiles the NFA.
func initLexer() (*lex.Lexer, error) {
	lexer := lex.NewLexer()
	for _, lit := range Literals {
		r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
		lexer.Add([]byte(r), token(lit))
	}
	for _, lit := range ClosingLiterals {
		r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
		lexer.Add([]byte(r), token(lit))
	}
	lexer.Add([]byte(`([Oo][Rr])`), token("OR_OPERATOR"))
	lexer.Add([]byte(`([Aa][Nn][Dd])`), token("AND_OPERATOR"))
	lexer.Add([]byte(`\.\s*==\s*`), token("EQUALS_SELF_OPERATOR"))
	lexer.Add([]byte(`\s*==\s*`), token("EQUALS_OPERATOR"))
	lexer.Add([]byte(`\[-?[0-9]+\]`), numberToken("ARRAY_INDEX", true))
	lexer.Add([]byte(`-?[0-9]+`), numberToken("ARRAY_INDEX", false))
	lexer.Add([]byte("( |\t|\n|\r)+"), skip)
	lexer.Add([]byte(`"[^ "]+"`), wrappedToken("PATH_KEY"))
	lexer.Add([]byte(`[^ \.\[\(\)=]+`), token("PATH_KEY"))
	lexer.Add([]byte(`\.`), token("TRAVERSE_OPERATOR"))
	err := lexer.Compile()
	if err != nil {
		return nil, err
	}
	return lexer, nil
}

type PathTokeniser interface {
	Tokenise(path string) ([]*lex.Token, error)
}

type pathTokeniser struct {
	lexer *lex.Lexer
}

func NewPathTokeniser() PathTokeniser {
	initTokens()
	var lexer, err = initLexer()
	if err != nil {
		panic(err)
	}
	return &pathTokeniser{lexer}
}

func (p *pathTokeniser) Tokenise(path string) ([]*lex.Token, error) {
	scanner, err := p.lexer.Scanner([]byte(path))

	if err != nil {
		return nil, err
	}
	var tokens []*lex.Token
	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {

		if tok != nil {
			token := tok.(*lex.Token)
			log.Debugf("Processing %v - %v", token.Value, Tokens[token.Type])
			tokens = append(tokens, token)
		}
		if err != nil {
			return nil, err
		}
	}
	var postProcessedTokens []*lex.Token = make([]*lex.Token, 0)

	for index, token := range tokens {
		for _, literalTokenDef := range append(Literals, "ARRAY_INDEX") {
			if index > 0 && token.Type == TokenIds[literalTokenDef] && tokens[index-1].Type != TokenIds["TRAVERSE_OPERATOR"] {
				postProcessedTokens = append(postProcessedTokens, &lex.Token{Type: TokenIds["TRAVERSE_OPERATOR"], Value: "."})
			}
		}

		postProcessedTokens = append(postProcessedTokens, token)
		for _, literalTokenDef := range append(ClosingLiterals, "ARRAY_INDEX") {
			if index != len(tokens)-1 && token.Type == TokenIds[literalTokenDef] && tokens[index+1].Type != TokenIds["TRAVERSE_OPERATOR"] {
				postProcessedTokens = append(postProcessedTokens, &lex.Token{Type: TokenIds["TRAVERSE_OPERATOR"], Value: "."})
			}
		}
	}

	return postProcessedTokens, nil
}
