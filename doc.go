// Package stub provides a small, dependency-free toolkit for working with
// stub files: template files used to generate boilerplate source code.
//
// A stub is a plain text file containing placeholders (for example
// "{{ NAME }}"). At generation time the placeholders are replaced with
// concrete values and the result is rendered to a string or written to a
// destination file.
//
// Stubs can be read from the real filesystem or from an embedded
// [io/fs.FS] (such as an embed.FS), which makes the package suitable both
// for local scaffolding and for CLIs that ship their stubs inside the
// binary.
//
// The package exposes two complementary styles:
//
//   - a functional core (stub.Render / stub.Generate with options), and
//   - a fluent, chainable layer (stub.New().From(...).To(...).Generate()).
//
// Both styles are built on the same underlying engine.
package stub
