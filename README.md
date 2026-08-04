<div align="center">

# 🧩 go-stub

### Generate boilerplate from stub templates — with zero dependencies.

*Keep a template file with placeholders, fill them with real values, and stamp out concrete source files. Inspired by Laravel's stubs, built the Go way.*

[![Go Reference](https://pkg.go.dev/badge/github.com/binafy/go-stub.svg)](https://pkg.go.dev/github.com/binafy/go-stub)
[![CI](https://github.com/binafy/go-stub/actions/workflows/ci.yml/badge.svg)](https://github.com/binafy/go-stub/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/binafy/go-stub)](https://goreportcard.com/report/github.com/binafy/go-stub)
[![Go Version](https://img.shields.io/github/go-mod/go-version/binafy/go-stub)](go.mod)
[![Coverage](https://img.shields.io/badge/coverage-97%25-brightgreen)](#)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

</div>

---

## ✨ Why go-stub?

```go
// A stub in, a file out. That's it.
stub.New().
    From("stubs/model.stub").
    To("models/user.go").
    Replaces(map[string]any{"PACKAGE": "models", "NAME": "User"}).
    Format().
    Generate()
```

|                              |                                                               |
|------------------------------|---------------------------------------------------------------|
| 🪶 **Zero dependencies**     | Pure standard library. Nothing to audit, nothing to break.    |
| 🎭 **Two API styles**        | A functional core *and* a fluent builder — over one engine.   |
| 📦 **Stubs from anywhere**   | The OS filesystem or any `io/fs.FS`, including `embed.FS`.    |
| 🛡️ **Safe by default**       | Never overwrites unless you say so. Explicit write policies.  |
| 🌳 **Directory scaffolding** | Render whole trees — placeholders work in file names too.     |
| 🧵 **Batch generation**      | Many files, shared options, one call.                         |
| 🔤 **Case helpers**          | `ToSnake`, `ToPascal`, `ToCamel`… derive names from one base. |
| ✅ **97% test coverage**     | Table-driven tests, race-checked on Linux/macOS/Windows.      |

## 📥 Installation

```bash
go get github.com/binafy/go-stub
```

> Requires **Go 1.21+**. Import as `stub "github.com/binafy/go-stub"`.

## 🚀 Quick start

**1.** Write a stub — `stubs/model.stub`:

```go
package {{ PACKAGE }}

type {{ NAME }} struct {
    ID   int
    Name string
}
```

**2.** Render it to a string…

```go
out, err := stub.Render("stubs/model.stub",
    stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "User"}),
)
```

**3.** …or write it straight to a file, gofmt'd:

```go
err := stub.Generate("stubs/model.stub", "models/user.go",
    stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "User"}),
    stub.WithFormat(),
)
```

```go
// models/user.go
package models

type User struct {
    ID   int
    Name string
}
```

## 🎭 Two styles, one engine

<table>
<tr><th>Functional</th><th>Fluent</th></tr>
<tr>
<td>

```go
stub.Generate(
    "stubs/model.stub",
    "models/user.go",
    stub.WithReplaces(map[string]any{
        "PACKAGE": "models",
        "NAME":    "User",
    }),
    stub.WithFormat(),
)
```

</td>
<td>

```go
stub.New().
    From("stubs/model.stub").
    To("models/user.go").
    Replaces(map[string]any{
        "PACKAGE": "models",
        "NAME":    "User",
    }).
    Format().
    Generate()
```

</td>
</tr>
</table>

Pick whichever reads better at the call site — they do exactly the same thing.

## 📦 Embedded stubs

Ship your stubs *inside* your binary — perfect for CLIs and code generators:

```go
//go:embed stubs
var stubs embed.FS

err := stub.New().
    FromFS(stubs, "stubs/model.stub").
    To("models/user.go").
    Replace("NAME", "User").
    Generate()
```

> Functional equivalents: `stub.RenderFS(fsys, path, …)` and `stub.GenerateFS(fsys, src, dst, …)`.

## 🔖 Placeholders

Default delimiters are `{{ }}`, and **inner whitespace is ignored** — `{{ NAME }}` and `{{NAME}}` are the same key. Values of any type are formatted with `fmt.Sprint`.

```go
stub.Render("t.stub", stub.WithReplace("COUNT", 42))            // 42
stub.Render("t.stub", stub.WithDelimiters("<", ">"),           // custom markers
    stub.WithReplace("NAME", "User"))
```

> 💡 **Unknown keys are left untouched**, not blanked or errored — so a half-configured run never silently drops content.

Want the opposite? `WithStrict()` turns any unresolved placeholder into an error listing exactly which keys were missing:

```go
_, err := stub.Render("t.stub", stub.WithReplace("A", "x"), stub.WithStrict())
if errors.Is(err, stub.ErrMissingKeys) {
    var mk *stub.MissingKeysError
    errors.As(err, &mk)
    fmt.Println("missing:", mk.Keys) // e.g. [B C]
}
```

## 🛡️ Write policies

`Generate` refuses to clobber an existing file unless you opt in:

| Option               | When the destination already exists      |
|----------------------|------------------------------------------|
| _(default)_          | returns `ErrExists` — nothing is written |
| `WithForce()`        | overwrite it                             |
| `WithSkipExisting()` | leave it, return `nil`                   |
| `WithAppend()`       | append the rendered output               |

```go
if errors.Is(err, stub.ErrExists) {
    // ask before overwriting…
}
```

## 🌳 Directory scaffolding

Render an entire tree in one call — and **placeholders in file names are rendered too**. `WithTrimSuffix` drops a trailing extension:

```
stubs/module/
├── {{ NAME }}.go.stub
└── repository/
    └── {{ NAME }}_repo.go.stub
```

```go
stub.GenerateDir("stubs/module", "internal/user",
    stub.WithReplaces(map[string]any{"PACKAGE": "user", "NAME": "User"}),
    stub.WithTrimSuffix(".stub"),
)
```

```
internal/user/
├── User.go
└── repository/
    └── User_repo.go
```

> Reading from an `fs.FS`? Use `stub.GenerateDirFS(fsys, srcDir, dstDir, …)`.

## 🧵 Batch generation

Many source→destination pairs, shared options, per-job overrides:

```go
err := stub.GenerateJobs([]stub.Job{
    {Src: "stubs/model.stub", Dst: "models/user.go",
        Opts: []stub.Option{stub.WithReplace("NAME", "User")}},
    {FS: stubs, Src: "stubs/repo.stub", Dst: "models/user_repo.go",
        Opts: []stub.Option{stub.WithForce()}},
}, stub.WithReplace("PACKAGE", "models")) // shared by every job
```

Jobs run in order and stop at the first error — which tells you exactly which job failed.

## 🔤 String-case helpers

Derive every casing you need from a single base name. camelCase and acronym boundaries are handled (`HTTPServer` → `http_server`):

```go
stub.ToSnake("UserName")          // "user_name"
stub.ToScreamingSnake("UserName") // "USER_NAME"
stub.ToKebab("UserName")          // "user-name"
stub.ToPascal("user_name")        // "UserName"
stub.ToCamel("user_name")         // "userName"
```

```go
base := "user profile"
stub.Generate("stubs/model.stub", "models/user_profile.go",
    stub.WithReplaces(map[string]any{
        "NAME":    stub.ToPascal(base), // UserProfile
        "TABLE":   stub.ToSnake(base),  // user_profile
        "PACKAGE": "models",
    }),
)
```

## 📚 API at a glance

| Symbol                                                                | Purpose                   |
|-----------------------------------------------------------------------|---------------------------|
| `Render` · `RenderFS`                                                 | render a stub to a string |
| `Generate` · `GenerateFS`                                             | render a stub to a file   |
| `GenerateDir` · `GenerateDirFS`                                       | render a whole stub tree  |
| `GenerateJobs` · `Job`                                                | batch generation          |
| `New` → `Builder`                                                     | fluent API                |
| `WithReplace` · `WithReplaces` · `WithDelimiters` · `WithStrict`      | placeholders              |
| `WithForce` · `WithSkipExisting` · `WithAppend`                       | write policies            |
| `WithFormat` · `WithTrimSuffix` · `WithDirPerm` · `WithFilePerm`      | output                    |
| `ToSnake` · `ToScreamingSnake` · `ToKebab` · `ToPascal` · `ToCamel`   | case helpers              |
| `ErrExists` · `ErrMissingKeys` · `MissingKeysError` · `ErrUnsafePath` | error contract            |

📖 Full documentation on [pkg.go.dev](https://pkg.go.dev/github.com/binafy/go-stub). Runnable example in [`examples/basic`](examples/basic).

## 🛠️ Development

```bash
make test   # run tests (with -race in CI)
make lint   # go vet + gofmt check
make cover  # coverage report
```

## 🤝 Contributing

Issues and pull requests are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for the workflow, and make sure `make test` and `make lint` pass before opening a PR. Security reports go through [SECURITY.md](SECURITY.md).

## 📄 License

Released under the [MIT License](LICENSE) © Binafy.

<div align="center">
<sub>Built with ❤️ for gophers who hate writing the same file twice.</sub>
</div>
