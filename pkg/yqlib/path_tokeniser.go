package yqlib

import (
	"fmt"
	"strconv"

	lex "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func skip(*lex.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

type TokenType uint32

const (
	OperationToken = 1 << iota
	OpenBracket
	CloseBracket
	OpenCollect
	CloseCollect
	OpenCollectObject
	CloseCollectObject
	SplatOrEmptyCollect
)

type Token struct {
	TokenType            TokenType
	Operation            *Operation
	AssignOperation      *Operation // e.g. tag (GetTag) op becomes AssignTag if '=' follows it
	CheckForPostTraverse bool       // e.g. [1]cat should really be [1].cat

}

func (t *Token) toString() string {
	if t.TokenType == OperationToken {
		log.Debug("toString, its an op")
		return t.Operation.toString()
	} else if t.TokenType == OpenBracket {
		return "("
	} else if t.TokenType == CloseBracket {
		return ")"
	} else if t.TokenType == OpenCollect {
		return "["
	} else if t.TokenType == CloseCollect {
		return "]"
	} else if t.TokenType == OpenCollectObject {
		return "{"
	} else if t.TokenType == CloseCollectObject {
		return "}"
	} else if t.TokenType == SplatOrEmptyCollect {
		return "[]?"
	} else {
		return "NFI"
	}
}

func pathToken(wrapped bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		value = value[1:]
		if wrapped {
			value = unwrap(value)
		}
		log.Debug("PathToken %v", value)
		op := &Operation{OperationType: TraversePath, Value: value, StringValue: value}
		return &Token{TokenType: OperationToken, Operation: op, CheckForPostTraverse: true}, nil
	}
}

func documentToken() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		numberString = numberString[1:]
		var number, errParsingInt = strconv.ParseInt(numberString, 10, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}
		log.Debug("documentToken %v", string(m.Bytes))
		op := &Operation{OperationType: DocumentFilter, Value: number, StringValue: numberString}
		return &Token{TokenType: OperationToken, Operation: op, CheckForPostTraverse: true}, nil
	}
}

func opToken(op *OperationType) lex.Action {
	return opTokenWithPrefs(op, nil, nil)
}

func opAssignableToken(opType *OperationType, assignOpType *OperationType) lex.Action {
	return opTokenWithPrefs(opType, assignOpType, nil)
}

func opTokenWithPrefs(op *OperationType, assignOpType *OperationType, preferences interface{}) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		log.Debug("opTokenWithPrefs %v", string(m.Bytes))
		value := string(m.Bytes)
		op := &Operation{OperationType: op, Value: op.Type, StringValue: value, Preferences: preferences}
		var assign *Operation
		if assignOpType != nil {
			assign = &Operation{OperationType: assignOpType, Value: assignOpType.Type, StringValue: value, Preferences: preferences}
		}
		return &Token{TokenType: OperationToken, Operation: op, AssignOperation: assign}, nil
	}
}

func literalToken(pType TokenType, checkForPost bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &Token{TokenType: pType, CheckForPostTraverse: checkForPost}, nil
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
		op := &Operation{OperationType: TraversePath, Value: number, StringValue: numberString}
		return &Token{TokenType: OperationToken, Operation: op, CheckForPostTraverse: true}, nil
	}
}

func numberValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		var number, errParsingInt = strconv.ParseInt(numberString, 10, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		return &Token{TokenType: OperationToken, Operation: CreateValueOperation(number, numberString)}, nil
	}
}

func floatValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		var number, errParsingInt = strconv.ParseFloat(numberString, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}
		return &Token{TokenType: OperationToken, Operation: CreateValueOperation(number, numberString)}, nil
	}
}

func booleanValue(val bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &Token{TokenType: OperationToken, Operation: CreateValueOperation(val, string(m.Bytes))}, nil
	}
}

func stringValue(wrapped bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		if wrapped {
			value = unwrap(value)
		}
		return &Token{TokenType: OperationToken, Operation: CreateValueOperation(value, value)}, nil
	}
}

func nullValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &Token{TokenType: OperationToken, Operation: CreateValueOperation(nil, string(m.Bytes))}, nil
	}
}

func selfToken() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		op := &Operation{OperationType: SelfReference}
		return &Token{TokenType: OperationToken, Operation: op}, nil
	}
}

func initLexer() (*lex.Lexer, error) {
	lexer := lex.NewLexer()
	lexer.Add([]byte(`\(`), literalToken(OpenBracket, false))
	lexer.Add([]byte(`\)`), literalToken(CloseBracket, true))

	lexer.Add([]byte(`\.\[\]`), pathToken(false))
	lexer.Add([]byte(`\.\.`), opToken(RecursiveDescent))

	lexer.Add([]byte(`,`), opToken(Union))
	lexer.Add([]byte(`:\s*`), opToken(CreateMap))
	lexer.Add([]byte(`length`), opToken(Length))
	lexer.Add([]byte(`sortKeys`), opToken(SortKeys))
	lexer.Add([]byte(`select`), opToken(Select))
	lexer.Add([]byte(`has`), opToken(Has))
	lexer.Add([]byte(`explode`), opToken(Explode))
	lexer.Add([]byte(`or`), opToken(Or))
	lexer.Add([]byte(`and`), opToken(And))
	lexer.Add([]byte(`not`), opToken(Not))

	lexer.Add([]byte(`documentIndex`), opToken(GetDocumentIndex))

	lexer.Add([]byte(`style`), opAssignableToken(GetStyle, AssignStyle))

	lexer.Add([]byte(`tag`), opAssignableToken(GetTag, AssignTag))
	lexer.Add([]byte(`filename`), opToken(GetFilename))
	lexer.Add([]byte(`fileIndex`), opToken(GetFileIndex))
	lexer.Add([]byte(`path`), opToken(GetPath))

	lexer.Add([]byte(`lineComment`), opTokenWithPrefs(GetComment, AssignComment, &CommentOpPreferences{LineComment: true}))

	lexer.Add([]byte(`headComment`), opTokenWithPrefs(GetComment, AssignComment, &CommentOpPreferences{HeadComment: true}))

	lexer.Add([]byte(`footComment`), opTokenWithPrefs(GetComment, AssignComment, &CommentOpPreferences{FootComment: true}))

	lexer.Add([]byte(`comments\s*=`), opTokenWithPrefs(AssignComment, nil, &CommentOpPreferences{LineComment: true, HeadComment: true, FootComment: true}))

	lexer.Add([]byte(`collect`), opToken(Collect))

	lexer.Add([]byte(`\s*==\s*`), opToken(Equals))
	lexer.Add([]byte(`\s*=\s*`), opTokenWithPrefs(Assign, nil, &AssignOpPreferences{false}))

	lexer.Add([]byte(`del`), opToken(DeleteChild))

	lexer.Add([]byte(`\s*\|=\s*`), opTokenWithPrefs(Assign, nil, &AssignOpPreferences{true}))

	lexer.Add([]byte(`\.\[-?[0-9]+\]`), arrayIndextoken(true))

	lexer.Add([]byte("( |\t|\n|\r)+"), skip)

	lexer.Add([]byte(`d[0-9]+`), documentToken())
	lexer.Add([]byte(`\."[^ "]+"`), pathToken(true))
	lexer.Add([]byte(`\.[^ \}\{\:\[\],\|\.\[\(\)=]+`), pathToken(false))
	lexer.Add([]byte(`\.`), selfToken())

	lexer.Add([]byte(`\|`), opToken(Pipe))

	lexer.Add([]byte(`-?\d+(\.\d+)`), floatValue())
	lexer.Add([]byte(`-?[1-9](\.\d+)?[Ee][-+]?\d+`), floatValue())
	lexer.Add([]byte(`-?\d+`), numberValue())

	lexer.Add([]byte(`[Tt][Rr][Uu][Ee]`), booleanValue(true))
	lexer.Add([]byte(`[Ff][Aa][Ll][Ss][Ee]`), booleanValue(false))

	lexer.Add([]byte(`[Nn][Uu][Ll][Ll]`), nullValue())
	lexer.Add([]byte(`~`), nullValue())

	lexer.Add([]byte(`"[^ "]*"`), stringValue(true))

	lexer.Add([]byte(`\[\]`), literalToken(SplatOrEmptyCollect, true))

	lexer.Add([]byte(`\[`), literalToken(OpenCollect, false))
	lexer.Add([]byte(`\]`), literalToken(CloseCollect, true))
	lexer.Add([]byte(`\{`), literalToken(OpenCollectObject, false))
	lexer.Add([]byte(`\}`), literalToken(CloseCollectObject, true))
	lexer.Add([]byte(`\*`), opTokenWithPrefs(Multiply, nil, &MultiplyPreferences{AppendArrays: false}))
	lexer.Add([]byte(`\*\+`), opTokenWithPrefs(Multiply, nil, &MultiplyPreferences{AppendArrays: true}))
	lexer.Add([]byte(`\+`), opToken(Add))
	lexer.Add([]byte(`\+=`), opToken(AddAssign))

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
		return nil, fmt.Errorf("Parsing expression: %v", err)
	}
	var tokens []*Token
	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {

		if tok != nil {
			token := tok.(*Token)
			log.Debugf("Tokenising %v", token.toString())
			tokens = append(tokens, token)
		}
		if err != nil {
			return nil, fmt.Errorf("Parsing expression: %v", err)
		}
	}
	var postProcessedTokens = make([]*Token, 0)

	skipNextToken := false

	for index := range tokens {
		if skipNextToken {
			skipNextToken = false
		} else {
			postProcessedTokens, skipNextToken = p.handleToken(tokens, index, postProcessedTokens)
		}
	}

	return postProcessedTokens, nil
}

