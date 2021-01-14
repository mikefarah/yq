package yqlib

import (
	"fmt"
	"strconv"
	"strings"

	lex "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func skip(*lex.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

type tokenType uint32

const (
	operationToken = 1 << iota
	openBracket
	closeBracket
	openCollect
	closeCollect
	openCollectObject
	closeCollectObject
	traverseArrayCollect
)

type token struct {
	TokenType            tokenType
	Operation            *Operation
	AssignOperation      *Operation // e.g. tag (GetTag) op becomes AssignTag if '=' follows it
	CheckForPostTraverse bool       // e.g. [1]cat should really be [1].cat

}

func (t *token) toString() string {
	if t.TokenType == operationToken {
		log.Debug("toString, its an op")
		return t.Operation.toString()
	} else if t.TokenType == openBracket {
		return "("
	} else if t.TokenType == closeBracket {
		return ")"
	} else if t.TokenType == openCollect {
		return "["
	} else if t.TokenType == closeCollect {
		return "]"
	} else if t.TokenType == openCollectObject {
		return "{"
	} else if t.TokenType == closeCollectObject {
		return "}"
	} else if t.TokenType == traverseArrayCollect {
		return ".["

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
		op := &Operation{OperationType: traversePathOpType, Value: value, StringValue: value, Preferences: traversePreferences{}}
		return &token{TokenType: operationToken, Operation: op, CheckForPostTraverse: true}, nil
	}
}

func opToken(op *operationType) lex.Action {
	return opTokenWithPrefs(op, nil, nil)
}

func opAssignableToken(opType *operationType, assignOpType *operationType) lex.Action {
	return opTokenWithPrefs(opType, assignOpType, nil)
}

func assignOpToken(updateAssign bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		log.Debug("assignOpToken %v", string(m.Bytes))
		value := string(m.Bytes)
		op := &Operation{OperationType: assignOpType, Value: assignOpType.Type, StringValue: value, UpdateAssign: updateAssign}
		return &token{TokenType: operationToken, Operation: op}, nil
	}
}

func multiplyWithPrefs() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		prefs := multiplyPreferences{}
		options := string(m.Bytes)
		if strings.Contains(options, "+") {
			prefs.AppendArrays = true
		}
		if strings.Contains(options, "?") {
			prefs.TraversePrefs = traversePreferences{DontAutoCreate: true}
		}
		op := &Operation{OperationType: multiplyOpType, Value: multiplyOpType.Type, StringValue: options, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op}, nil
	}
}

func opTokenWithPrefs(op *operationType, assignOpType *operationType, preferences interface{}) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		log.Debug("opTokenWithPrefs %v", string(m.Bytes))
		value := string(m.Bytes)
		op := &Operation{OperationType: op, Value: op.Type, StringValue: value, Preferences: preferences}
		var assign *Operation
		if assignOpType != nil {
			assign = &Operation{OperationType: assignOpType, Value: assignOpType.Type, StringValue: value, Preferences: preferences}
		}
		return &token{TokenType: operationToken, Operation: op, AssignOperation: assign}, nil
	}
}

func assignAllCommentsOp(updateAssign bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		log.Debug("assignAllCommentsOp %v", string(m.Bytes))
		value := string(m.Bytes)
		op := &Operation{
			OperationType: assignCommentOpType,
			Value:         assignCommentOpType.Type,
			StringValue:   value,
			UpdateAssign:  updateAssign,
			Preferences:   commentOpPreferences{LineComment: true, HeadComment: true, FootComment: true},
		}
		return &token{TokenType: operationToken, Operation: op}, nil
	}
}

func literalToken(pType tokenType, checkForPost bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &token{TokenType: pType, CheckForPostTraverse: checkForPost}, nil
	}
}

func unwrap(value string) string {
	return value[1 : len(value)-1]
}

func numberValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		var number, errParsingInt = strconv.ParseInt(numberString, 10, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		return &token{TokenType: operationToken, Operation: createValueOperation(number, numberString)}, nil
	}
}

func floatValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		var number, errParsingInt = strconv.ParseFloat(numberString, 64) // nolint
		if errParsingInt != nil {
			return nil, errParsingInt
		}
		return &token{TokenType: operationToken, Operation: createValueOperation(number, numberString)}, nil
	}
}

func booleanValue(val bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &token{TokenType: operationToken, Operation: createValueOperation(val, string(m.Bytes))}, nil
	}
}

func stringValue(wrapped bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		if wrapped {
			value = unwrap(value)
		}
		return &token{TokenType: operationToken, Operation: createValueOperation(value, value)}, nil
	}
}

