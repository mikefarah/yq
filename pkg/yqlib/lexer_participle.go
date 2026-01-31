package yqlib

import (
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

var participleYqRules = []*participleYqRule{
	{"LINE_COMMENT", `line_?comment|lineComment`, opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{LineComment: true}), 0},
	{"HEAD_COMMENT", `head_?comment|headComment`, opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{HeadComment: true}), 0},
	{"FOOT_COMMENT", `foot_?comment|footComment`, opTokenWithPrefs(getCommentOpType, assignCommentOpType, commentOpPreferences{FootComment: true}), 0},

	{"OpenBracket", `\(`, literalToken(openBracket, false), 0},
	{"CloseBracket", `\)`, literalToken(closeBracket, true), 0},
	{"OpenTraverseArrayCollect", `\.\[`, literalToken(traverseArrayCollect, false), 0},

	{"OpenCollect", `\[`, literalToken(openCollect, false), 0},
	{"CloseCollect", `\]\??`, literalToken(closeCollect, true), 0},

	{"OpenCollectObject", `\{`, literalToken(openCollectObject, false), 0},
	{"CloseCollectObject", `\}`, literalToken(closeCollectObject, true), 0},

	{"RecursiveDecentIncludingKeys", `\.\.\.`, recursiveDecentOpToken(true), 0},
	{"RecursiveDecent", `\.\.`, recursiveDecentOpToken(false), 0},

	{"GetVariable", `\$[a-zA-Z_\-0-9]+`, getVariableOpToken(), 0},
	{"AssignAsVariable", `as`, opTokenWithPrefs(assignVariableOpType, nil, assignVarPreferences{}), 0},
	{"AssignRefVariable", `ref`, opTokenWithPrefs(assignVariableOpType, nil, assignVarPreferences{IsReference: true}), 0},

	{"CreateMap", `:\s*`, opToken(createMapOpType), 0},
	simpleOp("length", lengthOpType),
	simpleOp("line", lineOpType),
	simpleOp("column", columnOpType),
	simpleOp("eval", evalOpType),
	simpleOp("to_?number", toNumberOpType),

	{"MapValues", `map_?values`, opToken(mapValuesOpType), 0},
	simpleOp("map", mapOpType),
	simpleOp("filter", filterOpType),
	simpleOp("pick", pickOpType),
	simpleOp("omit", omitOpType),

	{"FlattenWithDepth", `flatten\([0-9]+\)`, flattenWithDepth(), 0},
	{"Flatten", `flatten`, opTokenWithPrefs(flattenOpType, nil, flattenPreferences{depth: -1}), 0},

	simpleOp("format_datetime", formatDateTimeOpType),
	simpleOp("now", nowOpType),
	simpleOp("tz", tzOpType),
	simpleOp("from_?unix", fromUnixOpType),
	simpleOp("to_?unix", toUnixOpType),
	simpleOp("with_dtf", withDtFormatOpType),
	simpleOp("error", errorOpType),
	simpleOp("shuffle", shuffleOpType),
	simpleOp("sortKeys", sortKeysOpType),
	simpleOp("sort_?keys", sortKeysOpType),

	{"ArrayToMap", "array_?to_?map", expressionOpToken(`(.[] | select(. != null) ) as $i ireduce({}; .[$i | key] = $i)`), 0},
	{"Root", "root", expressionOpToken(`parent(-1)`), 0},
	{"YamlEncodeWithIndent", `to_?yaml\([0-9]+\)`, encodeParseIndent(YamlFormat), 0},
	{"XMLEncodeWithIndent", `to_?xml\([0-9]+\)`, encodeParseIndent(XMLFormat), 0},
	{"JSONEncodeWithIndent", `to_?json\([0-9]+\)`, encodeParseIndent(JSONFormat), 0},

	{"YamlDecode", `from_?yaml|@yamld|from_?json|@jsond`, decodeOp(YamlFormat), 0},
	{"YamlEncode", `to_?yaml|@yaml`, encodeWithIndent(YamlFormat, 2), 0},

	{"JSONEncode", `to_?json`, encodeWithIndent(JSONFormat, 2), 0},
	{"JSONEncodeNoIndent", `@json`, encodeWithIndent(JSONFormat, 0), 0},

	{"PropertiesDecode", `from_?props|@propsd`, decodeOp(PropertiesFormat), 0},
	{"PropsEncode", `to_?props|@props`, encodeWithIndent(PropertiesFormat, 2), 0},

	{"XmlDecode", `from_?xml|@xmld`, decodeOp(XMLFormat), 0},
	{"XMLEncode", `to_?xml`, encodeWithIndent(XMLFormat, 2), 0},
	{"XMLEncodeNoIndent", `@xml`, encodeWithIndent(XMLFormat, 0), 0},

	{"CSVDecode", `from_?csv|@csvd`, decodeOp(CSVFormat), 0},
	{"CSVEncode", `to_?csv|@csv`, encodeWithIndent(CSVFormat, 0), 0},

	{"TSVDecode", `from_?tsv|@tsvd`, decodeOp(TSVFormat), 0},
	{"TSVEncode", `to_?tsv|@tsv`, encodeWithIndent(TSVFormat, 0), 0},

	{"Base64d", `@base64d`, decodeOp(Base64Format), 0},
	{"Base64", `@base64`, encodeWithIndent(Base64Format, 0), 0},

	{"Urid", `@urid`, decodeOp(UriFormat), 0},
	{"Uri", `@uri`, encodeWithIndent(UriFormat, 0), 0},
	{"SH", `@sh`, encodeWithIndent(ShFormat, 0), 0},

	{"LoadXML", `load_?xml|xml_?load`, loadOp(NewXMLDecoder(ConfiguredXMLPreferences)), 0},

	{"LoadBase64", `load_?base64`, loadOp(NewBase64Decoder()), 0},

	{"LoadProperties", `load_?props`, loadOp(NewPropertiesDecoder()), 0},
	simpleOp("load_?str|str_?load", loadStringOpType),
	{"LoadYaml", `load`, loadOp(NewYamlDecoder(LoadYamlPreferences)), 0},

	{"SplitDocument", `splitDoc|split_?doc`, opToken(splitDocumentOpType), 0},

	simpleOp("select", selectOpType),
	simpleOp("has", hasOpType),
	simpleOp("unique_?by", uniqueByOpType),
	simpleOp("unique", uniqueOpType),

	simpleOp("group_?by", groupByOpType),
	simpleOp("explode", explodeOpType),
	simpleOp("or", orOpType),
	simpleOp("and", andOpType),
	simpleOp("not", notOpType),
	simpleOp("ireduce", reduceOpType),

	simpleOp("join", joinStringOpType),
	simpleOp("sub", subStringOpType),
	simpleOp("match", matchOpType),
	simpleOp("capture", captureOpType),
	simpleOp("test", testOpType),

	simpleOp("sort_?by", sortByOpType),
	simpleOp("sort", sortOpType),
	simpleOp("first", firstOpType),

	simpleOp("reverse", reverseOpType),

	simpleOp("any_c", anyConditionOpType),
	simpleOp("any", anyOpType),

	simpleOp("all_c", allConditionOpType),
	simpleOp("all", allOpType),

	simpleOp("contains", containsOpType),
	simpleOp("split", splitStringOpType),

	simpleOp("parents", getParentsOpType),
	{"ParentWithLevel", `parent\(-?[0-9]+\)`, parentWithLevel(), 0},
	{"ParentWithDefaultLevel", `parent`, parentWithDefaultLevel(), 0},

	simpleOp("keys", keysOpType),
	simpleOp("key", getKeyOpType),
	simpleOp("is_?key", isKeyOpType),

	simpleOp("file_?name|fileName", getFilenameOpType),
	simpleOp("file_?index|fileIndex|fi", getFileIndexOpType),
	simpleOp("path", getPathOpType),
	simpleOp("set_?path", setPathOpType),
	simpleOp("del_?paths", delPathsOpType),

	simpleOp("to_?entries|toEntries", toEntriesOpType),
	simpleOp("from_?entries|fromEntries", fromEntriesOpType),
	simpleOp("with_?entries|withEntries", withEntriesOpType),

	simpleOp("with", withOpType),

	simpleOp("collect", collectOpType),
	simpleOp("del", deleteChildOpType),

	assignableOp("style", getStyleOpType, assignStyleOpType),
	assignableOp("tag|type", getTagOpType, assignTagOpType),
	simpleOp("kind", getKindOpType),
	assignableOp("anchor", getAnchorOpType, assignAnchorOpType),
	assignableOp("alias", getAliasOpType, assignAliasOpType),

	{"ALL_COMMENTS", `comments\s*=`, assignAllCommentsOp(false), 0},
	{"ALL_COMMENTS_ASSIGN_RELATIVE", `comments\s*\|=`, assignAllCommentsOp(true), 0},

	{"Block", `;`, opToken(blockOpType), 0},
	{"Alternative", `\/\/`, opToken(alternativeOpType), 0},

	{"DocumentIndex", `documentIndex|document_?index|di`, opToken(getDocumentIndexOpType), 0},

	{"Uppercase", `upcase|ascii_?upcase`, opTokenWithPrefs(changeCaseOpType, nil, changeCasePrefs{ToUpperCase: true}), 0},
	{"Downcase", `downcase|ascii_?downcase`, opTokenWithPrefs(changeCaseOpType, nil, changeCasePrefs{ToUpperCase: false}), 0},
	simpleOp("trim", trimOpType),
	simpleOp("to_?string", toStringOpType),

	{"HexValue", `0[xX][0-9A-Fa-f]+`, hexValue(), 0},
	{"FloatValueScientific", `-?[1-9](\.\d+)?[Ee][-+]?\d+`, floatValue(), 0},
	{"FloatValue", `-?\d+(\.\d+)`, floatValue(), 0},

	{"NumberValue", `-?\d+`, numberValue(), 0},

	{"TrueBooleanValue", `[Tt][Rr][Uu][Ee]`, booleanValue(true), 0},
	{"FalseBooleanValue", `[Ff][Aa][Ll][Ss][Ee]`, booleanValue(false), 0},

	{"NullValue", `[Nn][Uu][Ll][Ll]|~`, nullValue(), 0},

	{"QuotedStringValue", `"([^"\\]*(\\.[^"\\]*)*)"`, stringValue(), 0},

	{"StrEnvOp", `strenv\([^\)]+\)`, envOp(true), 0},
	{"EnvOp", `env\([^\)]+\)`, envOp(false), 0},

	{"EnvSubstWithOptions", `envsubst\((ne|nu|ff| |,)+\)`, envSubstWithOptions(), 0},
	simpleOp("envsubst", envsubstOpType),

	{"Equals", `\s*==\s*`, opToken(equalsOpType), 0},
	{"NotEquals", `\s*!=\s*`, opToken(notEqualsOpType), 0},

	{"GreaterThanEquals", `\s*>=\s*`, opTokenWithPrefs(compareOpType, nil, compareTypePref{OrEqual: true, Greater: true}), 0},
	{"LessThanEquals", `\s*<=\s*`, opTokenWithPrefs(compareOpType, nil, compareTypePref{OrEqual: true, Greater: false}), 0},

	{"GreaterThan", `\s*>\s*`, opTokenWithPrefs(compareOpType, nil, compareTypePref{OrEqual: false, Greater: true}), 0},
	{"LessThan", `\s*<\s*`, opTokenWithPrefs(compareOpType, nil, compareTypePref{OrEqual: false, Greater: false}), 0},

	simpleOp("min", minOpType),
	simpleOp("max", maxOpType),

	{"AssignRelative", `\|=[c]*`, assignOpToken(true), 0},
	{"Assign", `=[c]*`, assignOpToken(false), 0},

	{`whitespace`, `[ \t\n]+`, nil, 0},

	{"WrappedPathElement", `\."[^ "]+"\??`, pathToken(true), 0},
	{"PathElement", `\.[^ ;\}\{\:\[\],\|\.\[\(\)=\n!]+\??`, pathToken(false), 0},
	{"Pipe", `\|`, opToken(pipeOpType), 0},
	{"Self", `\.`, opToken(selfReferenceOpType), 0},

	{"Union", `,`, opToken(unionOpType), 0},

	{"MultiplyAssign", `\*=[\+|\?cdn]*`, multiplyWithPrefs(multiplyAssignOpType), 0},
	{"Multiply", `\*[\+|\?cdn]*`, multiplyWithPrefs(multiplyOpType), 0},

	{"Divide", `\/`, opToken(divideOpType), 0},

	{"Modulo", `%`, opToken(moduloOpType), 0},

	{"AddAssign", `\+=`, opToken(addAssignOpType), 0},
	{"Add", `\+`, opToken(addOpType), 0},

	{"SubtractAssign", `\-=`, opToken(subtractAssignOpType), 0},
	{"Subtract", `\-`, opToken(subtractOpType), 0},
	{"Comment", `#.*`, nil, 0},

	simpleOp("pivot", pivotOpType),
}