func (p *pathTokeniser) handleToken(tokens []*Token, index int, postProcessedTokens []*Token) (tokensAccum []*Token, skipNextToken bool) {
	skipNextToken = false
	token := tokens[index]
	if token.TokenType == SplatOrEmptyCollect {
		if index > 0 && tokens[index-1].TokenType == OperationToken &&
			tokens[index-1].Operation.OperationType == TraversePath {
			// must be a splat without a preceding dot , e.g. .a[]
			// lets put a pipe in front of it, and convert it to a traverse "[]" token
			pipeOp := &Operation{OperationType: ShortPipe, Value: "PIPE"}

			postProcessedTokens = append(postProcessedTokens, &Token{TokenType: OperationToken, Operation: pipeOp})

			traverseOp := &Operation{OperationType: TraversePath, Value: "[]", StringValue: "[]"}
			token = &Token{TokenType: OperationToken, Operation: traverseOp, CheckForPostTraverse: true}

		} else {
			// gotta be a collect empty array, we need to split this into two tokens
			// one OpenCollect, the other CloseCollect
			postProcessedTokens = append(postProcessedTokens, &Token{TokenType: OpenCollect})
			token = &Token{TokenType: CloseCollect, CheckForPostTraverse: true}
		}
	}

	if index != len(tokens)-1 && token.AssignOperation != nil &&
		tokens[index+1].TokenType == OperationToken &&
		tokens[index+1].Operation.OperationType == Assign {
		token.Operation = token.AssignOperation
		skipNextToken = true
	}

	postProcessedTokens = append(postProcessedTokens, token)

	if index != len(tokens)-1 && token.CheckForPostTraverse &&
		tokens[index+1].TokenType == OperationToken &&
		tokens[index+1].Operation.OperationType == TraversePath {
		op := &Operation{OperationType: ShortPipe, Value: "PIPE"}
		postProcessedTokens = append(postProcessedTokens, &Token{TokenType: OperationToken, Operation: op})
	}
	return postProcessedTokens, skipNextToken
}
