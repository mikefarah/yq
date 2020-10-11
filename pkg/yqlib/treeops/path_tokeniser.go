package treeops

import (
	"strconv"

	lex "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func skip(*lex.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

type Token struct {
	PathElementType PathElementType
	OperationType   OperationType
	Value           interface{}
	StringValue     string
	AgainstSelf     bool

	CheckForPreTraverse bool // this token can sometimes have the traverse '.' missing in frnot of it
	// e.g. a[1] should really be a.[1]
	CheckForPostTraverse bool // samething but for post, e.g. [1]cat should really be [1].cat

}

func pathToken(wrapped bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		if wrapped {
			value = unwrap(value)
		}
		return &Token{PathElementType: PathKey, OperationType: None, Value: value, StringValue: value}, nil
	}
}

func opToken(op OperationType, againstSelf bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		return &Token{PathElementType: Operation, OperationType: op, Value: value, StringValue: value, AgainstSelf: againstSelf}, nil
	}
}

func literalToken(pType PathElementType, literal string, checkForPre bool, checkForPost bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &Token{PathElementType: pType, Value: literal, StringValue: literal, CheckForPreTraverse: checkForPre, CheckForPostTraverse: checkForPost}, nil
	}
}

func unwrap(value string) string {
	return value[1 : len(value)-1]
}

func arrayIndextoken(wrapped bool, checkForPre bool, checkForPost bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		if wrapped {
			numberString = unwrap(numberString)
		}
		var number, errParsingInt = strconv.ParseInt(numberString, 10, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}
		return &Token{PathElementType: ArrayIndex, Value: number, StringValue: numberString, CheckForPreTraverse: checkForPre, CheckForPostTraverse: checkForPost}, nil
	}
}

// Creates the lexer object and compiles the NFA.
func initLexer() (*lex.Lexer, error) {
	lexer := lex.NewLexer()
	lexer.Add([]byte(`\(`), literalToken(OpenBracket, "(", true, false))
	lexer.Add([]byte(`\)`), literalToken(CloseBracket, ")", false, true))

	lexer.Add([]byte(`\[\+\]`), literalToken(PathKey, "[+]", true, true))
	lexer.Add([]byte(`\[\*\]`), literalToken(PathKey, "[*]", true, true))
	lexer.Add([]byte(`\*\*`), literalToken(PathKey, "**", false, false))

	lexer.Add([]byte(`([Oo][Rr])`), opToken(Or, false))
	lexer.Add([]byte(`([Aa][Nn][Dd])`), opToken(And, false))

	lexer.Add([]byte(`\.\s*==\s*`), opToken(Equals, true))
	lexer.Add([]byte(`\s*==\s*`), opToken(Equals, false))

	lexer.Add([]byte(`\.\s*.-\s*`), opToken(DeleteChild, true))
	lexer.Add([]byte(`\s*.-\s*`), opToken(DeleteChild, false))

	lexer.Add([]byte(`\.\s*:=\s*`), opToken(Assign, true))
	lexer.Add([]byte(`\s*:=\s*`), opToken(Assign, false))

	lexer.Add([]byte(`\[-?[0-9]+\]`), arrayIndextoken(true, true, true))
	lexer.Add([]byte(`-?[0-9]+`), arrayIndextoken(false, false, false))
	lexer.Add([]byte("( |\t|\n|\r)+"), skip)

	lexer.Add([]byte(`"[^ "]+"`), pathToken(true))
	lexer.Add([]byte(`[^ \.\[\(\)=]+`), pathToken(false))

	lexer.Add([]byte(`\.`), opToken(Traverse, false))
	err := lexer.Compile()
	if err != nil {
		return nil, err
	}
	return lexer, nil
}

type PathTokeniser interface {
	Tokenise(path string) ([]*Token, error)
}

type pathTokeniser struct {
	lexer *lex.Lexer
}

func NewPathTokeniser() PathTokeniser {
	var lexer, err = initLexer()
	if err != nil {
		panic(err)
	}
	return &pathTokeniser{lexer}
}

func (p *pathTokeniser) Tokenise(path string) ([]*Token, error) {
	scanner, err := p.lexer.Scanner([]byte(path))

	if err != nil {
		return nil, err
	}
	var tokens []*Token
	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {

		if tok != nil {
			token := tok.(*Token)
			log.Debugf("Tokenising %v", token.Value)
			tokens = append(tokens, token)
		}
		if err != nil {
			return nil, err
		}
	}
	var postProcessedTokens = make([]*Token, 0)

	for index, token := range tokens {
		if index > 0 && token.CheckForPreTraverse &&
			(tokens[index-1].PathElementType == PathKey || tokens[index-1].PathElementType == CloseBracket) {
			postProcessedTokens = append(postProcessedTokens, &Token{PathElementType: Operation, OperationType: Traverse, Value: "."})
		}
		if token.PathElementType == Operation && token.AgainstSelf {
			postProcessedTokens = append(postProcessedTokens, &Token{PathElementType: SelfReference, Value: "SELF"})
		}

		postProcessedTokens = append(postProcessedTokens, token)

		if index != len(tokens)-1 && token.CheckForPostTraverse &&
			tokens[index+1].PathElementType == PathKey {
			postProcessedTokens = append(postProcessedTokens, &Token{PathElementType: Operation, OperationType: Traverse, Value: "."})
		}
	}

	return postProcessedTokens, nil
}