type yqAction func(lexer.Token) (*token, error)

type participleYqRule struct {
	Name                string
	Pattern             string
	CreateYqToken       yqAction
	ParticipleTokenType lexer.TokenType
}

type participleLexer struct {
	lexerDefinition lexer.StringDefinition
}

func simpleOp(name string, opType *operationType) *participleYqRule {
	return &participleYqRule{strings.ToUpper(string(name[1])) + name[1:], name, opToken(opType), 0}
}

func assignableOp(name string, opType *operationType, assignOpType *operationType) *participleYqRule {
	return &participleYqRule{strings.ToUpper(string(name[1])) + name[1:], name, opTokenWithPrefs(opType, assignOpType, nil), 0}
}

func newParticipleLexer() expressionTokeniser {
	simpleRules := make([]lexer.SimpleRule, len(participleYqRules))
	for i, yqRule := range participleYqRules {
		simpleRules[i] = lexer.SimpleRule{Name: yqRule.Name, Pattern: yqRule.Pattern}
	}
	lexerDefinition := lexer.MustSimple(simpleRules)
	symbols := lexerDefinition.Symbols()

	for _, yqRule := range participleYqRules {
		yqRule.ParticipleTokenType = symbols[yqRule.Name]
	}

	return &participleLexer{lexerDefinition}
}

