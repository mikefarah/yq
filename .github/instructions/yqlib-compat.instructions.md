---
applyTo: "pkg/yqlib/operator_*.go,pkg/yqlib/operation.go,pkg/yqlib/lexer_participle.go,pkg/yqlib/encoder_*.go,pkg/yqlib/decoder_*.go,pkg/yqlib/format.go,pkg/yqlib/candidate_node.go,pkg/yqlib/printer*.go"
---

# Library compatibility review

Check whether the change alters results for expressions and documents that work today:

- **Operators:** changed result for existing input, precedence changes in `operation.go`, lexer changes in `lexer_participle.go` that alter how existing expressions parse, removed or renamed operators/aliases, changed null/empty/missing-key handling, changed ordering, changed error vs. empty-result behaviour.
- **Encoders/decoders:** changed whitespace, indentation, quoting/style, comment placement, key order, number/boolean/null formatting, document separators, or round-trip fidelity. Even "more correct" output changes can break CI diffs.
- **Formats:** removed or renamed format names/aliases in `format.go`; build tags (`yq_no*`) must keep working.
- **Shared code** (`candidate_node.go`, printers, navigation): small tweaks affect every operator; ask which existing tests prove nothing else changed.

If behaviour must change, prefer an opt-in preference/flag, keeping the old behaviour as default, and a deprecation path rather than removal. Request characterisation tests of the old behaviour where `pkg/yqlib/*_test.go` does not already cover it, and call out any edited existing expectation as a potential breaking change.

Finish with **Compatibility risk: None / Low / Medium / High** and reasoning.
