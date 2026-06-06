# Contributing to tinyttl

Thanks for your interest in contributing to `tinyttl`.

## Getting started

1. Fork the repository
2. Create a branch for your change
3. Make your changes
4. Run tests locally
5. Open a pull request

## Development

Make sure all tests pass before opening a pull request:

```bash
go test
```

If your change affects performance, please run benchmarks as well:

```bash
go test -run=^$ -bench=. -benchmem
```

## Contribution guidelines

Please keep contributions focused and easy to review.

When possible:

- add or update tests
- update documentation if behavior changes
- keep pull requests scoped to a single improvement
- include benchmark results for performance-related changes

## Pull requests

When opening a pull request, please include:

- a short summary of the change
- why the change is needed
- any relevant test or benchmark results

## Issues

If you find a bug or want to suggest an improvement, please open an issue first when appropriate.

Thanks for helping improve `tinyttl`.