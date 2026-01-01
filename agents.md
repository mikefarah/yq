# General rules
✅ **DO:**
- You can use ./yq with the `--debug-node-info` flag to get a deeper understanding of the ast.
- run ./scripts/format.sh to format the code; then ./scripts/check.sh lint and finally ./scripts/spelling.sh to check spelling.
- Add comprehensive tests to cover the changes
- Run test suite to ensure there is no regression
- Use UK english spelling

❌ **DON'T:**
- Git add or commit
- Add comments to functions that are self-explanatory



# Adding a New Encoder/Decoder

This guide explains how to add support for a new format (encoder/decoder) to yq without modifying `candidate_node.go`.

## Overview

The encoder/decoder architecture in yq is based on two main interfaces:

- **Encoder**: Converts a `CandidateNode` to output in a specific format
- **Decoder**: Reads input in a specific format and creates a `CandidateNode`

Each format is registered in `pkg/yqlib/format.go` and made available through factory functions.

## Architecture

### Key Files

- `pkg/yqlib/encoder.go` - Defines the `Encoder` interface
- `pkg/yqlib/decoder.go` - Defines the `Decoder` interface
- `pkg/yqlib/format.go` - Format registry and factory functions
- `pkg/yqlib/operator_encoder_decoder.go` - Encode/decode operators
- `pkg/yqlib/encoder_*.go` - Encoder implementations
- `pkg/yqlib/decoder_*.go` - Decoder implementations

### Interfaces

**Encoder Interface:**
```go
type Encoder interface {
    Encode(writer io.Writer, node *CandidateNode) error
    PrintDocumentSeparator(writer io.Writer) error
    PrintLeadingContent(writer io.Writer, content string) error
    CanHandleAliases() bool
}
```

**Decoder Interface:**
```go
type Decoder interface {
    Init(reader io.Reader) error
    Decode() (*CandidateNode, error)
}
```

## Step-by-Step: Adding a New Encoder/Decoder

### Step 1: Create the Encoder File

Create `pkg/yqlib/encoder_<format>.go` implementing the `Encoder` interface:
- `Encode()` - Convert a `CandidateNode` to your format and write to the output writer
- `PrintDocumentSeparator()` - Handle document separators if your format requires them
- `PrintLeadingContent()` - Handle leading content/comments if supported
- `CanHandleAliases()` - Return whether your format supports YAML aliases

See `encoder_json.go` or `encoder_base64.go` for examples.

### Step 2: Create the Decoder File

Create `pkg/yqlib/decoder_<format>.go` implementing the `Decoder` interface:
- `Init()` - Initialize the decoder with the input reader and set up any needed state
- `Decode()` - Decode one document from the input and return a `CandidateNode`, or `io.EOF` when finished

See `decoder_json.go` or `decoder_base64.go` for examples.

### Step 3: Create Tests (Mandatory)

Create a test file `pkg/yqlib/<format>_test.go` using the `formatScenario` pattern:
- Define test scenarios as `formatScenario` structs with fields: `description`, `input`, `expected`, `scenarioType`
- `scenarioType` can be `"decode"` (test decoding to YAML) or `"roundtrip"` (encode/decode preservation)
- Create a helper function `test<Format>Scenario()` that switches on `scenarioType`
- Create main test function `Test<Format>FormatScenarios()` that iterates over scenarios
- The main test function should use `documentScenarios` to ensure testcase documentation is generated.

Test coverage must include:
- Basic data types (scalars, arrays, objects/maps)
- Nested structures
- Edge cases (empty inputs, special characters, escape sequences)
- Format-specific features or syntax
- Round-trip tests: decode → encode → decode should preserve data

See `hcl_test.go` for a complete example.

### Step 4: Register the Format in format.go

Edit `pkg/yqlib/format.go`:

1. Add a new format variable:
   - `"<format>"` is the formal name (e.g., "json", "yaml")
   - `[]string{...}` contains short aliases (can be empty)
   - The first function creates an encoder (can be nil for encode-only formats)
   - The second function creates a decoder (can be nil for decode-only formats)

