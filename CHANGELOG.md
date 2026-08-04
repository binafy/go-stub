# Changelog

All notable changes to this project are documented in this file. The format is
based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Security

- `GenerateDir` / `GenerateDirFS` now reject rendered file names that would
  escape the destination directory (e.g. a `../` in a placeholder value),
  returning the new `ErrUnsafePath` sentinel.

### Changed

- Unexported the internal `Delimiters` type; the public way to set markers is
  still `WithDelimiters`. (No effect on existing code — the type was never used
  in an exported signature.)

### Added

- Render core: `Render` and `RenderFS` with a whitespace-insensitive
  placeholder engine (default `{{ }}` delimiters).
- Options: `WithReplace`, `WithReplaces`, `WithDelimiters`.
- Strict mode: `WithStrict` fails on unresolved placeholders, reported via
  `MissingKeysError` / the `ErrMissingKeys` sentinel.
- File generation: `Generate` and `GenerateFS` with auto-created parent
  directories.
- Write policies: `WithForce`, `WithSkipExisting`, `WithAppend`, and the
  `ErrExists` sentinel (default behavior fails on an existing destination).
- Output options: `WithFormat` (gofmt), `WithTrimSuffix`, `WithDirPerm`,
  `WithFilePerm`.
- Fluent builder: `New` / `Builder` mirroring the functional API.
- Directory scaffolding: `GenerateDir` / `GenerateDirFS`, rendering placeholders
  in both file contents and file names.
- Batch generation: `GenerateJobs` and the `Job` type.
- String-case helpers: `ToSnake`, `ToScreamingSnake`, `ToKebab`, `ToPascal`,
  `ToCamel`.
- Runnable example under `examples/basic` and testable godoc examples.

[Unreleased]: https://github.com/binafy/go-stub/commits/main
