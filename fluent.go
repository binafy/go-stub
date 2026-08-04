package stub

import (
	"errors"
	"io/fs"
	"os"
)

// errNoSource is returned when Render or Generate is called on a Builder that
// has no source set via From or FromFS.
var errNoSource = errors.New("stub: no source set (call From or FromFS)")

// errNoDestination is returned when Generate is called on a Builder that has no
// destination set via To.
var errNoDestination = errors.New("stub: no destination set (call To)")

// Builder is a fluent, chainable front-end over the functional core (Render,
// RenderFS, Generate, GenerateFS). A zero Builder is not usable; start one with
// New.
//
// Every configuration method returns the same Builder so calls can be chained:
//
//	err := stub.New().
//		From("stubs/model.stub").
//		To("models/user.go").
//		Replaces(map[string]any{"NAME": "User"}).
//		Generate()
//
// A Builder is not safe for concurrent modification: build and run each one
// from a single goroutine. Distinct Builders share no state, so independent
// goroutines may each use their own. The package-level functions (Render,
// Generate, and friends) hold no shared state and are safe for concurrent use.
type Builder struct {
	src        string
	fsys       fs.FS // nil means the operating system filesystem
	content    string
	hasContent bool // true when the source is in-memory content, set via Content
	dst        string
	opts       []Option
}

// New returns an empty Builder ready to be configured.
func New() *Builder {
	return &Builder{}
}

// From sets the stub source to a path on the operating system filesystem. It
// replaces any source previously set with From, FromFS, or Content.
func (b *Builder) From(path string) *Builder {
	b.src = path
	b.fsys = nil
	b.hasContent = false

	return b
}

// FromFS sets the stub source to a path within fsys (for example an embed.FS).
// It replaces any source previously set with From, FromFS, or Content.
func (b *Builder) FromFS(fsys fs.FS, path string) *Builder {
	b.src = path
	b.fsys = fsys
	b.hasContent = false

	return b
}

// Content sets the stub source to an in-memory string, rather than a file. It
// replaces any source previously set with From, FromFS, or Content.
func (b *Builder) Content(content string) *Builder {
	b.content = content
	b.hasContent = true
	b.src = ""
	b.fsys = nil

	return b
}

// To sets the destination path (on the operating system filesystem) used by
// Generate.
func (b *Builder) To(path string) *Builder {
	b.dst = path

	return b
}

// Replace registers a single placeholder replacement. See WithReplace.
func (b *Builder) Replace(key string, value any) *Builder {
	b.opts = append(b.opts, WithReplace(key, value))

	return b
}

// Replaces registers many placeholder replacements at once. See WithReplaces.
func (b *Builder) Replaces(replaces map[string]any) *Builder {
	b.opts = append(b.opts, WithReplaces(replaces))

	return b
}

// Delimiters overrides the placeholder markers. See WithDelimiters.
func (b *Builder) Delimiters(left, right string) *Builder {
	b.opts = append(b.opts, WithDelimiters(left, right))

	return b
}

// Force makes Generate overwrite an existing destination. See WithForce.
func (b *Builder) Force() *Builder {
	b.opts = append(b.opts, WithForce())

	return b
}

// SkipExisting makes Generate skip an existing destination. See WithSkipExisting.
func (b *Builder) SkipExisting() *Builder {
	b.opts = append(b.opts, WithSkipExisting())

	return b
}

// Append makes Generate append to an existing destination. See WithAppend.
func (b *Builder) Append() *Builder {
	b.opts = append(b.opts, WithAppend())

	return b
}

// Format runs the rendered output through gofmt before writing. See WithFormat.
func (b *Builder) Format() *Builder {
	b.opts = append(b.opts, WithFormat())

	return b
}

// Strict makes rendering fail on any unresolved placeholder. See WithStrict.
func (b *Builder) Strict() *Builder {
	b.opts = append(b.opts, WithStrict())

	return b
}

// DirPerm overrides the permission bits for created directories. See WithDirPerm.
func (b *Builder) DirPerm(perm os.FileMode) *Builder {
	b.opts = append(b.opts, WithDirPerm(perm))

	return b
}

// FilePerm overrides the permission bits for the created file. See WithFilePerm.
func (b *Builder) FilePerm(perm os.FileMode) *Builder {
	b.opts = append(b.opts, WithFilePerm(perm))

	return b
}

// Render reads and renders the configured source, returning the result. It
// dispatches to RenderContent for an in-memory source, RenderFS when the source
// was set with FromFS, otherwise to Render.
func (b *Builder) Render() (string, error) {
	if b.hasContent {
		return RenderContent(b.content, b.opts...)
	}
	if b.src == "" {
		return "", errNoSource
	}
	if b.fsys != nil {
		return RenderFS(b.fsys, b.src, b.opts...)
	}

	return Render(b.src, b.opts...)
}

// Generate renders the configured source and writes it to the configured
// destination. It dispatches to GenerateFS when the source was set with FromFS,
// GenerateContent for an in-memory source, otherwise to Generate.
func (b *Builder) Generate() error {
	if !b.hasContent && b.src == "" {
		return errNoSource
	}
	if b.dst == "" {
		return errNoDestination
	}
	if b.hasContent {
		cfg := newConfig(b.opts...)
		out, err := render(b.content, cfg)
		if err != nil {
			return err
		}

		return writeRendered(b.dst, out, cfg)
	}
	if b.fsys != nil {
		return GenerateFS(b.fsys, b.src, b.dst, b.opts...)
	}

	return Generate(b.src, b.dst, b.opts...)
}
