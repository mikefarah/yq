1. Install (golang)[https://golang.org/]
1. Run  `scripts/devtools.sh` to install the required devtools
2. Run `make [local] vendor` to install the vendor dependencies
2. Run `make [local] test` to ensure you can run the existing tests
3. Write unit tests - (see existing examples). Changes will not be accepted without corresponding unit tests.
4. Make the code changes.
5. `make [local] test` to lint code and run tests
6. Profit! ok no profit, but raise a PR and get kudos :)
