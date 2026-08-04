package stub

import (
	"errors"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
)

// Generate renders the stub file at src (read from the operating system
// filesystem) and writes the result to dst, applying the given options.
//
// Parent directories of dst are created as needed. When dst already exists the
// behavior is controlled by the write-policy options; by default Generate
// fails with ErrExists. Under WithStrict an unresolved placeholder fails with a
// *MissingKeysError and nothing is written.
func Generate(src, dst string, opts ...Option) error {
	cfg := newConfig(opts...)

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("stub: read %q: %w", src, err)
	}

	content, err := render(string(data), cfg)
	if err != nil {
		return err
	}

	return writeRendered(dst, content, cfg)
}

// GenerateFS renders the stub file at src (read from fsys, for example an
// embed.FS) and writes the result to dst on the operating system filesystem,
// applying the given options. It behaves like Generate in every other respect.
func GenerateFS(fsys fs.FS, src, dst string, opts ...Option) error {
	cfg := newConfig(opts...)

	data, err := fs.ReadFile(fsys, src)
	if err != nil {
		return fmt.Errorf("stub: read %q from fs: %w", src, err)
	}

	content, err := render(string(data), cfg)
	if err != nil {
		return err
	}

	return writeRendered(dst, content, cfg)
}

// writeRendered applies formatting and the write policy, then persists content to dst.
func writeRendered(dst, content string, cfg config) error {
	if cfg.format {
		formatted, err := format.Source([]byte(content))
		if err != nil {
			return fmt.Errorf("stub: format %q: %w", dst, err)
		}
		content = string(formatted)
	}

	exists, err := fileExists(dst)
	if err != nil {
		return err
	}

	if exists {
		switch cfg.policy {
		case policyError:
			return fmt.Errorf("%w: %q", ErrExists, dst)
		case policySkip:
			return nil
		case policyAppend:
			return appendFile(dst, content)
		case policyForce:
			// fall through to a normal (overwriting) write.
		}
	}

	if err := os.MkdirAll(filepath.Dir(dst), cfg.dirPerm); err != nil {
		return fmt.Errorf("stub: create dir for %q: %w", dst, err)
	}
	if err := os.WriteFile(dst, []byte(content), cfg.filePerm); err != nil {
		return fmt.Errorf("stub: write %q: %w", dst, err)
	}

	return nil
}

// fileExists reports whether path exists. A non-existence error is reported as
// (false, nil); any other stat error is returned.
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, fmt.Errorf("stub: stat %q: %w", path, err)
}

// appendFile appends content to an existing file. It reports a close error only
// when the write itself succeeded, so a flush failure is not silently dropped.
func appendFile(path, content string) (err error) {
	f, openErr := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if openErr != nil {
		return fmt.Errorf("stub: open %q for append: %w", path, openErr)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("stub: close %q: %w", path, cerr)
		}
	}()

	if _, werr := f.WriteString(content); werr != nil {
		return fmt.Errorf("stub: append %q: %w", path, werr)
	}

	return nil
}
