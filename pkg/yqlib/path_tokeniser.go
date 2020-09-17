package yqlib

import (
	"strings"

	lex "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

var Literals []string       // The tokens representing literal strings
var Keywords []string       // The keyword tokens
var Tokens []string         // All of the tokens (including literals and keywords)
var TokenIds map[string]int // A map from the token names to their int ids

func initTokens() {
	Literals = []string{
		"(",
		")",
		"[+]",
		"[*]",
		"**",
	}
	Tokens = []string{
		"OPERATION",   // ==, OR, AND
		"PATH",        // a.b.c
		"ARRAY_INDEX", // 1234
		"PATH_JOIN",   // "."
	}
	Tokens = append(Tokens, Literals...)
	TokenIds = make(map[string]int)
	for i, tok := range Tokens {
		TokenIds[tok] = i
	}
}

func skip(*lex.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

func token(name string) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(TokenIds[name], string(m.Bytes), m), nil
	}
}

// Creates the lexer object and compiles the NFA.
func initLexer() (*lex.Lexer, error) {
	lexer := lex.NewLexer()
	for _, lit := range Literals {
		r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
		lexer.Add([]byte(r), token(lit))
	}
	lexer.Add([]byte(`([Oo][Rr]|[Aa][Nn][Dd]|==)`), token("OPERATION"))
	lexer.Add([]byte(`\[-?[0-9]+\]`), token("ARRAY_INDEX"))
	lexer.Add([]byte("( |\t|\n|\r)+"), skip)
	lexer.Add([]byte(`"[^ "]+"`), token("PATH"))
	lexer.Add([]byte(`[^ \.\[\(\)=]+`), token("PATH"))
	lexer.Add([]byte(`\.`), skip)
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

	return tokens, nil
}