func pathToken(wrapped bool) yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		value := rawToken.Value
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

func recursiveDecentOpToken(includeMapKeys bool) yqAction {
	prefs := recursiveDescentPreferences{
		RecurseArray: true,
		TraversePreferences: traversePreferences{
			DontFollowAlias: true,
			IncludeMapKeys:  includeMapKeys,
		},
	}
	return opTokenWithPrefs(recursiveDescentOpType, nil, prefs)
}

func opTokenWithPrefs(opType *operationType, assignOpType *operationType, preferences interface{}) yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		value := rawToken.Value
		op := &Operation{OperationType: opType, Value: opType.Type, StringValue: value, Preferences: preferences}
		var assign *Operation
		if assignOpType != nil {
			assign = &Operation{OperationType: assignOpType, Value: assignOpType.Type, StringValue: value, Preferences: preferences}
		}
		return &token{TokenType: operationToken, Operation: op, AssignOperation: assign, CheckForPostTraverse: op.OperationType.CheckForPostTraverse}, nil
	}
}

func expressionOpToken(expression string) yqAction {
	return func(_ lexer.Token) (*token, error) {
		prefs := expressionOpPreferences{expression: expression}
		expressionOp := &Operation{OperationType: expressionOpType, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: expressionOp}, nil
	}
}

