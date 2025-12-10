# Contributing to SusyDB

Thank you for considering contributing to SusyDB! This document outlines how to contribute.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/susydb.git`
3. Create a branch: `git checkout -b feature/your-feature-name`

## Development Setup

```bash
# Build
go build ./...

# Run tests
go test ./...

# Run static analysis
go vet ./...

# Start the server
go run cmd/susydb/main.go --addr :7379
```

## Code Style

- Run `gofmt -s -w .` before committing
- Keep functions small and focused
- Add comments for exported functions
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines

## Pull Request Process

1. Ensure `go build ./...` and `go vet ./...` pass
2. Update documentation if needed
3. Write a clear PR description
4. Link any related issues

## Ideas for Contributions

- Add new commands (e.g., `LPUSH`, `LRANGE`)
- Improve benchmarking tools
- Add more tests
- Documentation improvements
- Performance optimizations

## Questions?

Open an issue with the `question` label.
