# go-stub

> A tiny, dependency-free Go toolkit for working with **stub files** — template
> files used to generate boilerplate source code.

Inspired by the "stub" concept popular in frameworks like Laravel: you keep a
template file with placeholders, fill the placeholders with real values, and
generate a concrete file from it. `go-stub` uses only the standard library and
works with stubs on disk or embedded in your binary via `embed.FS`.

## Features

- **Two API styles** over one engine: a functional core (`Render`, `Generate`)
  and a fluent builder (`New().From().To().Generate()`).
- **Stubs from anywhere**: the OS filesystem or any `io/fs.FS` (e.g. `embed.FS`).
- **Configurable placeholders**: default `{{ NAME }}`, whitespace-insensitive,
  with custom delimiters. Unknown keys are left untouched.
- **Safe file generation**: parent dirs auto-created, and explicit write
  policies (error / force / skip / append).
- **`gofmt` on output** for generated Go source.
- **Directory scaffolding**: render a whole stub tree, including placeholders in
  file names.
- **Batch generation** and **string-case helpers** (`ToSnake`, `ToPascal`, …).

## Status

| Phase | Scope | Status |
|-------|-------|--------|
| 0 | Project bootstrap (module, license, CI, layout) | ✅ Done |
| 1 | Render core (`Render`, options, placeholder engine) | ✅ Done |
| 2 | File generation (`Generate`, write policies) | ✅ Done |
| 3 | Fluent API (`New().From().To().Generate()`) | ✅ Done |
| 4 | Advanced (dir scaffolding, batch, case helpers) | ✅ Done |
| 5 | Docs, examples, `v0.1.0` release | ✅ Done |

## Installation

```bash
go get github.com/binafy/go-stub
```

Requires Go 1.21+.

## Quick start

Given a stub file `stubs/model.stub`:

```
package {{ PACKAGE }}

type {{ NAME }} struct{}
```

Render it to a string:

```go
out, err := stub.Render("stubs/model.stub",
    stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "User"}),
)
```

Or generate a file directly:

```go
err := stub.Generate("stubs/model.stub", "models/user.go",
    stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "User"}),
    stub.WithFormat(),
)
```

## Fluent API

```go
err := stub.New().
    From("stubs/model.stub").
    To("models/user.go").
    Replaces(map[string]any{"PACKAGE": "models", "NAME": "User"}).
    Format().
    Generate()
```

## Embedded stubs

```go
//go:embed stubs
var stubs embed.FS

err := stub.New().
    FromFS(stubs, "stubs/model.stub").
    To("models/user.go").
    Replace("NAME", "User").
    Generate()
```

Functional equivalents: `stub.RenderFS(fsys, path, ...)` and
`stub.GenerateFS(fsys, src, dst, ...)`.

## Placeholders

The default delimiters are `{{` and `}}`, and surrounding whitespace is
ignored, so `{{ NAME }}` and `{{NAME}}` are the same key. Values of any type are
formatted with `fmt.Sprint`. **Unknown keys are left in the output verbatim**
rather than erroring. Override the markers with `WithDelimiters`:

```go
stub.Render("t.stub", stub.WithDelimiters("<", ">"), stub.WithReplace("NAME", "User"))
```

## Write policies

When the destination already exists, `Generate` behaves according to the policy
option (default: fail):

| Option | Behavior when destination exists |
|--------|----------------------------------|
| _(none)_ | returns `ErrExists` (nothing is overwritten) |
| `WithForce()` | overwrite |
| `WithSkipExisting()` | leave it, return `nil` |
| `WithAppend()` | append the rendered output |

Check the default failure with `errors.Is(err, stub.ErrExists)`.

## Directory scaffolding

Render an entire tree of stubs at once. Placeholders in **file names** are
rendered too, and `WithTrimSuffix` strips a trailing extension:

```
stubs/module/
  {{ NAME }}.go.stub
  repository/{{ NAME }}_repo.go.stub
```

```go
err := stub.GenerateDir("stubs/module", "internal/user",
    stub.WithReplaces(map[string]any{"PACKAGE": "user", "NAME": "User"}),
    stub.WithTrimSuffix(".stub"),
)
// -> internal/user/User.go, internal/user/repository/User_repo.go
```

`GenerateDirFS(fsys, srcDir, dstDir, ...)` reads the tree from an `fs.FS`.

## Batch generation

Generate many source→destination pairs with shared options; each job may add or
override options and choose its own source filesystem:

```go
err := stub.GenerateJobs([]stub.Job{
    {Src: "stubs/model.stub", Dst: "models/user.go", Opts: []stub.Option{stub.WithReplace("NAME", "User")}},
    {FS: stubs, Src: "stubs/repo.stub", Dst: "models/user_repo.go", Opts: []stub.Option{stub.WithForce()}},
}, stub.WithReplace("PACKAGE", "models"))
```

Jobs run in order and stop at the first error, which names the failing index.

## String-case helpers

Derive multiple placeholder values from a single base name (camelCase and
acronym boundaries are detected, e.g. `HTTPServer` → `http_server`):

```go
stub.ToSnake("UserName")          // "user_name"
stub.ToScreamingSnake("UserName") // "USER_NAME"
stub.ToKebab("UserName")          // "user-name"
stub.ToPascal("user_name")        // "UserName"
stub.ToCamel("user_name")         // "userName"
```

## API reference

| Symbol | Purpose |
|--------|---------|
| `Render`, `RenderFS` | render a stub to a string |
| `Generate`, `GenerateFS` | render a stub to a file |
| `GenerateDir`, `GenerateDirFS` | render a whole stub tree |
| `GenerateJobs`, `Job` | batch generation |
| `New` → `Builder` | fluent API |
| `WithReplace`, `WithReplaces`, `WithDelimiters` | placeholder options |
| `WithForce`, `WithSkipExisting`, `WithAppend` | write policies |
| `WithFormat`, `WithTrimSuffix`, `WithDirPerm`, `WithFilePerm` | output options |
| `ToSnake`, `ToScreamingSnake`, `ToKebab`, `ToPascal`, `ToCamel` | case helpers |
| `ErrExists` | sentinel for an existing destination |

## Development

```bash
make test   # run tests (with -race in CI)
make lint   # go vet + gofmt check
make cover  # coverage report
```

## License

[MIT](LICENSE) © Binafy
