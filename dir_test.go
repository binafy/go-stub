package stub_test

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	stub "github.com/binafy/go-stub"
)

//go:embed testdata/scaffold
var embeddedDir embed.FS

func TestGenerateDir(t *testing.T) {
	dstDir := t.TempDir()

	err := stub.GenerateDir("testdata/scaffold", dstDir,
		stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "User"}),
		stub.WithTrimSuffix(".stub"),
	)
	if err != nil {
		t.Fatalf("GenerateDir() error = %v", err)
	}

	// File name placeholder + suffix trimmed: "{{ NAME }}.go.stub" -> "User.go".
	top := filepath.Join(dstDir, "User.go")
	if got := readFile(t, top); !contains(got, "package models") || !contains(got, "type User struct{}") {
		t.Errorf("top file content = %q", got)
	}

	// Nested structure preserved and suffix trimmed.
	inner := filepath.Join(dstDir, "inner", "note.txt")
	if got := readFile(t, inner); got != "inner User\n" {
		t.Errorf("inner file content = %q", got)
	}
}

func TestGenerateDirNoTrim(t *testing.T) {
	dstDir := t.TempDir()

	err := stub.GenerateDir("testdata/scaffold", dstDir,
		stub.WithReplaces(map[string]any{"PACKAGE": "m", "NAME": "X"}),
	)
	if err != nil {
		t.Fatalf("GenerateDir() error = %v", err)
	}

	// Without WithTrimSuffix the ".stub" extension is kept.
	if _, err := os.Stat(filepath.Join(dstDir, "X.go.stub")); err != nil {
		t.Errorf("expected X.go.stub to exist: %v", err)
	}
}

func TestGenerateDirFS(t *testing.T) {
	dstDir := t.TempDir()

	err := stub.GenerateDirFS(embeddedDir, "testdata/scaffold", dstDir,
		stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "Account"}),
		stub.WithTrimSuffix(".stub"),
	)
	if err != nil {
		t.Fatalf("GenerateDirFS() error = %v", err)
	}

	if got := readFile(t, filepath.Join(dstDir, "Account.go")); !contains(got, "type Account struct{}") {
		t.Errorf("content = %q", got)
	}
	if got := readFile(t, filepath.Join(dstDir, "inner", "note.txt")); got != "inner Account\n" {
		t.Errorf("inner content = %q", got)
	}
}

func TestGenerateDirMissing(t *testing.T) {
	if err := stub.GenerateDir("testdata/no-such-dir", t.TempDir()); err == nil {
		t.Fatal("GenerateDir() expected error for missing source dir")
	}
}
