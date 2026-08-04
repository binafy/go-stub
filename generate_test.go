package stub_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	stub "github.com/binafy/go-stub"
)

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	return string(data)
}

func TestGenerateWritesRenderedFile(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "nested", "user.go")

	err := stub.Generate("testdata/model.stub", dst,
		stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "User"}),
	)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	got := readFile(t, dst)
	if !contains(got, "package models") || !contains(got, "func NewUser() *User") {
		t.Errorf("generated file missing rendered content:\n%s", got)
	}
}

func TestGenerateErrorsWhenExists(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out.go")
	if err := os.WriteFile(dst, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := stub.Generate("testdata/model.stub", dst,
		stub.WithReplaces(map[string]any{"PACKAGE": "m", "NAME": "X"}),
	)
	if !errors.Is(err, stub.ErrExists) {
		t.Fatalf("Generate() error = %v, want ErrExists", err)
	}
	if got := readFile(t, dst); got != "original" {
		t.Errorf("destination modified on error: %q", got)
	}
}

func TestGenerateForceOverwrites(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out.go")
	if err := os.WriteFile(dst, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := stub.Generate("testdata/model.stub", dst,
		stub.WithReplaces(map[string]any{"PACKAGE": "m", "NAME": "X"}),
		stub.WithForce(),
	)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got := readFile(t, dst); got == "original" || !contains(got, "package m") {
		t.Errorf("force did not overwrite, got: %q", got)
	}
}

func TestGenerateSkipExisting(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out.go")
	if err := os.WriteFile(dst, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := stub.Generate("testdata/model.stub", dst,
		stub.WithReplaces(map[string]any{"PACKAGE": "m", "NAME": "X"}),
		stub.WithSkipExisting(),
	)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got := readFile(t, dst); got != "original" {
		t.Errorf("skip modified destination: %q", got)
	}
}

func TestGenerateAppend(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out.txt")
	if err := os.WriteFile(dst, []byte("HEAD\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	appendSrc := filepath.Join(t.TempDir(), "line.stub")
	if err := os.WriteFile(appendSrc, []byte("added {{ NAME }}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := stub.Generate(appendSrc, dst,
		stub.WithReplace("NAME", "row"),
		stub.WithAppend(),
	)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got := readFile(t, dst); got != "HEAD\nadded row\n" {
		t.Errorf("append result = %q", got)
	}
}

func TestGenerateWithFormat(t *testing.T) {
	src := filepath.Join(t.TempDir(), "messy.stub")
	// Deliberately badly formatted, but valid, Go source.
	if err := os.WriteFile(src, []byte("package {{ PKG }}\nfunc  F( ){\n}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(t.TempDir(), "formatted.go")

	err := stub.Generate(src, dst, stub.WithReplace("PKG", "x"), stub.WithFormat())
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	want := "package x\n\nfunc F() {\n}\n"
	if got := readFile(t, dst); got != want {
		t.Errorf("formatted output = %q, want %q", got, want)
	}
}

func TestGenerateWithFormatInvalidGo(t *testing.T) {
	src := filepath.Join(t.TempDir(), "invalid.stub")
	if err := os.WriteFile(src, []byte("this is not go"), 0o600); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(t.TempDir(), "out.go")

	err := stub.Generate(src, dst, stub.WithFormat())
	if err == nil {
		t.Fatal("Generate() expected format error, got nil")
	}
	if _, statErr := os.Stat(dst); statErr == nil {
		t.Error("destination should not be created when formatting fails")
	}
}

func TestGenerateMissingSource(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out.go")
	if err := stub.Generate("testdata/nope.stub", dst); err == nil {
		t.Fatal("Generate() expected error for missing source")
	}
}

func TestGenerateFS(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "acct.go")

	err := stub.GenerateFS(embedded, "testdata/model.stub", dst,
		stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "Account"}),
	)
	if err != nil {
		t.Fatalf("GenerateFS() error = %v", err)
	}
	if got := readFile(t, dst); !contains(got, "func NewAccount() *Account") {
		t.Errorf("GenerateFS() content = %q", got)
	}
}

func TestGenerateCustomPerms(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix file-mode bits are not enforced on Windows")
	}
	dst := filepath.Join(t.TempDir(), "sub", "perm.go")

	err := stub.Generate("testdata/model.stub", dst,
		stub.WithReplaces(map[string]any{"PACKAGE": "m", "NAME": "P"}),
		stub.WithDirPerm(0o700),
		stub.WithFilePerm(0o600),
	)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	fi, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Errorf("file perm = %o, want 600", perm)
	}
	di, err := os.Stat(filepath.Dir(dst))
	if err != nil {
		t.Fatal(err)
	}
	if perm := di.Mode().Perm(); perm != 0o700 {
		t.Errorf("dir perm = %o, want 700", perm)
	}
}

func TestGenerateFSMissingSource(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out.go")
	if err := stub.GenerateFS(embedded, "testdata/missing.stub", dst); err == nil {
		t.Fatal("GenerateFS() expected error for missing source")
	}
}
