package yqlib

import (
	"fmt"
	"regexp"
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
	AssignOperation      *Operation      // e.g. tag (GetTag) op becomes AssignTag if '=' follows it
	CheckForPostTraverse bool            // e.g. [1]cat should really be [1].cat
	Match                *machines.Match // match that created this token

}

func (t *token) toString(detail bool) string {
	if t.TokenType == operationToken {
		if detail {
			return fmt.Sprintf("%v (%v)", t.Operation.toString(), t.Operation.OperationType.Precedence)
		}
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
		prefs := traversePreferences{}

		if value[len(value)-1:] == "?" {
			prefs.OptionalTraverse = true
			value = value[:len(value)-1]
		}

		value = value[1:]
		if wrapped {
			value = unwrap(value)
		}
		log.Debug("PathToken %v", value)
		op := &Operation{OperationType: traversePathOpType, Value: value, StringValue: value, Preferences: prefs}
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
		prefs := assignPreferences{DontOverWriteAnchor: true}
		op := &Operation{OperationType: assignOpType, Value: assignOpType.Type, StringValue: value, UpdateAssign: updateAssign, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op}, nil
	}
}

func multiplyWithPrefs(op *operationType) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		prefs := multiplyPreferences{}
		options := string(m.Bytes)
		if strings.Contains(options, "+") {
			prefs.AppendArrays = true
		}
		if strings.Contains(options, "?") {
			prefs.TraversePrefs = traversePreferences{DontAutoCreate: true}
		}
		if strings.Contains(options, "n") {
			prefs.AssignPrefs = assignPreferences{OnlyWriteNull: true}
		}
		if strings.Contains(options, "d") {
			prefs.DeepMergeArrays = true
		}
		prefs.TraversePrefs.DontFollowAlias = true
		op := &Operation{OperationType: op, Value: multiplyOpType.Type, StringValue: options, Preferences: prefs}
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

func hasOptionParameter(value string, option string) bool {
	parameterParser := regexp.MustCompile(`.*\([^\)]*\)`)
	matches := parameterParser.FindStringSubmatch(value)
	if len(matches) == 0 {
		return false
	}
	parameterString := matches[0]
	optionParser := regexp.MustCompile(fmt.Sprintf("\\b%v\\b", option))
	return len(optionParser.FindStringSubmatch(parameterString)) > 0
}

func extractNumberParameter(value string) (int, error) {
	parameterParser := regexp.MustCompile(`.*\(([0-9]+)\)`)
	matches := parameterParser.FindStringSubmatch(value)
	var indent, errParsingInt = strconv.ParseInt(matches[1], 10, 32)
	if errParsingInt != nil {
		return 0, errParsingInt
	}
	return int(indent), nil
}

func envSubstWithOptions() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		noEmpty := hasOptionParameter(value, "ne")
		noUnset := hasOptionParameter(value, "nu")
		failFast := hasOptionParameter(value, "ff")
		envsubstOpType.Type = "ENVSUBST"
		prefs := envOpPreferences{NoUnset: noUnset, NoEmpty: noEmpty, FailFast: failFast}
		if noEmpty {
			envsubstOpType.Type = envsubstOpType.Type + "_NO_EMPTY"
		}
		if noUnset {
			envsubstOpType.Type = envsubstOpType.Type + "_NO_UNSET"
		}

		op := &Operation{OperationType: envsubstOpType, Value: envsubstOpType.Type, StringValue: value, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op}, nil
	}
}

func flattenWithDepth() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		var depth, errParsingInt = extractNumberParameter(value)
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		prefs := flattenPreferences{depth: depth}
		op := &Operation{OperationType: flattenOpType, Value: flattenOpType.Type, StringValue: value, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op}, nil
	}
}

