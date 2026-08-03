package stub

// Delimiters defines the opening and closing markers that wrap a placeholder
// inside a stub. The default is "{{" and "}}", so a placeholder looks like
// "{{ NAME }}". Surrounding whitespace inside the markers is ignored, meaning
// "{{ NAME }}" and "{{NAME}}" refer to the same key.
type Delimiters struct {
	Left  string
	Right string
}

// config holds the resolved settings used by a single render/generate call.
// It is assembled from the supplied Options.
type config struct {
	replaces   map[string]string
	delimiters Delimiters
}

// defaultDelimiters are used when no WithDelimiters option is provided.
var defaultDelimiters = Delimiters{Left: "{{", Right: "}}"}

// newConfig builds a config from the given options, applying defaults first.
func newConfig(opts ...Option) config {
	c := config{
		replaces:   make(map[string]string),
		delimiters: defaultDelimiters,
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
