package stub

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// toString converts an arbitrary value to its string representation.
func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	return fmt.Sprint(v)
}

// apply performs placeholder replacement on content according to cfg.
//
// It scans left to right for delimiter pairs. For each pair, the text between
// the delimiters is trimmed of surrounding whitespace and looked up in the
// replacements map. A match is substituted with its value; an unknown key is
// left in the output verbatim (including its delimiters), and an unterminated
// opening delimiter is copied through unchanged.
//
// It also returns the unresolved keys, in first-seen order with duplicates
// removed, so callers can enforce strict mode.
func apply(content string, cfg config) (string, []string) {
	left, right := cfg.delimiters.Left, cfg.delimiters.Right
	if left == "" || right == "" {
		return content, nil
	}

	var b strings.Builder
	b.Grow(len(content))

	var missing []string
	var seen map[string]struct{}

	for i := 0; i < len(content); {
		open := strings.Index(content[i:], left)
		if open < 0 {
			b.WriteString(content[i:])
			break
		}
		open += i
		b.WriteString(content[i:open])

		contentAfterLeft := open + len(left)
		closeRel := strings.Index(content[contentAfterLeft:], right)
		if closeRel < 0 {
			// No closing delimiter; emit the rest untouched.
			b.WriteString(content[open:])
			break
		}
		c := contentAfterLeft + closeRel

		key := strings.TrimSpace(content[contentAfterLeft:c])
		if val, ok := cfg.replaces[key]; ok {
			b.WriteString(val)
		} else {
			// Preserve the untouched placeholder, delimiters included.
			b.WriteString(content[open : c+len(right)])
			if _, dup := seen[key]; !dup {
				if seen == nil {
					seen = make(map[string]struct{})
				}
				seen[key] = struct{}{}
				missing = append(missing, key)
			}
		}
		i = c + len(right)
	}

	return b.String(), missing
}

// render applies the configuration to content and, when strict mode is on,
// fails with a *MissingKeysError if any placeholder was left unresolved.
func render(content string, cfg config) (string, error) {
	out, missing := apply(content, cfg)
	if cfg.strict && len(missing) > 0 {
		return out, &MissingKeysError{Keys: missing}
	}

	return out, nil
}

// Render reads the stub file at path from the operating system filesystem and
// returns the rendered result after applying the given options.
//
// It fails if the file cannot be read (the error wraps fs.ErrNotExist for a
// missing file) or, under WithStrict, if any placeholder is unresolved (a
// *MissingKeysError, matching errors.Is(err, ErrMissingKeys)).
func Render(path string, opts ...Option) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("stub: read %q: %w", path, err)
	}

	return render(string(data), newConfig(opts...))
}

// RenderFS reads the stub file at path from fsys (for example an embed.FS) and
// returns the rendered result after applying the given options. It mirrors
// Render in every other respect.
func RenderFS(fsys fs.FS, path string, opts ...Option) (string, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return "", fmt.Errorf("stub: read %q from fs: %w", path, err)
	}

	return render(string(data), newConfig(opts...))
}
