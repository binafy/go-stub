package stub_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	stub "github.com/binafy/go-stub"
)

var stdReplaces = stub.WithReplaces(map[string]any{"PACKAGE": "m", "NAME": "X"})

// When the destination's parent is a regular file, stat fails with a
// non-"not exist" error, which Generate must surface (covers the error branch
// in fileExists).
func TestGenerateStatErrorWhenParentIsFile(t *testing.T) {
	tmp := t.TempDir()
	blocker := filepath.Join(tmp, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(blocker, "child.go") // parent "blocker" is a file

	err := stub.Generate("testdata/model.stub", dst, stdReplaces)
	if err == nil {
		t.Fatal("expected error when destination parent is a file")
	}
	if errors.Is(err, stub.ErrExists) {
		t.Errorf("unexpected ErrExists: %v", err)
	}
}

// Forcing a write onto a path that is an existing directory must fail at the
// write step (covers the write-error branch in writeRendered).
func TestGenerateForceOntoDirectoryFails(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "adir")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	err := stub.Generate("testdata/model.stub", dir, stdReplaces, stub.WithForce())
	if err == nil {
		t.Fatal("expected error writing over a directory")
	}
}

// Appending to a path that is an existing directory must fail when opening it
// for append (covers the open-error branch in appendFile).
func TestGenerateAppendOntoDirectoryFails(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "adir")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	err := stub.Generate("testdata/model.stub", dir, stdReplaces, stub.WithAppend())
	if err == nil {
		t.Fatal("expected error appending to a directory")
	}
}