func encodeWithIndent(outputFormat PrinterOutputFormat) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)
		var indent, errParsingInt = extractNumberParameter(value)
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		prefs := encoderPreferences{format: outputFormat, indent: indent}
		op := &Operation{OperationType: encodeOpType, Value: encodeOpType.Type, StringValue: value, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op}, nil
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
		return &token{TokenType: pType, CheckForPostTraverse: checkForPost, Match: m}, nil
	}
}

func unwrap(value string) string {
	return value[1 : len(value)-1]
}

func numberValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		var number, errParsingInt = strconv.ParseInt(numberString, 10, 64)
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		return &token{TokenType: operationToken, Operation: createValueOperation(number, numberString)}, nil
	}
}

func hexValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var originalString = string(m.Bytes)
		var numberString = originalString[2:]
		log.Debugf("numberString: %v", numberString)
		var number, errParsingInt = strconv.ParseInt(numberString, 16, 64)
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		return &token{TokenType: operationToken, Operation: createValueOperation(number, originalString)}, nil
	}
}

func floatValue() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		var numberString = string(m.Bytes)
		var number, errParsingInt = strconv.ParseFloat(numberString, 64)
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
		value = strings.ReplaceAll(value, "\\\"", "\"")
		return &token{TokenType: operationToken, Operation: createValueOperation(value, value)}, nil
	}
}

