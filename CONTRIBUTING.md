# Before you begin 
Not all new PRs will be merged in 

It's recommended to check with the owner first (e.g. raise an issue) to discuss a new feature before developing, to ensure your hard efforts don't go to waste.

PRs to fix bugs and issues are almost always be welcome, just make sure you write tests as well.

** PRs that significantly refactor code and release pipelines PRs will generally _NOT_ be accepted **

Significant refactors take a lot of time to understand and can have all sorts of unintended side effects.

Release pipeline PRs are a security risk - it's too easy for a serious vulnerability to sneak in (either intended or not). If there is a new cool way of releasing things, raise an issue for discussion first - it will need to be gone over with a fine tooth comb.

At this stage, yq is not going to maintain any other release platforms other than GitHub and Docker - that said, I'm more than happy to put in other community maintained methods in the README for visibility :heart:


# Development

1. Install (Golang)[https://golang.org/]
1. Run  `scripts/devtools.sh` to install the required devtools
2. Run `make [local] vendor` to install the vendor dependencies
2. Run `make [local] test` to ensure you can run the existing tests
3. Write unit tests - (see existing examples). Changes will not be accepted without corresponding unit tests.
4. Make the code changes.
5. `make [local] test` to lint code and run tests
6. Profit! ok no profit, but raise a PR and get kudos :)


# Documentation

The various operator documentation (e.g. 'strings') are generated from the 'master' branch, and have a statically defined header (e.g. `pkg/yqlib/doc/operators/headers/add.md`) and the bulk of the docs are generated from the unit tests e.g. `pkg/yqlib/operator_add_test.go`.

The pipeline will run the tests and automatically concatenate the files together, and put them under 
`pkg/qylib/doc/add.md`. These files are checked in the master branch (and are copied to the gitbook branch as part of the release process).


Remaining static documentation is in the 'githook' branch (where the generated docs are copied across into)

## How to contribute

The first step is to find if what you want is automatically generated or not - start by looking in the master branch. 

Note that PRs with small changes (e.g. minor typos) may not be merged (see https://joel.net/how-one-guy-ruined-hacktoberfest2020-drama).

### Updating dynamic documentation from master
- Search for the documentation you want to update. If you find matches in a `*_test.go` file - update that, as that will automatically update the matching `*.md` file 
- Assuming you are updating a `*_test.go` file, once updated, run the test to regenerated the docs. E.g. for the 'Add' test generated docs, from the pkg/yqlib folder run:
`go test -run TestAddOperatorScenarios` which will run that test defined in the `operator_add_test.go` file.
- Ensure the tests still pass, and check the generated documentation have your update.
- Note: If the documentation is only in a `headers/*.md` file, then just update that directly
- Raise a PR to merge the changes into master!

### Updating static documentation from the gitbook branch
If you haven't found what you want to update in the master branch, then check the gitbook branch directly as there are a few pages in there that are not in master.

- Update the `*.md` files
- Raise a PR to merge the changes into gitbook.