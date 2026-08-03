package stub

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// GenerateDir renders every file under srcDir (read from the operating system
// filesystem) into dstDir, preserving the directory structure.
func GenerateDir(srcDir, dstDir string, opts ...Option) error {
	return generateTree(os.DirFS(srcDir), ".", dstDir, newConfig(opts...))
}

// GenerateDirFS behaves like GenerateDir but reads the source tree from srcDir
// within fsys (for example an embed.FS), writing to dstDir on the operating
// system filesystem.
func GenerateDirFS(fsys fs.FS, srcDir, dstDir string, opts ...Option) error {
	return generateTree(fsys, srcDir, dstDir, newConfig(opts...))
}

// generateTree walks root within fsys and generates every regular file into
// dstDir. Paths inside fsys always use forward slashes (io/fs convention).
func generateTree(fsys fs.FS, root, dstDir string, cfg config) error {
	return fs.WalkDir(fsys, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("stub: walk %q: %w", p, err)
		}
		if d.IsDir() {
			return nil
		}

		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return fmt.Errorf("stub: read %q: %w", p, err)
		}

		rel := relPath(p, root)
		rel = apply(rel, cfg) // render placeholders inside the path
		if cfg.trimSuffix != "" {
			rel = strings.TrimSuffix(rel, cfg.trimSuffix)
		}

		dst := filepath.Join(dstDir, filepath.FromSlash(rel))

		return writeRendered(dst, apply(string(data), cfg), cfg)
	})
}

// relPath returns p expressed relative to root using forward slashes. When root
// is "." (the whole tree), p is already relative.
func relPath(p, root string) string {
	if root == "." || root == "" {
		return p
	}

	return strings.TrimPrefix(p, path.Clean(root)+"/")
}