func envOp(strenv bool) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		preferences := envOpPreferences{}

		if strenv {
			// strenv( )
			value = value[7 : len(value)-1]
			preferences.StringValue = true
		} else {
			//env( )
			value = value[4 : len(value)-1]
		}

		envOperation := createValueOperation(value, value)
		envOperation.OperationType = envOpType
		envOperation.Preferences = preferences

		return &token{TokenType: operationToken, Operation: envOperation}, nil
	}
}

func nullValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return &token{TokenType: operationToken, Operation: createValueOperation(nil, string(m.Bytes))}, nil
	}
}

func selfToken() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		op := &Operation{OperationType: selfReferenceOpType}
		return &token{TokenType: operationToken, Operation: op}, nil
	}
}

func initLexer() (*lex.Lexer, error) {
	lexer := lex.NewLexer()
	lexer.Add([]byte(`\(`), literalToken(openBracket, false))
	lexer.Add([]byte(`\)`), literalToken(closeBracket, true))

	lexer.Add([]byte(`\.\[`), literalToken(traverseArrayCollect, false))
	lexer.Add([]byte(`\.\.`), opTokenWithPrefs(recursiveDescentOpType, nil, recursiveDescentPreferences{RecurseArray: true,
		TraversePreferences: traversePreferences{DontFollowAlias: true, IncludeMapKeys: false}}))

	lexer.Add([]byte(`\.\.\.`), opTokenWithPrefs(recursiveDescentOpType, nil, recursiveDescentPreferences{RecurseArray: true,
		TraversePreferences: traversePreferences{DontFollowAlias: true, IncludeMapKeys: true}}))

	lexer.Add([]byte(`,`), opToken(unionOpType))
	lexer.Add([]byte(`:\s*`), opToken(createMapOpType))
	lexer.Add([]byte(`length`), opToken(lengthOpType))
	lexer.Add([]byte(`sortKeys`), opToken(sortKeysOpType))
	lexer.Add([]byte(`select`), opToken(selectOpType))
	lexer.Add([]byte(`has`), opToken(hasOpType))
	lexer.Add([]byte(`explode`), opToken(explodeOpType))
	lexer.Add([]byte(`or`), opToken(orOpType))
	lexer.Add([]byte(`and`), opToken(andOpType))
	lexer.Add([]byte(`not`), opToken(notOpType))
	lexer.Add([]byte(`\/\/`), opToken(alternativeOpType))

	lexer.Add([]byte(`documentIndex`), opToken(getDocumentIndexOpType))
	lexer.Add([]byte(`di`), opToken(getDocumentIndexOpType))
	lexer.Add([]byte(`splitDoc`), opToken(splitDocumentOpType))

	lexer.Add([]byte(`join`), opToken(joinStringOpType))

	lexer.Add([]byte(`style`), opAssignableToken(getStyleOpType, assignStyleOpType))

	lexer.Add([]byte(`tag`), opAssignableToken(getTagOpType, assignTagOpType))
	lexer.Add([]byte(`anchor`), opAssignableToken(getAnchorOpType, assignAnchorOpType))
	lexer.Add([]byte(`alias`), opAssignableToken(getAliasOptype, assignAliasOpType))
	lexer.Add([]byte(`filename`), opToken(getFilenameOpType))
	lexer.Add([]byte(`fileIndex`), opToken(getFileIndexOpType))
	lexer.Add([]byte(`fi`), opToken(getFileIndexOpType))
	lexer.Add([]byte(`path`), opToken(getPathOpType))

	lexer.Add([]byte(`lineComment`), opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{LineComment: true}))

	lexer.Add([]byte(`headComment`), opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{HeadComment: true}))

	lexer.Add([]byte(`footComment`), opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{FootComment: true}))

	lexer.Add([]byte(`comments\s*=`), assignAllCommentsOp(false))
	lexer.Add([]byte(`comments\s*\|=`), assignAllCommentsOp(true))

	lexer.Add([]byte(`collect`), opToken(collectOpType))

	lexer.Add([]byte(`\s*==\s*`), opToken(equalsOpType))
	lexer.Add([]byte(`\s*=\s*`), assignOpToken(false))

	lexer.Add([]byte(`del`), opToken(deleteChildOpType))

	lexer.Add([]byte(`\s*\|=\s*`), assignOpToken(true))

	lexer.Add([]byte("( |\t|\n|\r)+"), skip)

	lexer.Add([]byte(`\."[^ "]+"`), pathToken(true))
	lexer.Add([]byte(`\.[^ \}\{\:\[\],\|\.\[\(\)=]+`), pathToken(false))
	lexer.Add([]byte(`\.`), selfToken())

	lexer.Add([]byte(`\|`), opToken(pipeOpType))

	lexer.Add([]byte(`-?\d+(\.\d+)`), floatValue())
	lexer.Add([]byte(`-?[1-9](\.\d+)?[Ee][-+]?\d+`), floatValue())
	lexer.Add([]byte(`-?\d+`), numberValue())

	lexer.Add([]byte(`[Tt][Rr][Uu][Ee]`), booleanValue(true))
	lexer.Add([]byte(`[Ff][Aa][Ll][Ss][Ee]`), booleanValue(false))

	lexer.Add([]byte(`[Nn][Uu][Ll][Ll]`), nullValue())
	lexer.Add([]byte(`~`), nullValue())

	lexer.Add([]byte(`"[^"]*"`), stringValue(true))
	lexer.Add([]byte(`strenv\([^\)]+\)`), envOp(true))
	lexer.Add([]byte(`env\([^\)]+\)`), envOp(false))

	lexer.Add([]byte(`\[`), literalToken(openCollect, false))
	lexer.Add([]byte(`\]`), literalToken(closeCollect, true))
	lexer.Add([]byte(`\{`), literalToken(openCollectObject, false))
	lexer.Add([]byte(`\}`), literalToken(closeCollectObject, true))
	lexer.Add([]byte(`\*[\+|\?]*`), multiplyWithPrefs())
	lexer.Add([]byte(`\+`), opToken(addOpType))
	lexer.Add([]byte(`\+=`), opToken(addAssignOpType))

	err := lexer.Compile()
	if err != nil {
		return nil, err
	}
	return lexer, nil
}

