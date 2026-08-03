package stub

import "os"

// Delimiters defines the opening and closing markers that wrap a placeholder
// inside a stub. The default is "{{" and "}}", so a placeholder looks like
// "{{ NAME }}". Surrounding whitespace inside the markers is ignored, meaning
// "{{ NAME }}" and "{{NAME}}" refer to the same key.
type Delimiters struct {
	Left  string
	Right string
}

// writePolicy decides what happens when the destination file already exists
// during a Generate call.
type writePolicy int

const (
	// policyError (the default) makes Generate fail with ErrExists when the
	// destination already exists.
	policyError writePolicy = iota
	// policyForce overwrites an existing destination.
	policyForce
	// policySkip leaves an existing destination untouched and returns nil.
	policySkip
	// policyAppend appends the rendered output to an existing destination.
	policyAppend
)

// config holds the resolved settings used by a single render/generate call.
// It is assembled from the supplied Options.
type config struct {
	replaces   map[string]string
	delimiters Delimiters
	policy     writePolicy
	format     bool
	dirPerm    os.FileMode
	filePerm   os.FileMode
}

// defaultDelimiters are used when no WithDelimiters option is provided.
var defaultDelimiters = Delimiters{Left: "{{", Right: "}}"}

// newConfig builds a config from the given options, applying defaults first.
func newConfig(opts ...Option) config {
	c := config{
		replaces:   make(map[string]string),
		delimiters: defaultDelimiters,
		policy:     policyError,
		dirPerm:    0o755,
		filePerm:   0o644,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&c)
		}
	}

	return c
}

// Option customizes a render or generate operation.
type Option func(*config)

// WithReplace registers a single placeholder replacement. The value is
// converted to its string form using the same rules as fmt.Sprint.
func WithReplace(key string, value any) Option {
	return func(c *config) {
		c.replaces[key] = toString(value)
	}
}

// WithReplaces registers many placeholder replacements at once. Values are
// converted to strings using the same rules as fmt.Sprint. Keys merge with
// any previously registered replacements, overriding on conflict.
func WithReplaces(replaces map[string]any) Option {
	return func(c *config) {
		for k, v := range replaces {
			c.replaces[k] = toString(v)
		}
	}
}

// WithDelimiters overrides the placeholder markers. Empty left or right values
// are ignored and fall back to the current delimiters.
func WithDelimiters(left, right string) Option {
	return func(c *config) {
		if left != "" {
			c.delimiters.Left = left
		}
		if right != "" {
			c.delimiters.Right = right
		}
	}
}

// WithForce makes Generate overwrite the destination file when it already
// exists. Without it, generating over an existing file fails with ErrExists.
func WithForce() Option {
	return func(c *config) {
		c.policy = policyForce
	}
}

// WithSkipExisting makes Generate a no-op (returning nil) when the destination
// file already exists, instead of failing.
func WithSkipExisting() Option {
	return func(c *config) {
		c.policy = policySkip
	}
}

// WithAppend makes Generate append the rendered output to the destination file
// when it already exists, rather than replacing it.
func WithAppend() Option {
	return func(c *config) {
		c.policy = policyAppend
	}
}

// WithFormat runs the rendered output through go/format (gofmt) before writing.
// Generate fails if the output is not valid Go source. Use it only for stubs
// that produce Go code.
func WithFormat() Option {
	return func(c *config) {
		c.format = true
	}
}

// WithDirPerm overrides the permission bits used when Generate creates parent
// directories for the destination (default 0o755).
func WithDirPerm(perm os.FileMode) Option {
	return func(c *config) {
		c.dirPerm = perm
	}
}

// WithFilePerm overrides the permission bits used when Generate creates the
// destination file (default 0o644). It applies only when the file is created,
// not when appending to an existing one.
func WithFilePerm(perm os.FileMode) Option {
	return func(c *config) {
		c.filePerm = perm
	}
}
