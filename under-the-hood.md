# yq Under the hood

## 1. High-level overview

- Every yq run turns a text expression into an `*ExpressionNode` tree exactly once, then reuses that tree against every input document.
- Documents are decoded into `*CandidateNode` trees, which is the same node type the expression tree operates on.
- `DataTreeNavigator.GetMatchingNodes` walks the expression tree against a `Context` of candidate nodes, producing a new `Context` of matches which the encoder then serialises.

```mermaid
flowchart TD
    A["yq expression string<br/>e.g. .a.b[] | select(.x > 1)"] --> B["Expression Parsing<br/>(see Diagram 2)<br/>📄 expression_parser.go:ParseExpression"]
    B --> C["Expression Tree<br/>(*ExpressionNode)<br/>📄 expression_parser.go:createExpressionTree"]
    C --> D["Load YAML/JSON/etc document(s)<br/>via Decoder → *CandidateNode<br/>📄 decoder_yaml.go:Decode"]
    D --> E["DataTreeNavigator.GetMatchingNodes<br/>(see Diagram 3)<br/>📄 data_tree_navigator.go:GetMatchingNodes"]
    E --> F["Resulting Context<br/>(list of matching *CandidateNode)"]
    F --> G["Encoder writes result<br/>(YAML/JSON/etc) via Printer<br/>📄 encoder_yaml.go:Encode"]
```

## 2. Expression parsing sub-diagram (tokenising → RPN → tree)

- **Tokenise**: a participle-generated lexer converts the raw string into an ordered list of typed tokens (operators, brackets, path segments, literals).
- **Postfix (shunting-yard)**: tokens are reordered into Reverse Polish Notation using an operator stack, so operator precedence and bracket nesting are resolved before tree-building.
- **Tree build**: the postfix `Operation` list is walked once, using a node stack, popping 1 or 2 operands per operator to build a binary `*ExpressionNode` tree (unary ops just leave `LHS` nil).

Implemented in `pkg/yqlib/expression_parser.go`, `pkg/yqlib/lexer_participle.go`, `pkg/yqlib/expression_postfix.go`.

```mermaid
flowchart TD
    subgraph Tokenise ["1. Tokenise — pkg/yqlib/lexer_participle.go"]
        A1["Raw expression string"] --> A2["participle lexer<br/>matches regex rules:<br/>operators, brackets,<br/>traverse paths, strings, numbers<br/>📄 lexer_participle.go:newParticipleLexer / Tokenise"]
        A2 --> A3["Infix token stream<br/>[]*token<br/>e.g. TRAVERSE(a) TRAVERSE(b)<br/>PIPE SELECT( TRAVERSE(x) GT 1 )<br/>📄 lexer_participle.go:token struct"]
    end

    subgraph Postfix ["2. Convert to Postfix / RPN — pkg/yqlib/expression_postfix.go"]
        B1["Infix tokens"] --> B2["Shunting-yard algorithm:<br/>walk tokens, maintain operator stack<br/>+ output result list<br/>📄 expression_postfix.go:ConvertToPostfix"]
        B2 --> B3{"Token type?"}
        B3 -->|"open bracket / collect"| B4["push onto opStack<br/>📄 expression_postfix.go:ConvertToPostfix (openBracket case)"]
        B3 -->|"close bracket / collect"| B5["pop opStack → result<br/>until matching opener found,<br/>emit COLLECT / SHORT_PIPE ops<br/>📄 expression_postfix.go:popOpToResult"]
        B3 -->|"operator"| B6["pop higher-precedence ops<br/>from opStack → result,<br/>then push current op<br/>📄 expression_postfix.go:ConvertToPostfix (default case)"]
        B4 --> B7
        B5 --> B7
        B6 --> B7["continue to next token"]
        B7 --> B2
        B2 --> B8["[]*Operation in postfix/RPN order"]
    end

    subgraph Tree ["3. Build Expression Tree — pkg/yqlib/expression_parser.go"]
        C1["Postfix []*Operation"] --> C2["walk operations left→right,<br/>use a node stack<br/>📄 expression_parser.go:createExpressionTree"]
        C2 --> C3{"Operation.NumArgs"}
        C3 -->|"0 (e.g. SELF, value)"| C4["push leaf ExpressionNode"]
        C3 -->|"1 (e.g. NOT, LENGTH)"| C5["pop 1 node as RHS,<br/>push new node"]
        C3 -->|"2 (e.g. PIPE, ADD, EQUALS)"| C6["pop 2 nodes as LHS/RHS,<br/>push new node"]
        C4 --> C7["single remaining stack item"]
        C5 --> C7
        C6 --> C7
        C7 --> C8["Root *ExpressionNode<br/>(binary tree of Operation + LHS/RHS)<br/>📄 expression_parser.go:ExpressionNode struct"]
    end

    A3 --> B1
    B8 --> C1
```

