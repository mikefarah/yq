---
applyTo: "cmd/**/*.go,yq.go,scripts/**,pkg/yqlib/doc/usage/**,pkg/yqlib/doc/operators/headers/**,README.md"
---

# CLI and script compatibility review

Pipelines depend on yq's CLI contract. Flag as potential breaking changes:

- Removed, renamed or re-purposed flags/short flags in `cmd/root.go`; changed flag defaults; changed positional-argument handling.
- Changed exit codes, or output moving between stdout and stderr; changed error or warning text that scripts may match.
- Changed environment variable handling, auto-detection of input/output format (by file extension), or in-place (`-i`) semantics.
- Changed shell completion, version output, or `scripts/` behaviour used by releases and CI.
- New default behaviour that should be opt-in via a new flag.

Removals need a deprecation path: keep the old flag working, warn on stderr (never stdout), document it, and remove only in a later major release.

CLI behaviour is pinned by `acceptance_tests/*.sh` and `cmd/*_test.go`. Call out any edit to those expectations, and request a characterisation test when changed behaviour is not covered. Documentation edits that describe changed behaviour should be matched by a note in the PR.

Finish with **Compatibility risk: None / Low / Medium / High** and reasoning.
