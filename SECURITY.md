# Security Policy

## Reporting a Vulnerability

Please **do not** report security vulnerabilities through public GitHub issues.

Instead, use GitHub's private vulnerability reporting feature:
👉 https://github.com/mikefarah/yq/security

This allows vulnerabilities to be triaged and addressed confidentially before any public disclosure.

## Scope

### HTTP / TLS / Network vulnerabilities

yq is a command-line YAML/JSON/TOML processor that reads from files or standard input and writes to standard output. **yq does not include any HTTP or network libraries** and makes no network connections at runtime. CVEs related to HTTP, TLS, or networking are therefore **not applicable** to yq.

### Dependency version bumps

yq uses [Dependabot](https://docs.github.com/en/code-security/dependabot) to automatically raise pull requests for:

- Go module dependencies
- Go toolchain version
- Docker base images

Please **do not** raise pull requests or issues solely to bump dependency or Go versions — Dependabot handles this automatically and the maintainers merge those PRs regularly.