Key structs:
- `token` — raw lexed unit (bracket, operator, literal, path segment). Defined in [pkg/yqlib/lexer_participle.go](pkg/yqlib/lexer_participle.go).
- `Operation` — an `operationType` (e.g. `traverseOpType`, `pipeOpType`, `selectOpType`) plus preferences/value. Defined in [pkg/yqlib/operation.go](pkg/yqlib/operation.go).
- `ExpressionNode{Operation, LHS, RHS, Parent}` — the final AST, always binary (unary ops just leave LHS nil). Defined in [pkg/yqlib/expression_parser.go](pkg/yqlib/expression_parser.go).

## 3. Applying the expression tree to a YAML document

- Documents are decoded into `*CandidateNode`s and seeded into a `Context` (a `list.List` of matching nodes) — the starting point for evaluation.
- `GetMatchingNodes` is called recursively: each `ExpressionNode` looks up its `Operation`'s handler function and invokes it with the current `Context`.
- Handlers differ by operator: pipes chain two sub-evaluations, traversal looks up child keys/indices, `select` filters nodes, binary ops cross/zip and compare — each returning a new `Context` that feeds the next step.

Implemented in `pkg/yqlib/data_tree_navigator.go` and the various `operator_*.go` handler files, driven by `pkg/yqlib/stream_evaluator.go` or `pkg/yqlib/all_at_once_evaluator.go`.

```mermaid
flowchart TD
    A["Decoder reads YAML doc(s)<br/>→ *CandidateNode(s)<br/>📄 decoder_yaml.go:Decode"] --> B["Context{ MatchingNodes: list.List }<br/>seeded with document root node(s)<br/>📄 context.go:Context struct"]
    B --> C["DataTreeNavigator.GetMatchingNodes(context, rootExpressionNode)<br/>📄 data_tree_navigator.go:GetMatchingNodes"]

    C --> D{"expressionNode nil?"}
    D -->|"yes"| E["return context unchanged"]
    D -->|"no"| F["look up Operation.OperationType.Handler<br/>📄 operation.go:operationType.Handler"]

    F --> G["call handler(navigator, context, expressionNode)"]

    G --> H{"handler type,<br/>e.g. pipeOperator,<br/>traverseOperator,<br/>selectOperator..."}

    H -->|"Pipe A | B"| I["recurse: GetMatchingNodes(context, LHS)<br/>→ newContext<br/>then GetMatchingNodes(newContext, RHS)<br/>📄 operator_pipe.go:pipeOperator"]
    H -->|"Traverse .a"| J["for each CandidateNode in context:<br/>look up child key/index 'a'<br/>collect results into new Context<br/>📄 operator_traverse_path.go:traverseOperator / traverseArrayOperator"]
    H -->|"Select(expr)"| K["for each node: recurse into RHS<br/>with node as context;<br/>keep node if result truthy<br/>📄 operator_select.go:selectOperator"]
    H -->|"Binary op (==, +, and)"| L["recurse LHS & RHS,<br/>cross or zip results,<br/>compute per pair<br/>📄 operator_compare.go:equalsOperator / operator_add.go:addOperator"]

    I --> M["returned Context<br/>(new list of matching CandidateNodes)"]
    J --> M
    K --> M
    L --> M

    M --> N["Encoder/Printer renders<br/>each CandidateNode back to YAML/JSON<br/>📄 encoder_yaml.go:Encode / printer.go:PrintResults"]
```

## 4. End-to-end example: `.a.b[] | select(.x > 1)`

- Ties diagrams 2 and 3 together: the same expression is tokenised, converted to postfix, and built into a tree once.
- That tree is then evaluated against the loaded document by recursively descending through `PIPE` → traversal → `SELECT`.
- Only `CandidateNode`s surviving the `select` predicate reach the printer.

```mermaid
flowchart TD
    A["'.a.b[] | select(.x > 1)'"] --> B["Tokenise → infix tokens<br/>📄 lexer_participle.go:Tokenise"]
    B --> C["Shunting-yard → postfix Operations:<br/>a, TRAVERSE, b, TRAVERSE,<br/>[], TRAVERSE_ARRAY,<br/>x, TRAVERSE, 1, GT, SELECT,<br/>PIPE<br/>📄 expression_postfix.go:ConvertToPostfix"]
    C --> D["Build tree:<br/>PIPE(<br/>  TRAVERSE_ARRAY(TRAVERSE(TRAVERSE(SELF,a),b)),<br/>  SELECT(GT(TRAVERSE(SELF,x),1))<br/>)<br/>📄 expression_parser.go:createExpressionTree"]
    D --> E["Load yaml → root CandidateNode<br/>Context = {root}<br/>📄 decoder_yaml.go:Decode"]
    E --> F["GetMatchingNodes(context, PIPE node)<br/>📄 data_tree_navigator.go:GetMatchingNodes"]
    F --> G["LHS traverses a→b→[]<br/>→ Context of array elements<br/>📄 operator_traverse_path.go:traverseOperator / traverseArrayOperator"]
    G --> H["RHS select() runs GT(.x, 1) per element<br/>→ Context of elements where x>1<br/>📄 operator_select.go:selectOperator, operator_compare.go:compareOperator"]
    H --> I["Printer encodes surviving CandidateNodes<br/>📄 printer.go:PrintResults, encoder_yaml.go:Encode"]
```
