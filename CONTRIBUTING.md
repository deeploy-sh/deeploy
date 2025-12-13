# Contributing

Thanks for considering contributing to deeploy.

## Before You Contribute

deeploy is open source under the Apache 2.0 license. We may offer commercial products in the future (hosted solutions, paid plugins, etc.). By contributing, you agree that your contributions may be used in both open-source and commercial aspects of the project.

Questions? Reach out before contributing.

## Development

### Prerequisites

- Go 1.23+
- [Task](https://taskfile.dev)

### Run

```bash
task dev:server  # Server daemon
task dev:tui     # TUI client
```

## Style Guide

### Go

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- No inline conditionals - use multiline

```go
// Bad
if err := doSomething(); err != nil {
    return err
}

// Good
err := doSomething()
if err != nil {
    return err
}
```

### Git

- Present tense ("Add feature" not "Added feature")
- Imperative mood ("Fix bug" not "Fixes bug")
- Keep it short

## Issues & PRs

- Clear, descriptive titles
- Explain the why, not just the what
