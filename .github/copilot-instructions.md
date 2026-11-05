# yq review instructions

yq has millions of users, many running it unattended in CI/CD pipelines. Existing behaviour must not change disruptively. Use UK English spelling.

## Backward-compatibility review (highest priority)

For every PR, check whether existing behaviour changes. Treat the following as potential breaking changes:

- Output changes: formatting, whitespace, indentation, comments, quoting, key ordering, null/empty handling.
- Changed stdout/stderr split, exit codes, or error message text.
- Removed or renamed CLI flags, operators, formats, or env vars; changed flag defaults.
- Changed operator semantics, or encoder/decoder behaviour (`pkg/yqlib/encoder_*.go`, `decoder_*.go`).
- New behaviour that is on by default and should be opt-in (flag or preference).
- Removal without a deprecation path (keep the old name working, warn on stderr, document, remove in a later major release).

## Tests as a compatibility signal

1. **Modified expectations.** Explicitly call out any PR that updates, loosens, skips or deletes an existing expectation: Go tests (`pkg/yqlib/*_test.go`, `cmd/*_test.go`), doc-driven examples (`expected` in `expressionScenario`/`formatScenario`, which generate `pkg/yqlib/doc/**`), or acceptance scripts (`acceptance_tests/*.sh`). Flag it as a potential breaking change and ask the author to justify it. Do not accept "the test was wrong" without evidence.
2. **Coverage gaps.** If behaviour changes but no existing test covers the old behaviour, say so and suggest adding characterisation tests that pin the old behaviour (and the new opt-in behaviour) before the change.
3. Bug fixes must include a new test that exposes the bug, alongside the fix itself.

## Intentional behaviour changes

If the PR description states the behaviour change is the purpose (for example fixing formatting), do not object to the change itself. Instead check that: the old behaviour is pinned by a characterisation test; changed expectations match only the stated intent, with no unrelated drift; the change is opt-in, or the description explains user impact and a migration path; and a verdict is still given (usually Medium or High).

## Required verdict

End the review with:

**Compatibility risk: None / Low / Medium / High** followed by concise reasoning listing the affected behaviours, who might be impacted (for example CI scripts parsing output), and any requested mitigations.

Also ensure generated docs under `pkg/yqlib/doc/` are not hand-edited to hide a behaviour change.