2. Add the format to the `Formats` slice in the same file

See existing formats in `format.go` for the exact structure.

### Step 5: Handle Encoder Configuration (if needed)

If your format has preferences/configuration options:

1. Create a preferences struct with your configuration fields
2. Update the encoder to accept preferences in its factory function
3. Update `format.go` to pass the configured preferences
4. Update `operator_encoder_decoder.go` if special indent handling is needed (see existing formats like JSON and YAML for the pattern)

This pattern is optional and only needed if your format has user-configurable options.

## Build Tags

Use build tags to allow optional compilation of formats:
- Add `//go:build !yq_no<format>` at the top of your encoder and decoder files
- Create a no-build version in `pkg/yqlib/no_<format>.go` that returns nil for encoder/decoder factories

This allows users to compile yq without certain formats using: `go build -tags yq_no<format>`

## Working with CandidateNode

The `CandidateNode` struct represents a YAML node with:
- `Kind`: The node type (ScalarNode, SequenceNode, MappingNode)
- `Tag`: The YAML tag (e.g., "!!str", "!!int", "!!map")
- `Value`: The scalar value (for ScalarNode only)
- `Content`: Child nodes (for SequenceNode and MappingNode)

Key methods:
- `node.guessTagFromCustomType()` - Infer the tag from Go type
- `node.AsList()` - Convert to a list for processing
- `node.CreateReplacement()` - Create a new replacement node
- `NewCandidate()` - Create a new CandidateNode

## Key Points

✅ **DO:**
- Implement only the `Encoder` and `Decoder` interfaces
- Register your format in `format.go` only
- Keep format-specific logic in your encoder/decoder files
- Use the candidate_node style attribute to store style information for round-trip. Ask if this needs to be updated with new styles.
- Use build tags for optional compilation
- Add comprehensive tests
- Run the specific encoder/decoder test (e.g. <format>_test.go) whenever you make ay changes to the encoder_<format> or decoder_<format>
- Handle errors gracefully
- Add the no build directive, like the xml encoder and decoder, that enables a minimal yq builds. e.g.  `//go:build !yq_<format>`. Be sure to also update the build_small-yq.sh and build-tinygo-yq.sh to not include the new format.

❌ **DON'T:**
- Modify `candidate_node.go` to add format-specific logic
- Add format-specific fields to `CandidateNode`
- Create special cases in core navigation or evaluation logic
- Bypass the encoder/decoder interfaces
- Use candidate_node tag attribute for anything other than indicate the data type

## Examples

Refer to existing format implementations for patterns:

- **Simple encoder/decoder**: `encoder_json.go`, `decoder_json.go`
- **Complex with preferences**: `encoder_yaml.go`, `decoder_yaml.go`
- **Encoder-only**: `encoder_sh.go` (ShFormat has nil decoder)
- **String-only operations**: `encoder_base64.go`, `decoder_base64.go`

## Testing Your Implementation (Mandatory)

Tests must be implemented in `<format>_test.go` following the `formatScenario` pattern:

1. **Create test scenarios** using the `formatScenario` struct with fields:
   - `description`: Brief description of what's being tested
   - `input`: Sample input in your format
   - `expected`: Expected output (typically in YAML for decode tests)
   - `scenarioType`: Either `"decode"` or `"roundtrip"`

2. **Test coverage must include:**
   - Basic data types (scalars, arrays, objects/maps)
   - Nested structures
   - Edge cases (empty inputs, special characters, escape sequences)
   - Format-specific features or syntax
   - Round-trip tests: decode → encode → decode should preserve data

3. **Test function pattern:**
   - `test<Format>Scenario()`: Helper function that switches on `scenarioType`
   - `Test<Format>FormatScenarios()`: Main test function that iterates over scenarios

4. **Example from existing formats:**
   - See `hcl_test.go` for a complete example
   - See `yaml_test.go` for YAML-specific patterns
   - See `json_test.go` for more complex scenarios