func flattenWithDepth() yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		value := rawToken.Value
		var depth, errParsingInt = extractNumberParameter(value)
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		prefs := flattenPreferences{depth: depth}
		op := &Operation{OperationType: flattenOpType, Value: flattenOpType.Type, StringValue: value, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op, CheckForPostTraverse: flattenOpType.CheckForPostTraverse}, nil
	}
}

func assignAllCommentsOp(updateAssign bool) yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		log.Debug("assignAllCommentsOp %v", rawToken.Value)
		value := rawToken.Value
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

func assignOpToken(updateAssign bool) yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		log.Debug("assignOpToken %v", rawToken.Value)
		value := rawToken.Value
		prefs := assignPreferences{DontOverWriteAnchor: true}
		if strings.Contains(value, "c") {
			prefs.ClobberCustomTags = true
		}
		op := &Operation{OperationType: assignOpType, Value: assignOpType.Type, StringValue: value, UpdateAssign: updateAssign, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op}, nil
	}
}

func booleanValue(val bool) yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		return &token{TokenType: operationToken, Operation: createValueOperation(val, rawToken.Value)}, nil
	}
}

func nullValue() yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		return &token{TokenType: operationToken, Operation: createValueOperation(nil, rawToken.Value)}, nil
	}
}

func stringValue() yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		log.Debug("rawTokenvalue: %v", rawToken.Value)
		value := unwrap(rawToken.Value)
		log.Debug("unwrapped: %v", value)
		value = processEscapeCharacters(value)
		return &token{TokenType: operationToken, Operation: &Operation{
			OperationType: stringInterpolationOpType,
			StringValue:   value,
			Value:         value,
		}}, nil
	}
}

func envOp(strenv bool) yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		value := rawToken.Value
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

		return &token{TokenType: operationToken, Operation: envOperation, CheckForPostTraverse: envOpType.CheckForPostTraverse}, nil
	}
}

func envSubstWithOptions() yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		value := rawToken.Value
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

func multiplyWithPrefs(op *operationType) yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		prefs := multiplyPreferences{}
		prefs.AssignPrefs = assignPreferences{}
		options := rawToken.Value
		if strings.Contains(options, "+") {
			prefs.AppendArrays = true
		}
		if strings.Contains(options, "?") {
			prefs.TraversePrefs = traversePreferences{DontAutoCreate: true}
		}
		if strings.Contains(options, "n") {
			prefs.AssignPrefs.OnlyWriteNull = true
		}
		if strings.Contains(options, "d") {
			prefs.DeepMergeArrays = true
		}
		if strings.Contains(options, "c") {
			prefs.AssignPrefs.ClobberCustomTags = true
		}
		prefs.TraversePrefs.DontFollowAlias = true
		prefs.TraversePrefs.ExactKeyMatch = true
		op := &Operation{OperationType: op, Value: multiplyOpType.Type, StringValue: options, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op}, nil
	}

}