type expressionTokeniser interface {
	Tokenise(expression string) ([]*token, error)
}

type expressionTokeniserImpl struct {
	lexer *lex.Lexer
}

func newExpressionTokeniser() expressionTokeniser {
	var lexer, err = initLexer()
	if err != nil {
		panic(err)
	}
	return &expressionTokeniserImpl{lexer}
}

func (p *expressionTokeniserImpl) Tokenise(expression string) ([]*token, error) {
	scanner, err := p.lexer.Scanner([]byte(expression))

	if err != nil {
		return nil, fmt.Errorf("Parsing expression: %v", err)
	}
	var tokens []*token
	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {

		if tok != nil {
			currentToken := tok.(*token)
			log.Debugf("Tokenising %v", currentToken.toString())
			tokens = append(tokens, currentToken)
		}
		if err != nil {
			return nil, fmt.Errorf("Parsing expression: %v", err)
		}
	}
	var postProcessedTokens = make([]*token, 0)

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

func (p *expressionTokeniserImpl) handleToken(tokens []*token, index int, postProcessedTokens []*token) (tokensAccum []*token, skipNextToken bool) {
	skipNextToken = false
	currentToken := tokens[index]

	if currentToken.TokenType == traverseArrayCollect {
		//need to put a traverse array then a collect currentToken
		// do this by adding traverse then converting currentToken to collect

		op := &Operation{OperationType: traverseArrayOpType, StringValue: "TRAVERSE_ARRAY"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})

		currentToken = &token{TokenType: openCollect}

	}

	if index != len(tokens)-1 && currentToken.AssignOperation != nil &&
		tokens[index+1].TokenType == operationToken &&
		tokens[index+1].Operation.OperationType == assignOpType {
		currentToken.Operation = currentToken.AssignOperation
		currentToken.Operation.UpdateAssign = tokens[index+1].Operation.UpdateAssign
		skipNextToken = true
	}

	postProcessedTokens = append(postProcessedTokens, currentToken)

	if index != len(tokens)-1 && currentToken.CheckForPostTraverse &&
		tokens[index+1].TokenType == operationToken &&
		tokens[index+1].Operation.OperationType == traversePathOpType {
		op := &Operation{OperationType: shortPipeOpType, Value: "PIPE"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}
	if index != len(tokens)-1 && currentToken.CheckForPostTraverse &&
		tokens[index+1].TokenType == openCollect {

		op := &Operation{OperationType: shortPipeOpType, Value: "PIPE"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})

		op = &Operation{OperationType: traverseArrayOpType}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}
	if index != len(tokens)-1 && currentToken.CheckForPostTraverse &&
		tokens[index+1].TokenType == traverseArrayCollect {

		op := &Operation{OperationType: shortPipeOpType, Value: "PIPE"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})

	}
	return postProcessedTokens, skipNextToken
}
