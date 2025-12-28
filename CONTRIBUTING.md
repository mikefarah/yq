# Before you begin 
Not all new PRs will be merged in 

It's recommended to check with the owner first (e.g. raise an issue) to discuss a new feature before developing, to ensure your hard efforts don't go to waste.

PRs to fix bugs and issues are almost always welcome :pray: please ensure you write tests as well.

The following types of PRs will _not_ be accepted:
- **Significant refactors** take a lot of time to understand and can have all sorts of unintended side effects. If you think there's a better way to do things (that requires significant changes) raise an issue for discussion first :)
- **Release pipeline PRs** are a security risk - it's too easy for a serious vulnerability to sneak in (either intended or not). If there is a new cool way of releasing things, raise an issue for discussion first - it will need to be gone over with a fine tooth comb.
- **Version bumps** are handled by dependabot, the bot will auto-raise PRs and they will be regularly merged in. 
- **New release platforms** At this stage, yq is not going to maintain any other release platforms other than GitHub and Docker - that said, I'm more than happy to put in other community maintained methods in the README for visibility :heart:


# Development

## Initial Setup

1. Install [Golang](https://golang.org/) (version 1.24.0 or later)
2. Run `scripts/devtools.sh` to install required development tools:
   - golangci-lint for code linting
   - gosec for security analysis
3. Run `make [local] vendor` to install vendor dependencies
4. Run `make [local] test` to ensure you can run the existing tests

## Development Workflow

1. **Write unit tests first** - Changes will not be accepted without corresponding unit tests (see Testing section below)
2. **Make your code changes**
3. **Run tests and linting**: `make [local] test` (this runs formatting, linting, security checks, and tests)
4. **Create your PR** and get kudos! :)

## Make Commands

- Use `make [local] <command>` for local development (runs in Docker container)
- Use `make <command>` for CI/CD environments
- Common commands:
  - `make [local] vendor` - Install dependencies
  - `make [local] test` - Run all checks and tests
  - `make [local] build` - Build the yq binary
  - `make [local] format` - Format code
  - `make [local] check` - Run linting and security checks

# Code Quality

## Linting and Formatting

The project uses strict linting rules defined in `.golangci.yml`. All code must pass:

- **Code formatting**: gofmt, goimports, gci
- **Linting**: revive, errorlint, gosec, misspell, and others
- **Security checks**: gosec security analysis
- **Spelling checks**: misspell detection

Run `make [local] check` to verify your code meets all quality standards.

## Code Style Guidelines

- Follow standard Go conventions
- Use meaningful variable names
- Add comments for public functions and complex logic
- Keep functions focused and reasonably sized
- Use the project's existing patterns and conventions

# Testing

## Test Structure

Tests in yq use the `expressionScenario` pattern. Each test scenario includes:
- `expression`: The yq expression to test
- `document`: Input YAML/JSON (optional)
- `expected`: Expected output
- `skipDoc`: Whether to skip documentation generation

## Writing Tests

1. **Find the appropriate test file** (e.g., `operator_add_test.go` for addition operations)
2. **Add your test scenario** to the `*OperatorScenarios` slice
3. **Run the specific test**: `go test -run TestAddOperatorScenarios` (replace with appropriate test name)
4. **Verify documentation generation** (see Documentation section)

## Test Examples

```go
var addOperatorScenarios = []expressionScenario{
    {
        skipDoc:    true,
        expression: `"foo" + "bar"`,
        expected: []string{
            "D0, P[], (!!str)::foobar\n",
        },
    },
    {
        document:   "apples: 3",
        expression: `.apples + 3`,
        expected: []string{
            "D0, P[apples], (!!int)::6\n",
        },
    },
}
```

## Running Tests

- **All tests**: `make [local] test`
- **Specific test**: `go test -run TestName`
- **With coverage**: `make [local] cover`

# Documentation

## Documentation Generation

The project uses a documentation system that combines static headers with dynamically generated content from tests.

### How It Works

1. **Static headers** are defined in `pkg/yqlib/doc/operators/headers/*.md`
2. **Dynamic content** is generated from test scenarios in `*_test.go` files
3. **Generated docs** are created in `pkg/yqlib/doc/*.md` by concatenating headers with test-generated content
4. **Documentation is synced** to the gitbook branch for the website

### Updating Operator Documentation

#### For Test-Generated Documentation

Most operator documentation is generated from tests. To update:

1. **Find the test file** (e.g., `operator_add_test.go`)
2. **Update test scenarios** - each `expressionScenario` with `skipDoc: false` becomes documentation
3. **Run the test** to regenerate docs:
   ```bash
   cd pkg/yqlib
   go test -run TestAddOperatorScenarios
   ```
4. **Verify the generated documentation** in `pkg/yqlib/doc/add.md`
5. **Create a PR** with your changes

#### For Header-Only Documentation

If documentation exists only in `headers/*.md` files:
1. **Update the header file directly** (e.g., `pkg/yqlib/doc/operators/headers/add.md`)
2. **Create a PR** with your changes

### Updating Static Documentation

For documentation not in the master branch:

1. **Check the gitbook branch** for additional pages
2. **Update the `*.md` files** directly
3. **Create a PR** to the gitbook branch

### Documentation Best Practices

- **Write clear, concise examples** in test scenarios
- **Use meaningful variable names** in examples
- **Include edge cases** and error conditions
- **Test your documentation changes** by running the specific test
- **Verify generated output** matches expectations

Note: PRs with small changes (e.g. minor typos) may not be merged (see https://joel.net/how-one-guy-ruined-hacktoberfest2020-drama).

# Troubleshooting

## Common Setup Issues

### Docker/Podman Issues
- **Problem**: `make` commands fail with Docker errors
- **Solution**: Ensure Docker or Podman is running and accessible
- **Alternative**: Use `make local <command>` to run in containers

### Go Version Issues
- **Problem**: Build fails with Go version errors
- **Solution**: Ensure you have Go 1.24.0 or later installed
- **Check**: Run `go version` to verify

### Vendor Dependencies
- **Problem**: `make vendor` fails or dependencies are outdated
- **Solution**: 
  ```bash
  go mod tidy
  make [local] vendor
  ```

### Linting Failures
- **Problem**: `make check` fails with linting errors
- **Solution**: 
  ```bash
  make [local] format  # Auto-fix formatting
  # Manually fix remaining linting issues
  make [local] check   # Verify fixes
  ```

### Test Failures
- **Problem**: Tests fail locally but pass in CI
- **Solution**: 
  ```bash
  make [local] test    # Run in Docker container
  ```

- **Problem**: Tests fail with a VCS error:
  ```bash
  error obtaining VCS status: exit status 128
  Use -buildvcs=false to disable VCS stamping.
  ```
- **Solution**:
  Git security mechanisms prevent Golang from detecting the Git details inside
  the container; either build with the `local` option, or pass GOFLAGS to
  disable Golang buildvcs behaviour.
  ```bash
  make local test
  # OR
  make test GOFLAGS='-buildvcs=true'
  ```

### Documentation Generation Issues
- **Problem**: Generated docs don't update after test changes
- **Solution**: 
  ```bash
  cd pkg/yqlib
  go test -run TestSpecificOperatorScenarios
  # Check if generated file updated in pkg/yqlib/doc/
  ```

## Getting Help

- **Check existing issues**: Search GitHub issues for similar problems
- **Create an issue**: If you can't find a solution, create a detailed issue
- **Ask questions**: Use GitHub Discussions for general questions
- **Join the community**: Check the project's community channels