## Common Patterns

### Format with Indentation
Use preferences to control output formatting:
```go
type <format>Preferences struct {
    Indent int
}

func (prefs *<format>Preferences) Copy() <format>Preferences {
    return *prefs
}
```

### Multiple Documents
Decoders should support reading multiple documents:
```go
func (dec *<format>Decoder) Decode() (*CandidateNode, error) {
    if dec.finished {
        return nil, io.EOF
    }
    // ... decode next document ...
    if noMoreDocuments {
        dec.finished = true
    }
    return candidate, nil
}
```

---

# Adding a New Operator

This guide explains how to add a new operator to yq. Operators are the core of yq's expression language and process `CandidateNode` objects without requiring modifications to `candidate_node.go` itself.

## Overview

Operators transform data by implementing a handler function that processes a `Context` containing `CandidateNode` objects. Each operator is:

1. Defined as an `operationType` in `operation.go`
2. Registered in the lexer in `lexer_participle.go`
3. Implemented in its own `operator_<type>.go` file
4. Tested in `operator_<type>_test.go`
5. Documented in `pkg/yqlib/doc/operators/headers/<type>.md`

## Architecture

### Key Files

- `pkg/yqlib/operation.go` - Defines `operationType` and operator registry
- `pkg/yqlib/lexer_participle.go` - Registers operators with their syntax patterns
- `pkg/yqlib/operator_<type>.go` - Operator implementation
- `pkg/yqlib/operator_<type>_test.go` - Operator tests using `expressionScenario`
- `pkg/yqlib/doc/operators/headers/<type>.md` - Documentation header

### Core Types

**operationType:**
```go
type operationType struct {
    Type                 string          // Unique operator name (e.g., "REVERSE")
    NumArgs              uint            // Number of arguments (0 for no args)
    Precedence           uint            // Operator precedence (higher = higher precedence)
    Handler              operatorHandler // The function that executes the operator
    CheckForPostTraverse bool            // Whether to apply post-traversal logic
    ToString             func(*Operation) string // Custom string representation
}
```

**operatorHandler signature:**
```go
type operatorHandler func(*dataTreeNavigator, Context, *ExpressionNode) (Context, error)
```

**expressionScenario for tests:**
```go
type expressionScenario struct {
    description      string
    subdescription   string
    document         string
    expression       string
    expected         []string
    skipDoc          bool
    expectedError    string
}
```

## Step-by-Step: Adding a New Operator

### Step 1: Create the Operator Implementation File

Create `pkg/yqlib/operator_<type>.go` implementing the operator handler function:
- Implement the `operatorHandler` function signature
- Process nodes from `context.MatchingNodes`
- Return a new `Context` with results using `context.ChildContext()`
- Use `candidate.CreateReplacement()` or `candidate.CreateReplacementWithComments()` to create new nodes
- Handle errors gracefully with meaningful error messages

See `operator_reverse.go` or `operator_keys.go` for examples.

### Step 2: Register the Operator in operation.go

Add the operator type definition to `pkg/yqlib/operation.go`:

```go
var <type>OpType = &operationType{
    Type:       "<TYPE>",          // All caps, matches pattern in lexer
    NumArgs:    0,                 // 0 for no args, 1+ for args
    Precedence: 50,                // Typical range: 40-55
    Handler:    <type>Operator,    // Reference to handler function
}
```

**Precedence guidelines:**
- 10-20: Logical operators (OR, AND, UNION)
- 30: Pipe operator
- 40: Assignment and comparison operators
- 42: Arithmetic operators (ADD, SUBTRACT, MULTIPLY, DIVIDE)
- 50-52: Most other operators
- 55: High precedence (e.g., GET_VARIABLE)

**Optional fields:**
- `CheckForPostTraverse: true` - If your operator can have another directly after it without the pipe character. Most of the time this is false.
- `ToString: customToString` - Custom string representation (rarely needed)

### Step 3: Register the Operator in lexer_participle.go

