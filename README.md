# go-stub

> A tiny, dependency-free Go toolkit for working with **stub files** — template
> files used to generate boilerplate source code.

Inspired by the "stub" concept popular in frameworks like Laravel: you keep a
template file with placeholders, fill the placeholders with real values, and
generate a concrete file from it.

## Status

🚧 **Work in progress.** The public API is being built in phases and is not
yet stable.

| Phase | Scope | Status |
|-------|-------|--------|
| 0 | Project bootstrap (module, license, CI, layout) | ✅ Done |
| 1 | Render core (`Render`, options, placeholder engine) | ⬜ Planned |
| 2 | File generation (`Generate`, write policies) | ⬜ Planned |
| 3 | Fluent API (`New().From().To().Generate()`) | ⬜ Planned |
| 4 | Advanced (batch, embed helpers, case-aware keys) | ⬜ Planned |
| 5 | Docs, examples, `v0.1.0` release | ⬜ Planned |

## Installation

```bash
go get github.com/binafy/go-stub
```

Requires Go 1.21+.

## Planned usage

The API below is a preview of the target design and may change.

Functional style:

```go
out, err := stub.Render("stubs/model.stub",
    stub.WithReplaces(map[string]any{"NAME": "User"}),
)

err = stub.Generate("stubs/model.stub", "models/user.go",
    stub.WithReplaces(map[string]any{"NAME": "User"}),
)
```

Fluent style:

```go
err := stub.New().
    From("stubs/model.stub").
    To("models/user.go").
    Replaces(map[string]any{"NAME": "User"}).
    Generate()
```

Embedded stubs:

```go
//go:embed stubs/*.stub
var stubs embed.FS

err := stub.New().
    FromFS(stubs, "stubs/model.stub").
    To("models/user.go").
    Replaces(map[string]any{"NAME": "User"}).
    Generate()
```

## Development

```bash
make test   # run tests
make lint   # vet + format check
make cover  # coverage report
```

## License

[MIT](LICENSE) © Binafy