func getVariableOpToken() yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		value := rawToken.Value

		value = value[1:]

		getVarOperation := createValueOperation(value, value)
		getVarOperation.OperationType = getVariableOpType

		return &token{TokenType: operationToken, Operation: getVarOperation, CheckForPostTraverse: true}, nil
	}
}

func hexValue() yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		var originalString = rawToken.Value
		var numberString = originalString[2:]
		log.Debugf("numberString: %v", numberString)
		var number, errParsingInt = strconv.ParseInt(numberString, 16, 64)
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		return &token{TokenType: operationToken, Operation: createValueOperation(number, originalString)}, nil
	}
}

func floatValue() yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		var numberString = rawToken.Value
		var number, errParsingInt = strconv.ParseFloat(numberString, 64)
		if errParsingInt != nil {
			return nil, errParsingInt
		}
		return &token{TokenType: operationToken, Operation: createValueOperation(number, numberString)}, nil
	}
}

func numberValue() yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		var numberString = rawToken.Value
		var number, errParsingInt = strconv.ParseInt(numberString, 10, 64)
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		return &token{TokenType: operationToken, Operation: createValueOperation(number, numberString)}, nil
	}
}

func parentWithLevel() yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		value := rawToken.Value
		var level, errParsingInt = extractNumberParameter(value)
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		prefs := parentOpPreferences{Level: level}
		op := &Operation{OperationType: getParentOpType, Value: getParentOpType.Type, StringValue: value, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op, CheckForPostTraverse: true}, nil
	}
}

func parentWithDefaultLevel() yqAction {
	return func(_ lexer.Token) (*token, error) {
		prefs := parentOpPreferences{Level: 1}
		op := &Operation{OperationType: getParentOpType, Value: getParentOpType.Type, StringValue: getParentOpType.Type, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op, CheckForPostTraverse: true}, nil
	}
}

func encodeParseIndent(outputFormat *Format) yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		value := rawToken.Value
		var indent, errParsingInt = extractNumberParameter(value)
		if errParsingInt != nil {
			return nil, errParsingInt
		}

		prefs := encoderPreferences{format: outputFormat, indent: indent}
		op := &Operation{OperationType: encodeOpType, Value: encodeOpType.Type, StringValue: value, Preferences: prefs}
		return &token{TokenType: operationToken, Operation: op}, nil
	}
}

func encodeWithIndent(outputFormat *Format, indent int) yqAction {
	prefs := encoderPreferences{format: outputFormat, indent: indent}
	return opTokenWithPrefs(encodeOpType, nil, prefs)
}

func decodeOp(format *Format) yqAction {
	prefs := decoderPreferences{format: format}
	return opTokenWithPrefs(decodeOpType, nil, prefs)
}

func loadOp(decoder Decoder) yqAction {
	prefs := loadPrefs{decoder}
	return opTokenWithPrefs(loadOpType, nil, prefs)
}

func opToken(op *operationType) yqAction {
	return opTokenWithPrefs(op, nil, nil)
}

func literalToken(tt tokenType, checkForPost bool) yqAction {
	return func(rawToken lexer.Token) (*token, error) {
		return &token{TokenType: tt, CheckForPostTraverse: checkForPost, Match: rawToken.Value}, nil
	}
}

func (p *participleLexer) getYqDefinition(rawToken lexer.Token) *participleYqRule {
	for _, yqRule := range participleYqRules {
		if yqRule.ParticipleTokenType == rawToken.Type {
			return yqRule
		}
	}
	return &participleYqRule{}
}

func (p *participleLexer) Tokenise(expression string) ([]*token, error) {
	myLexer, err := p.lexerDefinition.LexString("", expression)
	if err != nil {
		return nil, err
	}
	tokens := make([]*token, 0)

	for {
		rawToken, e := myLexer.Next()
		if e != nil {
			return nil, e
		} else if rawToken.Type == lexer.EOF {
			return postProcessTokens(tokens), nil
		}

		definition := p.getYqDefinition(rawToken)
		if definition.CreateYqToken != nil {
			token, e := definition.CreateYqToken(rawToken)
			if e != nil {
				return nil, e
			}
			tokens = append(tokens, token)
		}

	}

}