Edit `pkg/yqlib/lexer_participle.go` to add the operator to the lexer rules:
- Use `simpleOp()` for simple keyword patterns
- Use object syntax for regex patterns or complex syntax
- Support optional characters with `_?` and aliases with `|`

See existing operators in `lexer_participle.go` for pattern examples.

### Step 4: Create Tests (Mandatory)

Create `pkg/yqlib/operator_<type>_test.go` using the `expressionScenario` pattern:
- Define test scenarios with `description`, `document`, `expression`, and `expected` fields
- `expected` is a slice of strings showing output format: `"D<doc>, P[<path>], (<tag>)::<value>\n"`
- Set `skipDoc: true` for edge cases you don't want in generated documentation
- Include `subdescription` for longer test names
- Set `expectedError` if testing error cases
- Create main test function that iterates over scenarios
- The main test function should use `documentScenarios` to ensure testcase documentation is generated.

Test coverage must include:
- Basic data types and nested structures
- Edge cases (empty inputs, special characters, type errors)
- Multiple outputs if applicable
- Format-specific features

See `operator_reverse_test.go` for a simple example and `operator_keys_test.go` for complex cases.

### Step 5: Create Documentation Header

Create `pkg/yqlib/doc/operators/headers/<type>.md`:
- Use the exact operator name as the title
- Include a concise 1-2 sentence summary
- Add additional context or examples if the operator is complex

See existing headers in `doc/operators/headers/` for examples.

## Working with Context and CandidateNode

### Context Management
- `context.ChildContext(results)` - Create child context with results
- `context.GetVariable("varName")` - Get variables stored in context
- `context.SetVariable("varName", value)` - Set variables in context

### CandidateNode Operations
- `candidate.CreateReplacement(ScalarNode, "!!str", stringValue)` - Create a replacement node
- `candidate.CreateReplacementWithComments(SequenceNode, "!!seq", candidate.Style)` - With style preserved
- `candidate.Kind` - The node type (ScalarNode, SequenceNode, MappingNode)
- `candidate.Tag` - The YAML tag (!!str, !!int, etc.)
- `candidate.Value` - The scalar value (for ScalarNode only)
- `candidate.Content` - Child nodes (for SequenceNode and MappingNode)
- `candidate.guessTagFromCustomType()` - Infer the tag from Go type
- `candidate.AsList()` - Convert to a list representation

## Key Points

✅ **DO:**
- Implement the operator handler with the correct signature
- Register in `operation.go` with appropriate precedence
- Add the lexer pattern in `lexer_participle.go`
- Write comprehensive tests covering normal and edge cases
- Create a documentation header in `doc/operators/headers/`
- Use `Context.ChildContext()` for proper context threading
- Handle all node types gracefully
- Return meaningful error messages

❌ **DON'T:**
- Modify `candidate_node.go` (operators shouldn't need this)
- Modify core navigation or evaluation logic
- Bypass the handler function pattern
- Add format-specific or operator-specific fields to `CandidateNode`
- Skip tests or documentation

## Examples

Refer to existing operator implementations for patterns:

- **No-argument operator**: `operator_reverse.go` - Processes arrays/sequences
- **Single-argument operator**: `operator_map.go` - Takes an expression argument
- **Complex multi-output**: `operator_keys.go` - Produces multiple results
- **With preferences**: `operator_to_number.go` - Configuration options
- **Error handling**: `operator_error.go` - Control flow with errors
- **String operations**: `operator_strings.go` - Multiple related operators

## Testing Patterns

Refer to existing test files for specific patterns:
- Basic expression tests in `operator_reverse_test.go`
- Multi-output tests in `operator_keys_test.go`
- Error handling tests in `operator_error_test.go`
- Tests with `skipDoc` flag to exclude from generated documentation

## Common Patterns

Refer to existing operator implementations for these patterns:
- Simple transformation: see `operator_reverse.go`
- Type checking: see `operator_error.go`
- Working with arguments: see `operator_map.go`
- Post-traversal operators: see `operator_with.go`
