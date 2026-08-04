package stub_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	stub "github.com/binafy/go-stub"
)

func TestBuilderRenderFromDisk(t *testing.T) {
	got, err := stub.New().
		From("testdata/model.stub").
		Replaces(map[string]any{"PACKAGE": "models", "NAME": "User"}).
		Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !contains(got, "package models") || !contains(got, "func NewUser() *User") {
		t.Errorf("Render() = %q", got)
	}
}

func TestBuilderRenderFromFS(t *testing.T) {
	got, err := stub.New().
		FromFS(embedded, "testdata/model.stub").
		Replace("PACKAGE", "models").
		Replace("NAME", "Account").
		Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !contains(got, "func NewAccount() *Account") {
		t.Errorf("Render() = %q", got)
	}
}

func TestBuilderGenerate(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "gen", "user.go")

	err := stub.New().
		From("testdata/model.stub").
		To(dst).
		Replaces(map[string]any{"PACKAGE": "models", "NAME": "User"}).
		Format().
		Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got := readFile(t, dst); !contains(got, "func NewUser() *User") {
		t.Errorf("generated = %q", got)
	}
}

func TestBuilderGenerateFromFS(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "acct.go")

	err := stub.New().
		FromFS(embedded, "testdata/model.stub").
		To(dst).
		Replaces(map[string]any{"PACKAGE": "models", "NAME": "Account"}).
		Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got := readFile(t, dst); !contains(got, "type Account struct") {
		t.Errorf("generated = %q", got)
	}
}

func TestBuilderForceAndPolicies(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out.go")
	if err := os.WriteFile(dst, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Default policy should refuse to overwrite.
	err := stub.New().From("testdata/model.stub").To(dst).
		Replaces(map[string]any{"PACKAGE": "m", "NAME": "X"}).
		Generate()
	if !errors.Is(err, stub.ErrExists) {
		t.Fatalf("expected ErrExists, got %v", err)
	}

	// Force should overwrite.
	err = stub.New().From("testdata/model.stub").To(dst).
		Replaces(map[string]any{"PACKAGE": "m", "NAME": "X"}).
		Force().
		Generate()
	if err != nil {
		t.Fatalf("Force Generate() error = %v", err)
	}
	if got := readFile(t, dst); got == "old" {
		t.Error("Force() did not overwrite")
	}
}

func TestBuilderCustomDelimiters(t *testing.T) {
	src := filepath.Join(t.TempDir(), "angle.stub")
	if err := os.WriteFile(src, []byte("Hi <NAME>"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := stub.New().
		From(src).
		Delimiters("<", ">").
		Replace("NAME", "there").
		Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if got != "Hi there" {
		t.Errorf("Render() = %q, want %q", got, "Hi there")
	}
}

func TestBuilderFromOverridesFromFS(t *testing.T) {
	// FromFS then From should switch back to the OS filesystem.
	got, err := stub.New().
		FromFS(embedded, "testdata/model.stub").
		From("testdata/model.stub").
		Replaces(map[string]any{"PACKAGE": "p", "NAME": "N"}).
		Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !contains(got, "package p") {
		t.Errorf("Render() = %q", got)
	}
}

func TestBuilderMissingSource(t *testing.T) {
	if _, err := stub.New().Render(); err == nil {
		t.Error("Render() expected error with no source")
	}
	if err := stub.New().To("x").Generate(); err == nil {
		t.Error("Generate() expected error with no source")
	}
}

func TestBuilderMissingDestination(t *testing.T) {
	err := stub.New().From("testdata/model.stub").Generate()
	if err == nil {
		t.Error("Generate() expected error with no destination")
	}
}

func TestBuilderPerms(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix file-mode bits are not enforced on Windows")
	}
	dst := filepath.Join(t.TempDir(), "d", "perm.go")

	err := stub.New().
		From("testdata/model.stub").
		To(dst).
		Replaces(map[string]any{"PACKAGE": "m", "NAME": "P"}).
		DirPerm(0o700).
		FilePerm(0o600).
		Generate()
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
}

func TestBuilderSkipAndAppend(t *testing.T) {
	// Skip existing.
	dst := filepath.Join(t.TempDir(), "keep.txt")
	if err := os.WriteFile(dst, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(t.TempDir(), "s.stub")
	if err := os.WriteFile(src, []byte("new {{ K }}"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := stub.New().From(src).To(dst).Replace("K", "v").SkipExisting().Generate(); err != nil {
		t.Fatalf("SkipExisting Generate() error = %v", err)
	}
	if got := readFile(t, dst); got != "keep" {
		t.Errorf("SkipExisting modified file: %q", got)
	}

	// Append.
	if err := stub.New().From(src).To(dst).Replace("K", "v").Append().Generate(); err != nil {
		t.Fatalf("Append Generate() error = %v", err)
	}
	if got := readFile(t, dst); got != "keepnew v" {
		t.Errorf("Append result = %q", got)
	}
}
