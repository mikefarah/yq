---
applyTo: "**/*_test.go,acceptance_tests/**,scripts/acceptance.sh,examples/**"
---

# Test review: changed expectations signal changed behaviour

Existing tests are the compatibility contract. For any change in these files:

- **Flag every modified, loosened, skipped, commented-out or deleted expectation** (`expected`, `expectedError`, assertions, acceptance script output/exit codes, changed `examples/` inputs). Call it out as a **potential breaking change**, ask the author to justify why users' existing output may change, and do not treat it as routine. Changes to `skipDoc` affect documentation visibility, not test execution; inspect them for hidden output changes rather than treating them alone as changed expectations.
- Watch for subtle weakening: exact match turned into substring, whitespace/comment/quote differences accepted, error text changed, `expected` reordered, scenarios removed instead of fixed.
- Adding new scenarios is fine; ask that they sit alongside, not replace, the old ones. If the PR states the change is intentional, check that edited expectations match only that stated intent.
- If production code changed behaviour but this PR adds no test pinning the old behaviour, request characterisation tests (in `pkg/yqlib/*_test.go` using `expressionScenario`/`formatScenario`, or `acceptance_tests/*.sh` for CLI behaviour such as exit codes, flags and stderr/stdout).
- Changes to doc-driven scenarios regenerate `pkg/yqlib/doc/**`; check whether changed outputs or errors for existing inputs indicate a behaviour change. Prose-only edits and added examples do not by themselves establish one.
- When a bug is found, require a new test that exposes the bug (fails without the fix), included in the same PR as the fix.