func getVariableOpToken() lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		value := string(m.Bytes)

		value = value[1:]

		getVarOperation := createValueOperation(value, value)
		getVarOperation.OperationType = getVariableOpType

		return &token{TokenType: operationToken, Operation: getVarOperation, CheckForPostTraverse: true}, nil
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
	lexer.Add([]byte(`line`), opToken(lineOpType))
	lexer.Add([]byte(`column`), opToken(columnOpType))

	lexer.Add([]byte(`eval`), opToken(evalOpType))

	lexer.Add([]byte(`map`), opToken(mapOpType))
	lexer.Add([]byte(`map_values`), opToken(mapValuesOpType))
	lexer.Add([]byte(`pick`), opToken(pickOpType))

	lexer.Add([]byte(`flatten\([0-9]+\)`), flattenWithDepth())
	lexer.Add([]byte(`flatten`), opTokenWithPrefs(flattenOpType, nil, flattenPreferences{depth: -1}))

	lexer.Add([]byte(`format_datetime`), opToken(formatDateTimeOpType))
	lexer.Add([]byte(`now`), opToken(nowOpType))
	lexer.Add([]byte(`tz`), opToken(tzOpType))
	lexer.Add([]byte(`with_dtf`), opToken(withDtFormatOpType))

	lexer.Add([]byte(`toyaml\([0-9]+\)`), encodeWithIndent(YamlOutputFormat))
	lexer.Add([]byte(`to_yaml\([0-9]+\)`), encodeWithIndent(YamlOutputFormat))

	lexer.Add([]byte(`toxml\([0-9]+\)`), encodeWithIndent(XMLOutputFormat))
	lexer.Add([]byte(`to_xml\([0-9]+\)`), encodeWithIndent(XMLOutputFormat))

	lexer.Add([]byte(`tojson\([0-9]+\)`), encodeWithIndent(JSONOutputFormat))
	lexer.Add([]byte(`to_json\([0-9]+\)`), encodeWithIndent(JSONOutputFormat))

	lexer.Add([]byte(`toyaml`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: YamlOutputFormat, indent: 2}))
	lexer.Add([]byte(`to_yaml`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: YamlOutputFormat, indent: 2}))
	// 0 indent doesn't work with yaml.
	lexer.Add([]byte(`@yaml`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: YamlOutputFormat, indent: 2}))

	lexer.Add([]byte(`tojson`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: JSONOutputFormat, indent: 2}))
	lexer.Add([]byte(`to_json`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: JSONOutputFormat, indent: 2}))
	lexer.Add([]byte(`@json`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: JSONOutputFormat, indent: 0}))

	lexer.Add([]byte(`toprops`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: PropsOutputFormat, indent: 2}))
	lexer.Add([]byte(`to_props`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: PropsOutputFormat, indent: 2}))
	lexer.Add([]byte(`@props`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: PropsOutputFormat, indent: 2}))

	lexer.Add([]byte(`tocsv`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: CSVOutputFormat}))
	lexer.Add([]byte(`to_csv`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: CSVOutputFormat}))
	lexer.Add([]byte(`@csv`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: CSVOutputFormat}))

	lexer.Add([]byte(`totsv`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: TSVOutputFormat}))
	lexer.Add([]byte(`to_tsv`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: TSVOutputFormat}))
	lexer.Add([]byte(`@tsv`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: TSVOutputFormat}))

	lexer.Add([]byte(`toxml`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: XMLOutputFormat}))
	lexer.Add([]byte(`to_xml`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: XMLOutputFormat, indent: 2}))
	lexer.Add([]byte(`@xml`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: XMLOutputFormat, indent: 0}))

	lexer.Add([]byte(`@base64`), opTokenWithPrefs(encodeOpType, nil, encoderPreferences{format: Base64OutputFormat}))
	lexer.Add([]byte(`@base64d`), opTokenWithPrefs(decodeOpType, nil, decoderPreferences{format: Base64InputFormat}))

	lexer.Add([]byte(`fromyaml`), opTokenWithPrefs(decodeOpType, nil, decoderPreferences{format: YamlInputFormat}))
	lexer.Add([]byte(`fromjson`), opTokenWithPrefs(decodeOpType, nil, decoderPreferences{format: YamlInputFormat}))
	lexer.Add([]byte(`fromxml`), opTokenWithPrefs(decodeOpType, nil, decoderPreferences{format: XMLInputFormat}))

	lexer.Add([]byte(`from_yaml`), opTokenWithPrefs(decodeOpType, nil, decoderPreferences{format: YamlInputFormat}))
	lexer.Add([]byte(`from_json`), opTokenWithPrefs(decodeOpType, nil, decoderPreferences{format: YamlInputFormat}))
	lexer.Add([]byte(`from_xml`), opTokenWithPrefs(decodeOpType, nil, decoderPreferences{format: XMLInputFormat}))
	lexer.Add([]byte(`from_props`), opTokenWithPrefs(decodeOpType, nil, decoderPreferences{format: PropertiesInputFormat}))

	lexer.Add([]byte(`sortKeys`), opToken(sortKeysOpType))
	lexer.Add([]byte(`sort_keys`), opToken(sortKeysOpType))

	lexer.Add([]byte(`load`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: false, decoder: NewYamlDecoder()}))

	lexer.Add([]byte(`xmlload`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: false, decoder: NewXMLDecoder(XMLPreferences.AttributePrefix, XMLPreferences.ContentName, XMLPreferences.StrictMode)}))
	lexer.Add([]byte(`load_xml`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: false, decoder: NewXMLDecoder(XMLPreferences.AttributePrefix, XMLPreferences.ContentName, XMLPreferences.StrictMode)}))
	lexer.Add([]byte(`loadxml`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: false, decoder: NewXMLDecoder(XMLPreferences.AttributePrefix, XMLPreferences.ContentName, XMLPreferences.StrictMode)}))

	lexer.Add([]byte(`load_base64`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: false, decoder: NewBase64Decoder()}))

	lexer.Add([]byte(`load_props`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: false, decoder: NewPropertiesDecoder()}))
	lexer.Add([]byte(`loadprops`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: false, decoder: NewPropertiesDecoder()}))

	lexer.Add([]byte(`strload`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: true}))
	lexer.Add([]byte(`load_str`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: true}))
	lexer.Add([]byte(`loadstr`), opTokenWithPrefs(loadOpType, nil, loadPrefs{loadAsString: true}))

	lexer.Add([]byte(`select`), opToken(selectOpType))
	lexer.Add([]byte(`has`), opToken(hasOpType))
	lexer.Add([]byte(`unique`), opToken(uniqueOpType))
	lexer.Add([]byte(`unique_by`), opToken(uniqueByOpType))
	lexer.Add([]byte(`group_by`), opToken(groupByOpType))
	lexer.Add([]byte(`explode`), opToken(explodeOpType))
	lexer.Add([]byte(`or`), opToken(orOpType))
	lexer.Add([]byte(`and`), opToken(andOpType))
	lexer.Add([]byte(`not`), opToken(notOpType))
	lexer.Add([]byte(`ireduce`), opToken(reduceOpType))
	lexer.Add([]byte(`;`), opToken(blockOpType))
	lexer.Add([]byte(`\/\/`), opToken(alternativeOpType))

	lexer.Add([]byte(`documentIndex`), opToken(getDocumentIndexOpType))
	lexer.Add([]byte(`document_index`), opToken(getDocumentIndexOpType))

	lexer.Add([]byte(`di`), opToken(getDocumentIndexOpType))

	lexer.Add([]byte(`splitDoc`), opToken(splitDocumentOpType))
	lexer.Add([]byte(`split_doc`), opToken(splitDocumentOpType))

	lexer.Add([]byte(`join`), opToken(joinStringOpType))
	lexer.Add([]byte(`sub`), opToken(subStringOpType))
	lexer.Add([]byte(`match`), opToken(matchOpType))
	lexer.Add([]byte(`capture`), opToken(captureOpType))
	lexer.Add([]byte(`test`), opToken(testOpType))

	lexer.Add([]byte(`upcase`), opTokenWithPrefs(changeCaseOpType, nil, changeCasePrefs{ToUpperCase: true}))
	lexer.Add([]byte(`ascii_upcase`), opTokenWithPrefs(changeCaseOpType, nil, changeCasePrefs{ToUpperCase: true}))

	lexer.Add([]byte(`downcase`), opTokenWithPrefs(changeCaseOpType, nil, changeCasePrefs{ToUpperCase: false}))
	lexer.Add([]byte(`ascii_downcase`), opTokenWithPrefs(changeCaseOpType, nil, changeCasePrefs{ToUpperCase: false}))

	lexer.Add([]byte(`sort`), opToken(sortOpType))
	lexer.Add([]byte(`sort_by`), opToken(sortByOpType))
	lexer.Add([]byte(`reverse`), opToken(reverseOpType))

	lexer.Add([]byte(`any`), opToken(anyOpType))
	lexer.Add([]byte(`any_c`), opToken(anyConditionOpType))
	lexer.Add([]byte(`all`), opToken(allOpType))
	lexer.Add([]byte(`all_c`), opToken(allConditionOpType))
	lexer.Add([]byte(`contains`), opToken(containsOpType))

	lexer.Add([]byte(`split`), opToken(splitStringOpType))

	lexer.Add([]byte(`parent`), opToken(getParentOpType))
	lexer.Add([]byte(`key`), opToken(getKeyOpType))
	lexer.Add([]byte(`keys`), opToken(keysOpType))

	lexer.Add([]byte(`style`), opAssignableToken(getStyleOpType, assignStyleOpType))

	lexer.Add([]byte(`tag`), opAssignableToken(getTagOpType, assignTagOpType))
	lexer.Add([]byte(`anchor`), opAssignableToken(getAnchorOpType, assignAnchorOpType))
	lexer.Add([]byte(`alias`), opAssignableToken(getAliasOptype, assignAliasOpType))
	lexer.Add([]byte(`filename`), opToken(getFilenameOpType))

	lexer.Add([]byte(`fileIndex`), opToken(getFileIndexOpType))
	lexer.Add([]byte(`file_index`), opToken(getFileIndexOpType))

	lexer.Add([]byte(`fi`), opToken(getFileIndexOpType))
	lexer.Add([]byte(`path`), opToken(getPathOpType))
	lexer.Add([]byte(`to_entries`), opToken(toEntriesOpType))
	lexer.Add([]byte(`from_entries`), opToken(fromEntriesOpType))
	lexer.Add([]byte(`with_entries`), opToken(withEntriesOpType))

	lexer.Add([]byte(`with`), opToken(withOpType))

	lexer.Add([]byte(`lineComment`), opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{LineComment: true}))
	lexer.Add([]byte(`line_comment`), opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{LineComment: true}))

	lexer.Add([]byte(`headComment`), opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{HeadComment: true}))
	lexer.Add([]byte(`head_comment`), opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{HeadComment: true}))

	lexer.Add([]byte(`footComment`), opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{FootComment: true}))
	lexer.Add([]byte(`foot_comment`), opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{FootComment: true}))

	lexer.Add([]byte(`comments\s*=`), assignAllCommentsOp(false))
	lexer.Add([]byte(`comments\s*\|=`), assignAllCommentsOp(true))

	lexer.Add([]byte(`collect`), opToken(collectOpType))

	lexer.Add([]byte(`\s*==\s*`), opToken(equalsOpType))
	lexer.Add([]byte(`\s*!=\s*`), opToken(notEqualsOpType))

	lexer.Add([]byte(`\s*>=\s*`), opTokenWithPrefs(compareOpType, nil, compareTypePref{OrEqual: true, Greater: true}))
	lexer.Add([]byte(`\s*>\s*`), opTokenWithPrefs(compareOpType, nil, compareTypePref{OrEqual: false, Greater: true}))

	lexer.Add([]byte(`\s*<=\s*`), opTokenWithPrefs(compareOpType, nil, compareTypePref{OrEqual: true, Greater: false}))
	lexer.Add([]byte(`\s*<\s*`), opTokenWithPrefs(compareOpType, nil, compareTypePref{OrEqual: false, Greater: false}))

	lexer.Add([]byte(`\s*=\s*`), assignOpToken(false))

	lexer.Add([]byte(`del`), opToken(deleteChildOpType))

	lexer.Add([]byte(`\s*\|=\s*`), assignOpToken(true))

	lexer.Add([]byte("( |\t|\n|\r)+"), skip)

	lexer.Add([]byte(`\."[^ "]+"\??`), pathToken(true))
	lexer.Add([]byte(`\.[^ ;\}\{\:\[\],\|\.\[\(\)=\n]+\??`), pathToken(false))
	lexer.Add([]byte(`\.`), selfToken())

	lexer.Add([]byte(`\|`), opToken(pipeOpType))

	lexer.Add([]byte(`0[xX][0-9A-Fa-f]+`), hexValue())
	lexer.Add([]byte(`-?\d+(\.\d+)`), floatValue())
	lexer.Add([]byte(`-?[1-9](\.\d+)?[Ee][-+]?\d+`), floatValue())
	lexer.Add([]byte(`-?\d+`), numberValue())

	lexer.Add([]byte(`[Tt][Rr][Uu][Ee]`), booleanValue(true))
	lexer.Add([]byte(`[Ff][Aa][Ll][Ss][Ee]`), booleanValue(false))

	lexer.Add([]byte(`[Nn][Uu][Ll][Ll]`), nullValue())
	lexer.Add([]byte(`~`), nullValue())

	lexer.Add([]byte(`"([^"\\]*(\\.[^"\\]*)*)"`), stringValue(true))
	lexer.Add([]byte(`strenv\([^\)]+\)`), envOp(true))
	lexer.Add([]byte(`env\([^\)]+\)`), envOp(false))

	lexer.Add([]byte(`envsubst\((ne|nu|ff| |,)+\)`), envSubstWithOptions())
	lexer.Add([]byte(`envsubst`), opToken(envsubstOpType))

	lexer.Add([]byte(`\[`), literalToken(openCollect, false))
	lexer.Add([]byte(`\]\??`), literalToken(closeCollect, true))
	lexer.Add([]byte(`\{`), literalToken(openCollectObject, false))
	lexer.Add([]byte(`\}`), literalToken(closeCollectObject, true))
	lexer.Add([]byte(`\*=[\+|\?dn]*`), multiplyWithPrefs(multiplyAssignOpType))
	lexer.Add([]byte(`\*[\+|\?dn]*`), multiplyWithPrefs(multiplyOpType))

	lexer.Add([]byte(`\+`), opToken(addOpType))
	lexer.Add([]byte(`\+=`), opToken(addAssignOpType))

	lexer.Add([]byte(`\-`), opToken(subtractOpType))
	lexer.Add([]byte(`\-=`), opToken(subtractAssignOpType))
	lexer.Add([]byte(`\$[a-zA-Z_-0-9]+`), getVariableOpToken())
	lexer.Add([]byte(`as`), opTokenWithPrefs(assignVariableOpType, nil, assignVarPreferences{}))
	lexer.Add([]byte(`ref`), opTokenWithPrefs(assignVariableOpType, nil, assignVarPreferences{IsReference: true}))

	err := lexer.CompileNFA()
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
		return nil, fmt.Errorf("parsing expression: %w", err)
	}
	var tokens []*token
	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {

		if tok != nil {
			currentToken := tok.(*token)
			log.Debugf("Tokenising %v", currentToken.toString(true))
			tokens = append(tokens, currentToken)
		}
		if err != nil {
			return nil, fmt.Errorf("parsing expression: %w", err)
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

	log.Debug("processing %v", currentToken.toString(true))

	if currentToken.TokenType == traverseArrayCollect {
		// `.[exp]`` works by creating a traversal array of [self, exp] and piping that into the traverse array operator
		//need to put a traverse array then a collect currentToken
		// do this by adding traverse then converting currentToken to collect

		log.Debug("  adding self")
		op := &Operation{OperationType: selfReferenceOpType, StringValue: "SELF"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})

		log.Debug("  adding traverse array")
		op = &Operation{OperationType: traverseArrayOpType, StringValue: "TRAVERSE_ARRAY"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})

		currentToken = &token{TokenType: openCollect}

	}

	if index != len(tokens)-1 && currentToken.AssignOperation != nil &&
		tokens[index+1].TokenType == operationToken &&
		tokens[index+1].Operation.OperationType == assignOpType {
		log.Debug("  its an update assign")
		currentToken.Operation = currentToken.AssignOperation
		currentToken.Operation.UpdateAssign = tokens[index+1].Operation.UpdateAssign
		skipNextToken = true
	}

	log.Debug("  adding token to the fixed list")
	postProcessedTokens = append(postProcessedTokens, currentToken)

	if index != len(tokens)-1 &&
		((currentToken.TokenType == openCollect && tokens[index+1].TokenType == closeCollect) ||
			(currentToken.TokenType == openCollectObject && tokens[index+1].TokenType == closeCollectObject)) {
		log.Debug("  adding empty")
		op := &Operation{OperationType: emptyOpType, StringValue: "EMPTY"}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}

	if index != len(tokens)-1 && currentToken.CheckForPostTraverse &&
		((tokens[index+1].TokenType == operationToken && (tokens[index+1].Operation.OperationType == traversePathOpType)) ||
			(tokens[index+1].TokenType == traverseArrayCollect)) {
		log.Debug("  adding pipe because the next thing is traverse")
		op := &Operation{OperationType: shortPipeOpType, Value: "PIPE", StringValue: "."}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}
	if index != len(tokens)-1 && currentToken.CheckForPostTraverse &&
		tokens[index+1].TokenType == openCollect {

		log.Debug("  adding traverArray because next is opencollect")
		op := &Operation{OperationType: traverseArrayOpType}
		postProcessedTokens = append(postProcessedTokens, &token{TokenType: operationToken, Operation: op})
	}
	return postProcessedTokens, skipNextToken
}
