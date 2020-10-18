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
	OperationType   *OperationType
	Value           interface{}
	StringValue     string
	PrefixSelf      bool

	CheckForPostTraverse bool // e.g. [1]cat should really be [1].cat
}

func pathToken(wrapped bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		value = value[1:len(value)]
		if wrapped {
			value = unwrap(value)
		}
		return &Token{PathElementType: PathKey, OperationType: None, Value: value, StringValue: value, CheckForPostTraverse: true}, nil
	}
}

func opToken(op *OperationType) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		return &Token{PathElementType: Operation, OperationType: op, Value: op.Type, StringValue: value}, nil
	}
}

func literalToken(pType PathElementType, literal string, checkForPost bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &Token{PathElementType: pType, OperationType: None, Value: literal, StringValue: literal, CheckForPostTraverse: checkForPost}, nil
	}
}

func unwrap(value string) string {
	return value[1 : len(value)-1]
}

func arrayIndextoken(precedingDot bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		startIndex := 1
		if precedingDot {
			startIndex = 2
		}
		numberString = numberString[startIndex : len(numberString)-1]
		var number, errParsingInt = strconv.ParseInt(numberString, 10, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}
		return &Token{PathElementType: PathKey, OperationType: None, Value: number, StringValue: numberString, CheckForPostTraverse: true}, nil
	}
}

func numberValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		var number, errParsingInt = strconv.ParseInt(numberString, 10, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}
		return &Token{PathElementType: Value, OperationType: None, Value: number, StringValue: numberString}, nil
	}
}

func floatValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		var number, errParsingInt = strconv.ParseFloat(numberString, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}
		return &Token{PathElementType: Value, OperationType: None, Value: number, StringValue: numberString}, nil
	}
}

func booleanValue(val bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &Token{PathElementType: Value, OperationType: None, Value: val, StringValue: string(m.Bytes)}, nil
	}
}

func stringValue(wrapped bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		if wrapped {
			value = unwrap(value)
		}
		return &Token{PathElementType: Value, OperationType: None, Value: value, StringValue: value}, nil
	}
}

func selfToken() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &Token{PathElementType: SelfReference, OperationType: None, Value: "SELF", StringValue: "SELF"}, nil
	}
}

// Creates the lexer object and compiles the NFA.
func initLexer() (*lex.Lexer, error) {
	lexer := lex.NewLexer()
	lexer.Add([]byte(`\(`), literalToken(OpenBracket, "(", false))
	lexer.Add([]byte(`\)`), literalToken(CloseBracket, ")", true))

	lexer.Add([]byte(`\.?\[\]`), literalToken(PathKey, "[]", true))
	lexer.Add([]byte(`\.\.`), opToken(RecursiveDescent))

	lexer.Add([]byte(`,`), opToken(Union))
	lexer.Add([]byte(`length`), opToken(Length))
	lexer.Add([]byte(`select`), opToken(Select))
	lexer.Add([]byte(`or`), opToken(Or))
	// lexer.Add([]byte(`and`), opToken())
	lexer.Add([]byte(`collect`), opToken(Collect))

	lexer.Add([]byte(`\s*==\s*`), opToken(Equals))

	lexer.Add([]byte(`\s*.-\s*`), opToken(DeleteChild))

	lexer.Add([]byte(`\s*\|=\s*`), opToken(Assign))

	lexer.Add([]byte(`\[-?[0-9]+\]`), arrayIndextoken(false))
	lexer.Add([]byte(`\.\[-?[0-9]+\]`), arrayIndextoken(true))

	lexer.Add([]byte("( |\t|\n|\r)+"), skip)

	lexer.Add([]byte(`\."[^ "]+"`), pathToken(true))
	lexer.Add([]byte(`\.[^ \[\],\|\.\[\(\)=]+`), pathToken(false))
	lexer.Add([]byte(`\.`), selfToken())

	lexer.Add([]byte(`\|`), opToken(Pipe))

	lexer.Add([]byte(`-?\d+(\.\d+)`), floatValue())
	lexer.Add([]byte(`-?[1-9](\.\d+)?[Ee][-+]?\d+`), floatValue())
	lexer.Add([]byte(`-?\d+`), numberValue())

	lexer.Add([]byte(`[Tt][Rr][Uu][Ee]`), booleanValue(true))
	lexer.Add([]byte(`[Ff][Aa][Ll][Ss][Ee]`), booleanValue(false))

	lexer.Add([]byte(`"[^ "]+"`), stringValue(true))

	lexer.Add([]byte(`\[`), literalToken(OpenCollect, "[", false))
	lexer.Add([]byte(`\]`), literalToken(CloseCollect, "]", true))
	lexer.Add([]byte(`\*`), opToken(Multiply))

	// lexer.Add([]byte(`[^ \,\|\.\[\(\)=]+`), stringValue(false))
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
			log.Debugf("Tokenising %v - %v", token.Value, token.OperationType.Type)
			tokens = append(tokens, token)
		}
		if err != nil {
			return nil, err
		}
	}
	var postProcessedTokens = make([]*Token, 0)

	for index, token := range tokens {

		postProcessedTokens = append(postProcessedTokens, token)

		if index != len(tokens)-1 && token.CheckForPostTraverse &&
			tokens[index+1].PathElementType == PathKey {
			postProcessedTokens = append(postProcessedTokens, &Token{PathElementType: Operation, OperationType: Pipe, Value: "PIPE"})
		}
	}

	return postProcessedTokens, nil
}
