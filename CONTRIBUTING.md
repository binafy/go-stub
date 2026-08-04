# Contributing to go-stub

Thanks for taking the time to contribute! 🎉 This document explains how to get a
change merged smoothly.

## Getting started

```bash
git clone https://github.com/binafy/go-stub
cd go-stub
go test ./...
```

The package depends only on the Go standard library — there is nothing else to
install.

## Development workflow

1. Create a branch: `git switch -c feat/short-description`.
2. Make your change with tests.
3. Run the full local check:

   ```bash
   make test   # go test (add -race when touching I/O or concurrency)
   make lint   # go vet + gofmt check
   ```

4. Open a pull request against `main`.

## Expectations for a pull request

- **Tests pass** on `go test -race ./...`.
- **`gofmt` clean** — run `make fmt` (or `gofmt -w .`).
- **`go vet` and `golangci-lint` clean** — see [`.golangci.yml`](.golangci.yml).
- **New behavior is covered by tests.** We keep coverage high; prefer
  table-driven tests, and use `t.TempDir()` for filesystem tests.
- **Public API changes are documented** with godoc comments and, where useful, a
  testable `Example`.
- **Update the docs** — the `README.md` and `CHANGELOG.md` when behavior changes.

## API stability

`go-stub` follows [Semantic Versioning](https://semver.org). Until `v1.0.0` the
API may still change, but we try to keep changes additive. Please flag any
breaking change explicitly in your PR description.

## Commit messages

Write clear, imperative commit subjects (e.g. "Add strict rendering mode").
Small, focused commits are easier to review than one large one.

## Reporting bugs

Use the issue templates. A minimal reproducer — a stub, the options you passed,
what you expected, and what you got — helps us fix things fast.

## License

By contributing, you agree that your contributions are licensed under the
[MIT License](LICENSE).
