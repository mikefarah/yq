package yqlib

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/mikefarah/yq/v4/test"
)

type participleLexerScenario struct {
	expression string
	tokens     []*token
}

var participleLexerScenarios = []participleLexerScenario{
	{
		expression: "to_entries[]",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: toEntriesOpType,
					Value:         "TO_ENTRIES",
					StringValue:   "to_entries",
				},
				CheckForPostTraverse: true,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traverseArrayOpType,
				},
			},
			{
				TokenType: openCollect,
				Match:     "[",
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: emptyOpType,
					StringValue:   "EMPTY",
				},
			},
			{
				TokenType:            closeCollect,
				CheckForPostTraverse: true,
				Match:                "]",
			},
		},
	},
	{
		expression: ".a!=",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "a",
					StringValue:   "a",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: notEqualsOpType,
					Value:         "NOT_EQUALS",
					StringValue:   "!=",
				},
			},
		},
	},
	{
		expression: ".[:3]",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: selfReferenceOpType,
					StringValue:   "SELF",
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traverseArrayOpType,
					StringValue:   "TRAVERSE_ARRAY",
				},
			},
			{
				TokenType: openCollect,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: valueOpType,
					Value:         0,
					StringValue:   "0",
					CandidateNode: &CandidateNode{
						Kind:  ScalarNode,
						Tag:   "!!int",
						Value: "0",
					},
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: createMapOpType,
					Value:         "CREATE_MAP",
					StringValue:   ":",
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: valueOpType,
					Value:         int64(3),
					StringValue:   "3",
					CandidateNode: &CandidateNode{
						Kind:  ScalarNode,
						Tag:   "!!int",
						Value: "3",
					},
				},
			},
			{
				TokenType:            closeCollect,
				CheckForPostTraverse: true,
				Match:                "]",
			},
		},
	},
	{
		expression: ".[-2:]",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: selfReferenceOpType,
					StringValue:   "SELF",
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traverseArrayOpType,
					StringValue:   "TRAVERSE_ARRAY",
				},
			},
			{
				TokenType: openCollect,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: valueOpType,
					Value:         int64(-2),
					StringValue:   "-2",
					CandidateNode: &CandidateNode{
						Kind:  ScalarNode,
						Tag:   "!!int",
						Value: "-2",
					},
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: createMapOpType,
					Value:         "CREATE_MAP",
					StringValue:   ":",
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: lengthOpType,
				},
			},
			{
				TokenType:            closeCollect,
				CheckForPostTraverse: true,
				Match:                "]",
			},
		},
	},
	{
		expression: ".a",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "a",
					StringValue:   "a",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
		},
	},
	{
		expression: ".a.b",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "a",
					StringValue:   "a",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: shortPipeOpType,
					Value:         "PIPE",
					StringValue:   ".",
					Preferences:   nil,
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "b",
					StringValue:   "b",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
		},
	},
	{
		expression: ".a.b?",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "a",
					StringValue:   "a",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: shortPipeOpType,
					Value:         "PIPE",
					StringValue:   ".",
					Preferences:   nil,
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "b",
					StringValue:   "b",
					Preferences: traversePreferences{
						OptionalTraverse: true,
					},
				},
				CheckForPostTraverse: true,
			},
		},
	},
	{
		expression: `.a."b?"`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "a",
					StringValue:   "a",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: shortPipeOpType,
					Value:         "PIPE",
					StringValue:   ".",
					Preferences:   nil,
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "b?",
					StringValue:   "b?",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
		},
	},
	{
		expression: `   .a  ."b?"`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "a",
					StringValue:   "a",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: shortPipeOpType,
					Value:         "PIPE",
					StringValue:   ".",
					Preferences:   nil,
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "b?",
					StringValue:   "b?",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
		},
	},
	{
		expression: `.a | .b`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "a",
					StringValue:   "a",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: pipeOpType,
					Value:         "PIPE",
					StringValue:   "|",
					Preferences:   nil,
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "b",
					StringValue:   "b",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
		},
	},
	{
		expression: "(.a)",
		tokens: []*token{
			{
				TokenType: openBracket,
				Match:     "(",
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "a",
					StringValue:   "a",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
			{
				TokenType:            closeBracket,
				Match:                ")",
				CheckForPostTraverse: true,
			},
		},
	},
	{
		expression: "..",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: recursiveDescentOpType,
					Value:         "RECURSIVE_DESCENT",
					StringValue:   "..",
					Preferences: recursiveDescentPreferences{
						RecurseArray: true,
						TraversePreferences: traversePreferences{
							DontFollowAlias: true,
							IncludeMapKeys:  false,
						},
					},
				},
			},
		},
	},
	{
		expression: "...",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: recursiveDescentOpType,
					Value:         "RECURSIVE_DESCENT",
					StringValue:   "...",
					Preferences: recursiveDescentPreferences{
						RecurseArray: true,
						TraversePreferences: traversePreferences{
							DontFollowAlias: true,
							IncludeMapKeys:  true,
						},
					},
				},
			},
		},
	},
	{
		expression: ".a,.b",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "a",
					StringValue:   "a",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: unionOpType,
					Value:         "UNION",
					StringValue:   ",",
					Preferences:   nil,
				},
			},
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: traversePathOpType,
					Value:         "b",
					StringValue:   "b",
					Preferences:   traversePreferences{},
				},
				CheckForPostTraverse: true,
			},
		},
	},
	{
		expression: "map_values",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: mapValuesOpType,
					Value:         "MAP_VALUES",
					StringValue:   "map_values",
					Preferences:   nil,
				},
				CheckForPostTraverse: mapValuesOpType.CheckForPostTraverse,
			},
		},
	},
	{
		expression: "mapvalues",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: mapValuesOpType,
					Value:         "MAP_VALUES",
					StringValue:   "mapvalues",
					Preferences:   nil,
				},
				CheckForPostTraverse: mapValuesOpType.CheckForPostTraverse,
			},
		},
	},
	{
		expression: "flatten(3)",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: flattenOpType,
					Value:         "FLATTEN_BY",
					StringValue:   "flatten(3)",
					Preferences:   flattenPreferences{depth: 3},
				},
				CheckForPostTraverse: flattenOpType.CheckForPostTraverse,
			},
		},
	},
	{
		expression: "flatten",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: flattenOpType,
					Value:         "FLATTEN_BY",
					StringValue:   "flatten",
					Preferences:   flattenPreferences{depth: -1},
				},
				CheckForPostTraverse: flattenOpType.CheckForPostTraverse,
			},
		},
	},
	{
		expression: "length",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: lengthOpType,
					Value:         "LENGTH",
					StringValue:   "length",
					Preferences:   nil,
				},
			},
		},
	},
	{
		expression: "format_datetime",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: formatDateTimeOpType,
					Value:         "FORMAT_DATE_TIME",
					StringValue:   "format_datetime",
					Preferences:   nil,
				},
			},
		},
	},
	{
		expression: "to_yaml(3)",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: encodeOpType,
					Value:         "ENCODE",
					StringValue:   "to_yaml(3)",
					Preferences: encoderPreferences{
						format: YamlFormat,
						indent: 3,
					},
				},
			},
		},
	},
	{
		expression: "tojson(2)",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: encodeOpType,
					Value:         "ENCODE",
					StringValue:   "tojson(2)",
					Preferences: encoderPreferences{
						format: JSONFormat,
						indent: 2,
					},
				},
			},
		},
	},
	{
		expression: "@yaml",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: encodeOpType,
					Value:         "ENCODE",
					StringValue:   "@yaml",
					Preferences: encoderPreferences{
						format: YamlFormat,
						indent: 2,
					},
				},
			},
		},
	},
	{
		expression: "to_props",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: encodeOpType,
					Value:         "ENCODE",
					StringValue:   "to_props",
					Preferences: encoderPreferences{
						format: PropertiesFormat,
						indent: 2,
					},
				},
			},
		},
	},
	{
		expression: "@base64d",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: decodeOpType,
					Value:         "DECODE",
					StringValue:   "@base64d",
					Preferences: decoderPreferences{
						format: Base64Format,
					},
				},
			},
		},
	},
	{
		expression: "@base64",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: encodeOpType,
					Value:         "ENCODE",
					StringValue:   "@base64",
					Preferences: encoderPreferences{
						format: Base64Format,
					},
				},
			},
		},
	},
	{
		expression: "@yamld",
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: decodeOpType,
					Value:         "DECODE",
					StringValue:   "@yamld",
					Preferences: decoderPreferences{
						format: YamlFormat,
					},
				},
			},
		},
	},
	{
		expression: `"string with a\n"`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: stringInterpolationOpType,
					Value:         "string with a\n",
					StringValue:   "string with a\n",
					Preferences:   nil,
				},
			},
		},
	},
	{
		expression: `"string with a \""`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: stringInterpolationOpType,
					Value:         `string with a "`,
					StringValue:   `string with a "`,
					Preferences:   nil,
				},
			},
		},
	},
	{
		expression: `"string with a\r"`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: stringInterpolationOpType,
					Value:         "string with a\r",
					StringValue:   "string with a\r",
					Preferences:   nil,
				},
			},
		},
	},
	{
		expression: `"string with a\t"`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: stringInterpolationOpType,
					Value:         "string with a\t",
					StringValue:   "string with a\t",
					Preferences:   nil,
				},
			},
		},
	},
	{
		expression: `"string with a\f"`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: stringInterpolationOpType,
					Value:         "string with a\f",
					StringValue:   "string with a\f",
					Preferences:   nil,
				},
			},
		},
	},
	{
		expression: `"string with a\v"`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: stringInterpolationOpType,
					Value:         "string with a\v",
					StringValue:   "string with a\v",
					Preferences:   nil,
				},
			},
		},
	},
	{
		expression: `"string with a\b"`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: stringInterpolationOpType,
					Value:         "string with a\b",
					StringValue:   "string with a\b",
					Preferences:   nil,
				},
			},
		},
	},
	{
		expression: `"string with a\a"`,
		tokens: []*token{
			{
				TokenType: operationToken,
				Operation: &Operation{
					OperationType: stringInterpolationOpType,
					Value:         "string with a\a",
					StringValue:   "string with a\a",
					Preferences:   nil,
				},
			},
		},
	},
}

func TestParticipleLexer(t *testing.T) {
	lexer := newParticipleLexer()

	for _, scenario := range participleLexerScenarios {
		actual, err := lexer.Tokenise(scenario.expression)
		if err != nil {
			t.Error(err)
		} else {
			test.AssertResultWithContext(t, repr.String(scenario.tokens, repr.Indent(" ")), repr.String(actual, repr.Indent(" ")), scenario.expression)
		}

	}
}
